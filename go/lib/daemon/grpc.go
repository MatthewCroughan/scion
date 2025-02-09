// Copyright 2020 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package daemon

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/scionproto/scion/go/co/reservation/translate"
	"github.com/scionproto/scion/go/lib/addr"
	col "github.com/scionproto/scion/go/lib/colibri"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/common"
	dkctrl "github.com/scionproto/scion/go/lib/ctrl/drkey"
	"github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
	"github.com/scionproto/scion/go/lib/drkey"
	"github.com/scionproto/scion/go/lib/serrors"
	colpath "github.com/scionproto/scion/go/lib/slayers/path/colibri"
	"github.com/scionproto/scion/go/lib/slayers/path/scion"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/snet/path"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/topology"
	"github.com/scionproto/scion/go/lib/util"
	libgrpc "github.com/scionproto/scion/go/pkg/grpc"
	colpb "github.com/scionproto/scion/go/pkg/proto/colibri"
	sdpb "github.com/scionproto/scion/go/pkg/proto/daemon"
)

// Service exposes the API to connect to a SCION daemon service.
type Service struct {
	// Address is the address of the SCION daemon to connect to.
	Address string
	// Metrics are the metric counters that should be incremented when using the
	// connector.
	Metrics Metrics
}

func (s Service) Connect(ctx context.Context) (Connector, error) {
	a, err := net.ResolveTCPAddr("tcp", s.Address)
	if err != nil {
		s.Metrics.incConnects(err)
		return nil, serrors.WrapStr("resolving addr", err)
	}
	conn, err := libgrpc.SimpleDialer{}.Dial(ctx, a)
	if err != nil {
		s.Metrics.incConnects(err)
		return nil, serrors.WrapStr("dialing", err)
	}
	s.Metrics.incConnects(nil)
	return grpcConn{conn: conn, metrics: s.Metrics}, nil
}

type grpcConn struct {
	conn    *grpc.ClientConn
	metrics Metrics
}

func (c grpcConn) LocalIA(ctx context.Context) (addr.IA, error) {
	asInfo, err := c.ASInfo(ctx, addr.IA{})
	if err != nil {
		return addr.IA{}, err
	}
	ia := asInfo.IA
	return ia, nil
}

func (c grpcConn) Paths(ctx context.Context, dst, src addr.IA,
	f PathReqFlags) ([]snet.Path, error) {

	client := sdpb.NewDaemonServiceClient(c.conn)
	response, err := client.Paths(ctx, &sdpb.PathsRequest{
		SourceIsdAs:      uint64(src.IAInt()),
		DestinationIsdAs: uint64(dst.IAInt()),
		Hidden:           f.Hidden,
		Refresh:          f.Refresh,
	})
	if err != nil {
		c.metrics.incPaths(err)
		return nil, err
	}
	paths, err := pathResponseToPaths(response.Paths, dst)
	c.metrics.incPaths(err)
	return paths, err
}

func (c grpcConn) ASInfo(ctx context.Context, ia addr.IA) (ASInfo, error) {
	client := sdpb.NewDaemonServiceClient(c.conn)
	response, err := client.AS(ctx, &sdpb.ASRequest{IsdAs: uint64(ia.IAInt())})
	if err != nil {
		c.metrics.incAS(err)
		return ASInfo{}, err
	}
	c.metrics.incAS(nil)
	return ASInfo{
		IA:  addr.IAInt(response.IsdAs).IA(),
		MTU: uint16(response.Mtu),
	}, nil
}

func (c grpcConn) IFInfo(ctx context.Context,
	_ []common.IFIDType) (map[common.IFIDType]*net.UDPAddr, error) {

	client := sdpb.NewDaemonServiceClient(c.conn)
	response, err := client.Interfaces(ctx, &sdpb.InterfacesRequest{})
	if err != nil {
		c.metrics.incInterface(err)
		return nil, err
	}
	result := make(map[common.IFIDType]*net.UDPAddr)
	for ifID, intf := range response.Interfaces {
		a, err := net.ResolveUDPAddr("udp", intf.Address.Address)
		if err != nil {
			c.metrics.incInterface(err)
			return nil, serrors.WrapStr("parsing reply", err, "raw_uri", intf.Address.Address)
		}
		result[common.IFIDType(ifID)] = a
	}
	c.metrics.incInterface(nil)
	return result, nil
}

func (c grpcConn) SVCInfo(ctx context.Context, _ []addr.HostSVC) (map[addr.HostSVC]string, error) {
	client := sdpb.NewDaemonServiceClient(c.conn)
	response, err := client.Services(ctx, &sdpb.ServicesRequest{})
	if err != nil {
		c.metrics.incServcies(err)
		return nil, err
	}
	result := make(map[addr.HostSVC]string)
	for st, si := range response.Services {
		svc := topoServiceTypeToSVCAddr(topology.ServiceTypeFromString(st))
		if svc == addr.SvcNone || len(si.Services) == 0 {
			continue
		}
		result[svc] = si.Services[0].Uri
	}
	c.metrics.incServcies(nil)
	return result, nil
}

