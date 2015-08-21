# Copyright 2015 ETH Zurich
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
:mod:`zookeeper` --- Library for interfacing with Zookeeper
===========================================================
"""
# Stdlib
import logging
import os.path
import threading
import time

# External packages
from kazoo.client import KazooClient, KazooRetry, KazooState
from kazoo.exceptions import (
    ConnectionLoss,
    LockTimeout,
    NoNodeError,
    NodeExistsError,
    SessionExpiredError,
)
from kazoo.handlers.threading import KazooTimeoutError

# SCION
from lib.errors import SCIONBaseError
from lib.thread import kill_self, thread_safety_net
from lib.util import timed


class ZkBaseError(SCIONBaseError):
    """
    Base exception class for all lib.zookeeper exceptions
    """
    pass


class ZkConnectionLoss(ZkBaseError):
    """Connection to Zookeeper is lost"""
    pass


class ZkNoNodeError(ZkBaseError):
    """A node doesn't exist"""
    pass


class ZkRetryLimit(ZkBaseError):
    """Operation retries hit limit"""
    pass


class Zookeeper(object):
    """
    A wrapper class for Zookeeper interfacing, using the `Kazoo python library
    <https://kazoo.readthedocs.org/en/latest/index.html>`_.

    As Kazoo's functionality is mostly unaware of connection-state changes,
    it requires quite a bit of delicate code to make work reliably.

    E.g. Kazoo's Lock will claim to be held, even if the Zookeeper connection
    has been lost in the meantime. This causes an immediate split-brain problem
    for anything relying on that lock for synchronization. This is also,
    unfortunately, no way to inform the local Lock object that the connection
    is down and therefore the Lock should be released.

    All of Kazoo's events are done via callbacks. These callbacks must not
    block. If they do, no more Kazoo events can happen.

    E.g. if a watch callback blocks, disconnection callbacks will not run.
    """

    def __init__(self, isd_id, ad_id, srv_type, srv_id,
                 zk_hosts, timeout=1.0, on_connect=None,
                 on_disconnect=None, handle_paths=None):
        """
        Setup the Zookeeper connection.

        :param int isd_id: The ID of the current ISD.
        :param int ad_id: The ID of the current AD.
        :param str srv_type: a service type from
                             :const:`lib.defines.SERVICE_TYPES`
        :param str srv_id: Service instance identifier.
        :param list zk_hosts: List of Zookeeper instances to connect to, in the
                              form of ``["host:port"..]``.
        :param float timeout: Zookeeper session timeout length (in seconds).
        :param on_connect: A function called everytime a connection is made to
                           Zookeeper.
        :param on_disconnect: A function called everytime a connection is lost
                              to Zookeeper.
        :param tuple handle_paths: A list of tuples of ZK paths, their
                                   corresponding handler functions, and sync
                                   states. It is ensured that paths exist on
                                   connect.
        """
        self._isd_id = isd_id
        self._ad_id = ad_id
        self._srv_id = srv_id
        self._timeout = timeout
        self._on_connect = on_connect
        self._on_disconnect = on_disconnect
        self._shared_caches = []
        if handle_paths:
            for path, handler, state_synced in handle_paths:
                shared_cache = ZkSharedCache(self, path, handler, state_synced)
                self._shared_caches.append(shared_cache)
        self._prefix = "/ISD%d-AD%d/%s" % (
            self._isd_id, self._ad_id, srv_type)
        # Keep track of our connection state
        self._connected = threading.Event()
        # Keep track of the kazoo lock
        self._lock = threading.Event()
        # Used to signal connection state changes
        self._state_event = threading.Semaphore(value=0)
        # Kazoo parties
        self._parties = {}
        # Kazoo lock (initialised later)
        self._zk_lock = None

        self._kazoo_setup(zk_hosts)
        self._setup_state_listener()
        self._kazoo_start()

    def _kazoo_setup(self, zk_hosts):
        """
        Create and configure Kazoo client

        :param list zk_hosts: List of Zookeeper instances to connect to, in the
                              form of ``["host:port"..]``.
        """
        # Disable exponential back-off
        kretry = KazooRetry(max_tries=-1, max_delay=1)
        # Stop kazoo from drowning the log with debug spam:
        logger = logging.getLogger("KazooClient")
        logger.setLevel(logging.ERROR)
        # (For low-level kazoo debugging):
        # import kazoo.loggingsupport
        # logger.setLevel(kazoo.loggingsupport.BLATHER)

        self._zk = KazooClient(hosts=",".join(zk_hosts),
                               timeout=self._timeout,
                               connection_retry=kretry,
                               logger=logger)

    def _kazoo_start(self):
        """
        Connect the Kazoo client to Zookeeper
        """
        logging.info("Connecting to Zookeeper")
        try:
            self._zk.start()
        except KazooTimeoutError:
            logging.critical(
                "Timed out connecting to Zookeeper on startup, exiting")
            kill_self()

    def _setup_state_listener(self):
        """
        Spawn state listener thread, to respond to state change notifications
        from Kazoo. We use a thread, as the listener callback must not block.
        """
        threading.Thread(
            target=thread_safety_net, args=(self._state_handler,),
            name="libZK._state_handler", daemon=True).start()
        # Listener called every time connection state changes
        self._zk.add_listener(self._state_listener)

    def _state_listener(self, new_state):
        """
        Called everytime the Kazoo connection state changes.
        """
        # Signal a connection state change
        self._state_event.release()
        # Tell kazoo not to remove this listener:
        return False

    def _state_handler(self, initial_state="startup"):
        """
        A thread worker function to wait for Kazoo connection state changes,
        and call the relevant method.
        """
        old_state = initial_state
        while True:
            # Wait for connection state change
            self._state_event.acquire()
            # Short-circuit handler if the state hasn't actually changed
            if old_state == self._zk.state:
                continue
            logging.debug("Kazoo old state: %s, new state: %s",
                          old_state, self._zk.state)
            old_state = self._zk.state
            if self._zk.state == KazooState.CONNECTED:
                self._state_connected()
            elif self._zk.state == KazooState.SUSPENDED:
                self._state_suspended()
            else:
                self._state_lost()

    def _state_connected(self):
        """
        Handles the Kazoo 'connected' event.
        """
        # Might be first connection, or reconnecting after a problem.
        logging.debug("Connection to Zookeeper succeeded (Session: %s)",
                      hex(self._zk.client_id[0]))
        try:
            self.ensure_path(self._prefix, abs=True)
            for shared_cache in self._shared_caches:
                self.ensure_path(shared_cache.path)
            for party in self._parties.values():
                party.autojoin()
        except ZkConnectionLoss:
            return
        self._connected.set()
        if self._on_connect:
            self._on_connect()

    def _state_suspended(self):
        """
        Handles the Kazoo 'connection suspended' event.

        This means that the connection to Zookeeper is down.
        """
        self._connected.clear()
        self._lock.clear()
        logging.info("Connection to Zookeeper suspended")
        if self._on_disconnect:
            self._on_disconnect()

    def _state_lost(self):
        """
        Handles the Kazoo 'connection lost' event.

        This means that the Zookeeper session is lost, so all setup needs to be
        re-done on connect.
        """
        self._connected.clear()
        self._lock.clear()
        logging.info("Connection to Zookeeper lost")
        if self._on_disconnect:
            self._on_disconnect()

    def is_connected(self):
        """
        Check if there is currently a connection to Zookeeper.
        """
        return self._connected.is_set()

    def wait_connected(self, timeout=None):
        """
        Wait until there is a connection to Zookeeper.
        """
        return self._connected.wait(timeout=timeout)

    def ensure_path(self, path, abs=False):
        """
        Ensure that a path exists in Zookeeper.

        :param str path: Path to ensure
        :param bool abs: Is the path abolute or relative?
        """
        full_path = path
        if not abs:
            full_path = os.path.join(self._prefix, path)
        try:
            self._zk.ensure_path(full_path)
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss

    def party_setup(self, prefix=None, autojoin=True):
        """
        Setup a `Kazoo Party
        <https://kazoo.readthedocs.org/en/latest/api/recipe/party.html>`_.

        Used to signal that a group of processes are in a similar state.

        :param str prefix: Path to create the party under. If not specified,
                           uses the default prefix for this server instance.
        :param bool autojoin: Join the party if True, also on reconnect
        :return: a ZkParty object
        :rtype: ZkParty
        :raises:
            ZkConnectionLoss: if the connection to ZK drops
        """
        if not self.is_connected():
            raise ZkConnectionLoss
        if prefix is None:
            prefix = self._prefix
        party_path = os.path.join(prefix, "party")
        self.ensure_path(party_path, abs=True)
        party = ZkParty(self._zk, party_path, self._srv_id, autojoin)
        self._parties[party_path] = party
        return party

    def get_lock(self, timeout=60.0):
        """
        Try to get the lock. Returns immediately if we already have the lock.

        :param float timeout: Time (in seconds) to wait for lock acquisition,
                              or ``None`` to wait forever.
        :type timeout: float or None.
        :return: ``True`` if we got the lock, or already had it, otherwise
                 ``False``.
        :rtype: :class:`bool`
        """
        if self._zk_lock is None:
            # First-time setup.
            lock_path = os.path.join(self._prefix, "lock")
            self._zk_lock = self._zk.Lock(lock_path, self._srv_id)
        if not self.is_connected():
            self.release_lock()
            return False
        elif self._lock.is_set():
            # We already have the lock
            return True
        try:
            if self._zk_lock.acquire(timeout=timeout):
                self._lock.set()
            else:
                pass
        except (LockTimeout, ConnectionLoss, SessionExpiredError):
            pass
        ret = self.have_lock()
        return ret

    def release_lock(self):
        """
        Release the lock
        """
        self._lock.clear()
        if self.is_connected():
            try:
                self._zk_lock.release()
            except (NoNodeError, ConnectionLoss, SessionExpiredError):
                pass
        # Hack suggested by https://github.com/python-zk/kazoo/issues/2
        self._zk_lock.is_acquired = False

    def have_lock(self):
        """
        Check if we currently hold the lock
        """
        return self.is_connected() and self._lock.is_set()

    def wait_lock(self):
        """
        Wait until we hold the lock
        """
        self._lock.wait()

    def store_shared_item(self, path, name, value):
        """
        Store an item in a shared path.

        :param str path: The path to store the item in. E.g. ``"shared"``
        :param str name: A name for the item entry. E.g. ``"item01"``
        :param bytes value: The value to store in the item.
        :raises:
            ZkConnectionLoss: if the connection to ZK drops
        """
        if not self.is_connected():
            raise ZkConnectionLoss
        path = os.path.join(self._prefix, path)
        # First, assume the path already exists (the normal case)
        try:
            self._zk.set("%s/%s" % (path, name), value)
            return
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss
        except NoNodeError:
            pass
        # Node doesn't exist, so create it instead.
        try:
            self._zk.create("%s/%s" % (path, name), value)
            return
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss
        except NodeExistsError:
            # Node was created between our check and our create, so assume that
            # the contents are recent, and return without error.
            return

    def get_shared_item(self, path, entry):
        """
        Retrieve a specific item from a shared path.

        :param str path: The path the item is stored in. E.g. ``"shared"``
        :param str entry: The name of the entry. E.g. ``"pcb0000002046"``
        :return: The value of the item
        :rtype: :class:`bytes`
        :raises:
            ZkConnectionLoss: if the connection to ZK drops
            ZkNoNodeError: if the entry does not exist
        """
        if not self.is_connected():
            raise ZkConnectionLoss
        entry_path = os.path.join(self._prefix, path, entry)
        try:
            data, _ = self._zk.get(entry_path)
        except NoNodeError:
            raise ZkNoNodeError
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss
        return data

    @timed(1.0)
    def get_shared_metadata(self, path):
        """
        List the items in a shared path, with their relevant metadata.

        :param str path: The path the items are stored in. E.g.  ``"shared"``
        :return: A list of (item, metadata) for each item in the shared path.
        :rtype: [(:class:`bytes`, :class:`ZnodeStat`),...] or ``[]``
        :raises:
            ZkConnectionLoss: if the connection to ZK drops
        """
        if not self.is_connected():
            return []
        path = os.path.join(self._prefix, path)
        entry_meta = []
        try:
            entries = self._zk.get_children(path)
            for entry in entries:
                entry_path = os.path.join(path, entry)
                meta = self._zk.exists(entry_path)
                if meta:
                    entry_meta.append((entry, meta))
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss
        return entry_meta

    @timed(1.0)
    def expire_shared_items(self, path, cutoff):
        """
        Delete items from a shared path that haven't been modified since
        `cutoff`

        :param str path: The path the items are stored in. E.g.  ``"shared"``
        :param int cutoff: Time (in seconds since epoch) before which to expire
                           items.
        :return: Number of items expired
        :rtype: int
        :raises:
            ZkConnectionLoss: if the connection to ZK drops
            ZkNoNodeError: if a node disappears unexpectedly
        """
        if not self.is_connected():
            return
        entries_meta = self.get_shared_metadata(path)
        if not entries_meta:
            return 0
        count = 0
        for entry, meta in entries_meta:
            if meta.last_modified < cutoff:
                count += 1
                try:
                    self._zk.delete(os.path.join(self._prefix, path, entry))
                except NoNodeError:
                    # This shouldn't happen, so raise an exception if it does.
                    raise ZkNoNodeError
                except (ConnectionLoss, SessionExpiredError):
                    raise ZkConnectionLoss
        return count

    def retry(self, desc, f, *args, _retries=4, _timeout=10.0, **kwargs):
        """
        Execute a given operation, retrying it if fails due to connection
        problems.

        :param str desc: Description of the operation
        :param function f: Function to call, passing in \*args and \*\*kwargs
        :param int _retries: Number of times to retry the operation, or `None`
                             to retry indefinitely.
        :param float _timeout: Number of seconds to wait for a connection, or
                               `None` to wait indefinitely.
        """
        count = -1
        while True:
            count += 1
            if _retries is not None and count > _retries:
                break
            if not self.wait_connected(timeout=_timeout):
                logging.warning("%s: No connection to ZK", desc)
                continue
            try:
                return f(*args, **kwargs)
            except ZkConnectionLoss:
                logging.warning("%s: Connection to ZK dropped", desc)
        raise ZkRetryLimit("%s: Failed %s times, giving up" %
                           (desc, 1+_retries))

    def run_shared_cache_handling(self):
        for shared_cache in self._shared_caches:
            shared_cache.run()


