// Copyright 2021 ETH Zurich
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

package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/peer"

	base "github.com/scionproto/scion/go/co/reservation"
	"github.com/scionproto/scion/go/co/reservation/e2e"
	"github.com/scionproto/scion/go/co/reservation/translate"
	"github.com/scionproto/scion/go/co/reservationstorage"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/colibri"
	"github.com/scionproto/scion/go/lib/colibri/coliquic"
	"github.com/scionproto/scion/go/lib/colibri/reservation"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/util"
	colpb "github.com/scionproto/scion/go/pkg/proto/colibri"
)

type ColibriService struct {
	Store reservationstorage.Store
}

var _ colpb.ColibriServer = (*ColibriService)(nil)

func (s *ColibriService) SetupSegment(ctx context.Context, msg *colpb.SegmentSetupRequest) (
	*colpb.SegmentSetupResponse, error) {

	msg.Base.Path.CurrentStep++
	// path, err := extractPath(ctx)
	// if err != nil {
	// 	log.Error("setup segment", "err", err)
	// 	return nil, err
	// }
	req, err := translate.SetupReq(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		// should send a message?
		return nil, err
	}
	res, err := s.Store.AdmitSegmentReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		// should send a message?
		return nil, err
	}
	pbRes := translate.PBufSetupResponse(res)
	return pbRes, nil
}

func (s *ColibriService) ConfirmSegmentIndex(ctx context.Context, msg *colpb.Request) (
	*colpb.Response, error) {

	msg.Path.CurrentStep++
	req, err := translate.Request(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		return nil, err
	}
	res, err := s.Store.ConfirmSegmentReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		return nil, err
	}
	pbRes := translate.PBufResponse(res)

	return pbRes, nil
}

func (s *ColibriService) ActivateSegmentIndex(ctx context.Context, msg *colpb.Request) (
	*colpb.Response, error) {

	msg.Path.CurrentStep++
	req, err := translate.Request(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		return nil, err
	}
	res, err := s.Store.ActivateSegmentReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		return nil, err
	}
	pbRes := translate.PBufResponse(res)

	return pbRes, nil
}

func (s *ColibriService) TeardownSegment(ctx context.Context, msg *colpb.Request) (
	*colpb.Response, error) {

	msg.Path.CurrentStep++
	req, err := translate.Request(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		return nil, err
	}
	res, err := s.Store.TearDownSegmentReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		return nil, err
	}
	pbRes := translate.PBufResponse(res)

	return pbRes, nil
}

func (s *ColibriService) CleanupSegmentIndex(ctx context.Context, msg *colpb.Request) (
	*colpb.Response, error) {

	msg.Path.CurrentStep++
	req, err := translate.Request(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		return nil, err
	}
	res, err := s.Store.CleanupSegmentReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		return nil, err
	}
	pbRes := translate.PBufResponse(res)

	return pbRes, nil
}

func (s *ColibriService) ListReservations(ctx context.Context, msg *colpb.ListRequest) (
	*colpb.ListResponse, error) {

	dstIA := addr.IAInt(msg.DstIa).IA()
	looks, err := s.Store.ListReservations(ctx, dstIA, reservation.PathType(msg.PathType))
	if err != nil {
		log.Error("colibri store while listing rsvs", "err", err)
		return &colpb.ListResponse{
			ErrorMessage: err.Error(),
		}, nil
	}
	return translate.PBufListResponse(looks), nil
}

func (s *ColibriService) SetupE2E(ctx context.Context, msg *colpb.E2ESetupRequest) (
	*colpb.E2ESetupResponse, error) {

	msg.Base.Path.CurrentStep++
	req, err := translate.E2ESetupRequest(msg)
	if err != nil {
		log.Error("translating e2e setup", "err", err)
		return nil, serrors.WrapStr("translating e2e setup", err)
	}
	res, err := s.Store.AdmitE2EReservation(ctx, req)
	if err != nil {
		log.Error("admitting e2e", "err", err)
		return nil, err
	}
	return translate.PBufE2ESetupResponse(res), nil
}