func (c grpcConn) RevNotification(ctx context.Context, revInfo *path_mgmt.RevInfo) error {
	client := sdpb.NewDaemonServiceClient(c.conn)
	_, err := client.NotifyInterfaceDown(ctx, &sdpb.NotifyInterfaceDownRequest{
		Id:    uint64(revInfo.IfID),
		IsdAs: uint64(revInfo.RawIsdas),
	})
	c.metrics.incIfDown(err)
	return err

}

func (c grpcConn) DRKeyGetLvl2Key(ctx context.Context, meta drkey.Lvl2Meta,
	valTime time.Time) (drkey.Lvl2Key, error) {

	client := sdpb.NewDaemonServiceClient(c.conn)

	lvl2Req := dkctrl.NewLvl2ReqFromMeta(meta, valTime)
	pbLvl2Req, err := lvl2reqToProtoRequest(lvl2Req)
	if err != nil {
		return drkey.Lvl2Key{}, err
	}

	reply, err := client.DRKeyLvl2(ctx, pbLvl2Req)
	if err != nil {
		return drkey.Lvl2Key{}, err
	}

	lvl2Key, err := getLvl2KeyFromReply(reply, meta)
	if err != nil {
		return drkey.Lvl2Key{}, err
	}
	return lvl2Key, nil
}

func (c grpcConn) ColibriListRsvs(ctx context.Context, dstIA addr.IA) (
	*col.StitchableSegments, error) {

	req := &sdpb.ColibriListRsvsRequest{
		Base: &colpb.ListStitchablesRequest{
			DstIa: uint64(dstIA.IAInt()),
		},
	}
	client := sdpb.NewDaemonServiceClient(c.conn)
	sdRes, err := client.ColibriListRsvs(ctx, req)
	if err != nil {
		return nil, err
	}
	if sdRes.Base.ErrorMessage != "" {
		return nil, fmt.Errorf(sdRes.Base.ErrorMessage)
	}

	stitchable, err := translate.StitchableSegments(sdRes.Base)
	if err != nil {
		return nil, err
	}
	return stitchable, nil
}