class ZkParty(object):
    """
    A wrapper for a `Kazoo Party
    <https://kazoo.readthedocs.org/en/latest/api/recipe/party.html>`_.
    """
    def __init__(self, zk, path, id_, autojoin_):
        """
        :param zk: A kazoo instance
        :param str path: The absolute path of the party
        :param str id_: The service id value to use in the party
        :param bool autojoin_: Join the party automatically
        """
        self._autojoin = autojoin_
        self._path = path
        try:
            self._party = zk.Party(path, id_)
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss
        self.autojoin()

    def join(self):
        try:
            self._party.join()
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss

    def autojoin(self):
        """
        If the autojoin parameter was set to True, join the party.
        """
        if self._autojoin:
            self.join()
        entries = self.list()
        names = set([entry.split("\0")[0] for entry in entries])
        logging.debug("Current party (%s) members are: %s", self._path,
                      sorted(names))

    def list(self):
        """
        List the current party member IDs
        """
        try:
            return set(self._party)
        except (ConnectionLoss, SessionExpiredError):
            raise ZkConnectionLoss


class ZkSharedCache(object):
    """
    Class for handling ZK's shared path.
    """
    def __init__(self, zk, path, handler, state_synced):
        """
        :param Zookeeper zk: A Zookeeper instance
        :param str path: The absolute path of the cache
        :param function handler: Handler for a list of cached objects
        :param threading.Event state_synced: state for synchronization
        """
        self.zk = zk
        self.path = path
        self.handler = handler
        self._state_synced = state_synced
        self._latest_entry = 0

    def handle_shared_entries(self):
        """
        A thread to handle Zookeeper connects/disconnects and the shared cache
        of entries.

        On connect, it registers us as in-service, and loads the shared cache
        of entries from ZK.

        While connected, it calls _read_cached_entries() to read updated entries
        from the cache.
        """
        while True:
            if not self.zk.is_connected():
                self._state_synced.clear()
                self.zk.wait_connected()
            else:
                time.sleep(0.5)
                if not self._state_synced.is_set():
                    # Make sure we re-read the entire cache
                    self._latest_entry = 0
            count = None
            try:
                count = self._read_cached_entries()
            except ZkConnectionLoss:
                self._state_synced.clear()
                continue
            if count:
                logging.debug("Processed %d new/updated entries from %s",
                              count, self.path)
            self._state_synced.set()

    def _read_cached_entries(self):
        """
        Read new/updated entries from the shared cache and send them for
        processing.
        """
        desc = "Fetching list of entries from shared path: %s" % self.path
        entries_meta = self.zk.get_shared_metadata(
            self.path,
            timed_desc=desc)
        if not entries_meta:
            return 0
        new = []
        newest = self._latest_entry
        for entry, meta in entries_meta:
            if meta.last_modified > self._latest_entry:
                new.append(entry)
            if meta.last_modified > newest:
                newest = meta.last_modified
        self._latest_entry = newest
        desc = "Processing %s new entries from shared path: %s" % (len(new),
                                                                   self.path)
        count = self._process_cached_entries(new, timed_desc=desc)
        return count

    @timed(1.0)
    def _process_cached_entries(self, entries):
        """
        Retrieve new beacons from the shared cache and send them for local
        processing (self.handler).

        :param entries: cached path segments.
        :param entries: list
        """
        new_entries = []
        for entry in entries:
            try:
                raw = self.zk.get_shared_item(self.path, entry)
            except ZkConnectionLoss:
                logging.warning("Unable to retrieve entry from shared "
                                "path %s: no connection to ZK" % self.path)
                break
            except ZkNoNodeError:
                logging.debug("Unable to retrieve entry from shared cache: "
                              "no such entry (%s/%s)" %
                              (self.path, entry))
                continue
            new_entries.append(raw)
        self.handler(new_entries)
        return len(new_entries)

    def run(self):
        """
        Run thread that handles shared path.
        """
        threading.Thread(
            target=thread_safety_net, args=(self.handle_shared_entries,),
            name="handle_shared_entries(%s)" % self.path, daemon=True).start()