func (s *ColibriService) CleanupE2EIndex(ctx context.Context, msg *colpb.Request) (
	*colpb.Response, error) {

	msg.Path.CurrentStep++
	req, err := translate.Request(msg)
	if err != nil {
		log.Error("error unmarshalling", "err", err)
		return nil, err
	}
	res, err := s.Store.CleanupE2EReservation(ctx, req)
	if err != nil {
		log.Error("colibri store returned an error", "err", err)
		return nil, err
	}
	pbRes := translate.PBufResponse(res)

	return pbRes, nil
}

func (s *ColibriService) ListStitchables(ctx context.Context, msg *colpb.ListStitchablesRequest) (
	*colpb.ListStitchablesResponse, error) {

	if _, err := checkLocalCaller(ctx); err != nil {
		return nil, err
	}

	dstIA := addr.IAInt(msg.DstIa).IA()
	stitchables, err := s.Store.ListStitchableSegments(ctx, dstIA)
	if err != nil {
		log.Error("colibri store while listing stitchables", "err", err)
		return &colpb.ListStitchablesResponse{
			ErrorMessage: err.Error(),
		}, nil
	}
	return translate.PBufStitchableResponse(stitchables), nil
}

// SetupReservation serves the intra AS clients, setting up or renewing an E2E reservation.
func (s *ColibriService) SetupReservation(ctx context.Context, msg *colpb.DaemonSetupRequest) (
	*colpb.DaemonSetupResponse, error) {

	clientAddr, err := checkLocalCaller(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	// build a valid E2E setup request now and query the store with it
	pbReq := &colpb.E2ESetupRequest{
		Base: &colpb.Request{
			Id:        msg.Id,
			Index:     msg.Index,
			Timestamp: util.TimeToSecs(now),
			Path:      &colpb.TransparentPath{},
		},
		RequestedBw: msg.RequestedBw,
		Params: &colpb.E2ESetupRequest_PathParams{
			Segments:       msg.Segments,
			CurrentSegment: 0,
			SrcIa:          msg.SrcIa,
			SrcHost:        clientAddr.IP,
			DstIa:          msg.DstIa,
			DstHost:        msg.DstHost,
		},
		Allocationtrail: nil,
	}
	req, err := translate.E2ESetupRequest(pbReq)
	if err != nil {
		log.Error("translating initial E2E setup from daemon to service", "err", err)
		return nil, err
	}

	res, err := s.Store.AdmitE2EReservation(ctx, req)
	if err != nil {
		log.Error("colibri store setting up an e2e reservation", "err", err)
		var trail []uint32
		var failedStep uint32
		if failure, ok := res.(*e2e.SetupResponseFailure); ok {
			trail = make([]uint32, len(failure.AllocTrail))
			for i, b := range failure.AllocTrail {
				trail[i] = uint32(b)
			}
			failedStep = uint32(failure.FailedStep)
		}
		return &colpb.DaemonSetupResponse{
			Failure: &colpb.DaemonSetupResponse_Failure{
				ErrorMessage: err.Error(),
				FailedStep:   failedStep,
				AllocTrail:   trail,
			},
		}, nil
	}
	pbMsg := &colpb.DaemonSetupResponse{}
	switch res := res.(type) {
	case *e2e.SetupResponseFailure:
		trail := make([]uint32, len(res.AllocTrail))
		for i, b := range res.AllocTrail {
			trail[i] = uint32(b)
		}
		pbMsg.Failure = &colpb.DaemonSetupResponse_Failure{
			ErrorMessage: res.Message,
			FailedStep:   uint32(res.FailedStep),
			AllocTrail:   trail,
		}
	case *e2e.SetupResponseSuccess:
		token, err := reservation.TokenFromRaw(res.Token)
		if err != nil {
			return nil, serrors.WrapStr("decoding token in colibri service", err)
		}
		path := e2e.DeriveColibriPath(&req.ID, token)
		egressId := ""
		if len(path.HopFields) > 0 {
			egressId = fmt.Sprintf("%d", path.HopFields[0].EgressId)
		}
		rawPath := make([]byte, path.Len())
		err = path.SerializeTo(rawPath)
		if err != nil {
			return nil, serrors.WrapStr("serializing a colibri path in colibri service", err)
		}
		// nexthop holds the interface id until the daemon resolves it with the topology
		pbMsg.Success = &colpb.DaemonSetupResponse_Success{
			Spath:   rawPath,
			NextHop: egressId,
		}
	}
	return pbMsg, nil
}

// CleanupReservation serves the intra AS clients, cleaning an E2E reservation.
func (s *ColibriService) CleanupReservation(ctx context.Context, msg *colpb.DaemonCleanupRequest) (
	*colpb.DaemonCleanupResponse, error) {

	if _, err := checkLocalCaller(ctx); err != nil {
		return nil, err
	}
	req := &base.Request{
		MsgId: base.MsgId{
			ID:        *translate.ID(msg.Id),
			Index:     reservation.IndexNumber(msg.Index),
			Timestamp: time.Now(),
		},
		Path: &base.TransparentPath{},
	}
	res, err := s.Store.CleanupE2EReservation(ctx, req)
	if err != nil {
		var failedStep uint32
		if failure, ok := res.(*base.ResponseFailure); ok {
			failedStep = uint32(failure.FailedStep)
		}
		return &colpb.DaemonCleanupResponse{
			Failure: &colpb.DaemonCleanupResponse_Failure{
				ErrorMessage: err.Error(),
				FailedStep:   uint32(failedStep),
			},
		}, nil
	}
	return &colpb.DaemonCleanupResponse{}, nil
}

func (s *ColibriService) AddAdmissionEntry(ctx context.Context,
	req *colpb.DaemonAdmissionEntry) (*colpb.DaemonAdmissionEntryResponse, error) {

	clientAddr, err := checkLocalCaller(ctx)
	if err != nil {
		return nil, err
	}
	// TODO(juagargi)
	// because we can't guarantee that the IP the client requested is reachable from this
	// service, checking that the connection from the endhost to this service uses the same
	// IP is wrong.
	// A new design for this check must be created and implemented. For now, the check is
	// completely disabled (commented code below).
	// if len(req.DstHost) > 0 {
	// 	// check that we have the same IP address in the DstHost field and the TCP connection
	// 	if !bytes.Equal(req.DstHost, clientAddr.IP) {
	// 		return nil, serrors.New("IP address in request not the same as connnection",
	// 			"req", net.IP(req.DstHost).String(), "conn", clientAddr.IP.String())
	// 	}
	// }
	if len(req.DstHost) == 0 {
		req.DstHost = clientAddr.IP
	}
	entry := &colibri.AdmissionEntry{
		DstHost:         req.DstHost,
		ValidUntil:      util.SecsToTime(req.ValidUntil),
		RegexpIA:        req.RegexpIa,
		RegexpHost:      req.RegexpHost,
		AcceptAdmission: req.Accept,
	}
	validUntil, err := s.Store.AddAdmissionEntry(ctx, entry)
	return &colpb.DaemonAdmissionEntryResponse{
		ValidUntil: util.TimeToSecs(validUntil),
	}, err
}

// extractPath returns the PacketPath, ingress and egress used with this RPC.
func extractPath(ctx context.Context) (base.PacketPath, error) {
	// TODO(juagargi) move from PacketPath to TransparentPath
	// TODO(juagargi) call this function to check that the transport path matches that
	// of base.Request.Path if the transport path is of colibri type.
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return nil, serrors.New("no peer found")
	}
	raddr, ok := p.Addr.(*snet.UDPAddr)
	if !ok || raddr == nil {
		return nil, serrors.New("no valid scion address found", "addr", p.Addr)
	}
	path, err := base.NewPacketPath(raddr.Path)
	if err != nil {
		return path, err
	}
	usage, ok, err := coliquic.UsageFromContext(ctx)
	_, _, _ = usage, ok, err
	return path, err
}

// checkLocalCaller prevents the service from doing anything if the caller is not from the local AS.
// We do it by checking the peer. We could instantiate the local ColibriService differently.
func checkLocalCaller(ctx context.Context) (*net.TCPAddr, error) {
	p, ok := peer.FromContext(ctx)
	if !ok || p == nil {
		return nil, serrors.New("no peer found")
	}
	tcpaddr, ok := p.Addr.(*net.TCPAddr)
	if !ok || tcpaddr == nil {
		return nil, serrors.New("no valid local tcp address found", "addr", p.Addr,
			"type", common.TypeOf(p.Addr))
	}
	return tcpaddr, nil
}