func (c grpcConn) ColibriSetupRsv(ctx context.Context, req *col.E2EReservationSetup) (
	snet.Path, error) {

	pbSegs := make([]*colpb.ReservationID, len(req.Segments))
	for i, r := range req.Segments {
		pbSegs[i] = translate.PBufID(&r)
	}
	pbReq := &sdpb.ColibriSetupRsvRequest{
		Base: &colpb.SetupReservationRequest{
			Id:          translate.PBufID(&req.Id),
			SrcIa:       uint64(req.SrcIA.IAInt()),
			DstIa:       uint64(req.DstIA.IAInt()),
			DstHost:     req.DstHost,
			Index:       uint32(req.Index),
			RequestedBw: uint32(req.RequestedBW),
			Segments:    pbSegs,
		},
	}
	client := sdpb.NewDaemonServiceClient(c.conn)
	sdRes, err := client.ColibriSetupRsv(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	if sdRes.Base.Failure != nil {
		trail := make([]reservation.BWCls, len(sdRes.Base.Failure.AllocTrail))
		for i, b := range sdRes.Base.Failure.AllocTrail {
			trail[i] = reservation.BWCls(b)
		}
		return nil, &col.E2ESetupError{
			E2EResponseError: col.E2EResponseError{
				Message:  sdRes.Base.Failure.ErrorMessage,
				FailedAS: int(sdRes.Base.Failure.FailedStep),
			},
			AllocationTrail: trail,
		}
	}
	nextHop, err := net.ResolveUDPAddr("udp", sdRes.Base.Success.NextHop)
	if err != nil {
		return nil, serrors.WrapStr("parsing next hop", err)
	}
	return &path.Path{
		SPath: spath.Path{
			Raw:  sdRes.Base.Success.Spath,
			Type: colpath.PathType,
		},
		NextHop: nextHop,
	}, nil
}

func (c grpcConn) ColibriCleanupRsv(ctx context.Context, id *reservation.ID,
	idx reservation.IndexNumber) error {

	if id == nil {
		return serrors.New("invalid nil ID")
	}
	if !id.IsE2EID() {
		return serrors.New("this id is not for an E2E reservation")
	}
	pbReq := &sdpb.ColibriCleanupRsvRequest{
		Base: &colpb.CleanupReservationRequest{
			Id:    translate.PBufID(id),
			Index: uint32(idx),
		},
	}
	client := sdpb.NewDaemonServiceClient(c.conn)
	sdRes, err := client.ColibriCleanupRsv(ctx, pbReq)
	if err != nil {
		return err
	}
	if sdRes.Base.Failure != nil {
		return &col.E2EResponseError{
			Message:  sdRes.Base.Failure.ErrorMessage,
			FailedAS: int(sdRes.Base.Failure.FailedStep),
		}
	}
	return nil
}

func (c grpcConn) ColibriAddAdmissionEntry(ctx context.Context, entry *col.AdmissionEntry) (
	time.Time, error) {
	req := &sdpb.ColibriAddAdmissionEntryRequest{
		Base: &colpb.AddAdmissionEntryRequest{
			DstHost:    entry.DstHost,
			ValidUntil: util.TimeToSecs(entry.ValidUntil),
			RegexpIa:   entry.RegexpIA,
			RegexpHost: entry.RegexpHost,
			Accept:     entry.AcceptAdmission,
		},
	}
	client := sdpb.NewDaemonServiceClient(c.conn)
	res, err := client.ColibriAddAdmissionEntry(ctx, req)
	if err != nil {
		return time.Time{}, err
	}
	return util.SecsToTime(res.Base.ValidUntil), nil
}

func (c grpcConn) Close(_ context.Context) error {
	return c.conn.Close()
}

func pathResponseToPaths(paths []*sdpb.Path, dst addr.IA) ([]snet.Path, error) {
	result := make([]snet.Path, 0, len(paths))
	for _, p := range paths {
		cp, err := convertPath(p, dst)
		if err != nil {
			return nil, err
		}
		result = append(result, cp)
	}
	return result, nil
}

func convertPath(p *sdpb.Path, dst addr.IA) (path.Path, error) {
	expiry := time.Unix(p.Expiration.Seconds, int64(p.Expiration.Nanos))
	if len(p.Interfaces) == 0 {
		return path.Path{
			Dst: dst,
			Meta: snet.PathMetadata{
				MTU:    uint16(p.Mtu),
				Expiry: expiry,
			},
		}, nil
	}
	underlayA, err := net.ResolveUDPAddr("udp", p.Interface.Address.Address)
	if err != nil {
		return path.Path{}, serrors.WrapStr("resolving underlay", err)
	}
	interfaces := make([]snet.PathInterface, len(p.Interfaces))
	for i, pi := range p.Interfaces {
		interfaces[i] = snet.PathInterface{
			ID: common.IFIDType(pi.Id),
			IA: addr.IAInt(pi.IsdAs).IA(),
		}
	}
	latency := make([]time.Duration, len(p.Latency))
	for i, v := range p.Latency {
		latency[i] = time.Second*time.Duration(v.Seconds) + time.Duration(v.Nanos)
	}
	geo := make([]snet.GeoCoordinates, len(p.Geo))
	for i, v := range p.Geo {
		geo[i] = snet.GeoCoordinates{
			Latitude:  v.Latitude,
			Longitude: v.Longitude,
			Address:   v.Address,
		}
	}
	linkType := make([]snet.LinkType, len(p.LinkType))
	for i, v := range p.LinkType {
		linkType[i] = linkTypeFromPB(v)
	}

	return path.Path{
		Dst:     dst,
		SPath:   spath.Path{Raw: p.Raw, Type: scion.PathType},
		NextHop: underlayA,
		Meta: snet.PathMetadata{
			Interfaces:   interfaces,
			MTU:          uint16(p.Mtu),
			Expiry:       expiry,
			Latency:      latency,
			Bandwidth:    p.Bandwidth,
			Geo:          geo,
			LinkType:     linkType,
			InternalHops: p.InternalHops,
			Notes:        p.Notes,
		},
	}, nil
}

func linkTypeFromPB(lt sdpb.LinkType) snet.LinkType {
	switch lt {
	case sdpb.LinkType_LINK_TYPE_DIRECT:
		return snet.LinkTypeDirect
	case sdpb.LinkType_LINK_TYPE_MULTI_HOP:
		return snet.LinkTypeMultihop
	case sdpb.LinkType_LINK_TYPE_OPEN_NET:
		return snet.LinkTypeOpennet
	default:
		return snet.LinkTypeUnset
	}
}

func topoServiceTypeToSVCAddr(st topology.ServiceType) addr.HostSVC {
	switch st {
	case topology.Control:
		return addr.SvcCS
	default:
		return addr.SvcNone
	}
}

func lvl2reqToProtoRequest(req dkctrl.Lvl2Req) (*sdpb.DRKeyLvl2Request, error) {
	baseReq, err := dkctrl.Lvl2reqToProtoRequest(req)
	if err != nil {
		return nil, err
	}
	return &sdpb.DRKeyLvl2Request{
		BaseReq: baseReq,
	}, nil
}

// getLvl2KeyFromReply decrypts and extracts the level 1 drkey from the reply.
func getLvl2KeyFromReply(rep *sdpb.DRKeyLvl2Response, meta drkey.Lvl2Meta) (drkey.Lvl2Key, error) {
	return dkctrl.GetLvl2KeyFromReply(rep.BaseRep, meta)
}
