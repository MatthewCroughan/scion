// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/pkg/discovery (interfaces: TopologyInformation)

// Package mock_discovery is a generated GoMock package.
package mock_discovery

import (
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	topology "github.com/scionproto/scion/go/lib/topology"
)

// MockTopologyInformation is a mock of TopologyInformation interface.
type MockTopologyInformation struct {
	ctrl     *gomock.Controller
	recorder *MockTopologyInformationMockRecorder
}

// MockTopologyInformationMockRecorder is the mock recorder for MockTopologyInformation.
type MockTopologyInformationMockRecorder struct {
	mock *MockTopologyInformation
}

// NewMockTopologyInformation creates a new mock instance.
func NewMockTopologyInformation(ctrl *gomock.Controller) *MockTopologyInformation {
	mock := &MockTopologyInformation{ctrl: ctrl}
	mock.recorder = &MockTopologyInformationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopologyInformation) EXPECT() *MockTopologyInformationMockRecorder {
	return m.recorder
}

// ColibriServices mocks base method.
func (m *MockTopologyInformation) ColibriServices() ([]*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ColibriServices")
	ret0, _ := ret[0].([]*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ColibriServices indicates an expected call of ColibriServices.
func (mr *MockTopologyInformationMockRecorder) ColibriServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ColibriServices", reflect.TypeOf((*MockTopologyInformation)(nil).ColibriServices))
}

// Gateways mocks base method.
func (m *MockTopologyInformation) Gateways() ([]topology.GatewayInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gateways")
	ret0, _ := ret[0].([]topology.GatewayInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Gateways indicates an expected call of Gateways.
func (mr *MockTopologyInformationMockRecorder) Gateways() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gateways", reflect.TypeOf((*MockTopologyInformation)(nil).Gateways))
}

// HiddenSegmentLookupAddresses mocks base method.
func (m *MockTopologyInformation) HiddenSegmentLookupAddresses() ([]*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HiddenSegmentLookupAddresses")
	ret0, _ := ret[0].([]*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HiddenSegmentLookupAddresses indicates an expected call of HiddenSegmentLookupAddresses.
func (mr *MockTopologyInformationMockRecorder) HiddenSegmentLookupAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HiddenSegmentLookupAddresses", reflect.TypeOf((*MockTopologyInformation)(nil).HiddenSegmentLookupAddresses))
}

// HiddenSegmentRegistrationAddresses mocks base method.
func (m *MockTopologyInformation) HiddenSegmentRegistrationAddresses() ([]*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HiddenSegmentRegistrationAddresses")
	ret0, _ := ret[0].([]*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HiddenSegmentRegistrationAddresses indicates an expected call of HiddenSegmentRegistrationAddresses.
func (mr *MockTopologyInformationMockRecorder) HiddenSegmentRegistrationAddresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HiddenSegmentRegistrationAddresses", reflect.TypeOf((*MockTopologyInformation)(nil).HiddenSegmentRegistrationAddresses))
}
