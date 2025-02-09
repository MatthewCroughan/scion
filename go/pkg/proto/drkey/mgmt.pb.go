// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.3
// source: proto/drkey/mgmt/v1/mgmt.proto

package drkey

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DRKeyLvl1Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValTime   *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=val_time,json=valTime,proto3" json:"val_time,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *DRKeyLvl1Request) Reset() {
	*x = DRKeyLvl1Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DRKeyLvl1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DRKeyLvl1Request) ProtoMessage() {}

func (x *DRKeyLvl1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DRKeyLvl1Request.ProtoReflect.Descriptor instead.
func (*DRKeyLvl1Request) Descriptor() ([]byte, []int) {
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP(), []int{0}
}

func (x *DRKeyLvl1Request) GetValTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ValTime
	}
	return nil
}

func (x *DRKeyLvl1Request) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type DRKeyLvl1Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EpochBegin *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=epoch_begin,json=epochBegin,proto3" json:"epoch_begin,omitempty"`
	EpochEnd   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=epoch_end,json=epochEnd,proto3" json:"epoch_end,omitempty"`
	Drkey      []byte                 `protobuf:"bytes,3,opt,name=drkey,proto3" json:"drkey,omitempty"`
	Timestamp  *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *DRKeyLvl1Response) Reset() {
	*x = DRKeyLvl1Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DRKeyLvl1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DRKeyLvl1Response) ProtoMessage() {}

func (x *DRKeyLvl1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DRKeyLvl1Response.ProtoReflect.Descriptor instead.
func (*DRKeyLvl1Response) Descriptor() ([]byte, []int) {
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP(), []int{1}
}

func (x *DRKeyLvl1Response) GetEpochBegin() *timestamppb.Timestamp {
	if x != nil {
		return x.EpochBegin
	}
	return nil
}

func (x *DRKeyLvl1Response) GetEpochEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.EpochEnd
	}
	return nil
}

func (x *DRKeyLvl1Response) GetDrkey() []byte {
	if x != nil {
		return x.Drkey
	}
	return nil
}

func (x *DRKeyLvl1Response) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type DRKeyLvl2Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Protocol string                      `protobuf:"bytes,1,opt,name=protocol,proto3" json:"protocol,omitempty"`
	ReqType  uint32                      `protobuf:"varint,2,opt,name=req_type,json=reqType,proto3" json:"req_type,omitempty"`
	ValTime  *timestamppb.Timestamp      `protobuf:"bytes,3,opt,name=val_time,json=valTime,proto3" json:"val_time,omitempty"`
	SrcIa    uint64                      `protobuf:"varint,4,opt,name=src_ia,json=srcIa,proto3" json:"src_ia,omitempty"`
	DstIa    uint64                      `protobuf:"varint,5,opt,name=dst_ia,json=dstIa,proto3" json:"dst_ia,omitempty"`
	SrcHost  *DRKeyLvl2Request_DRKeyHost `protobuf:"bytes,6,opt,name=src_host,json=srcHost,proto3" json:"src_host,omitempty"`
	DstHost  *DRKeyLvl2Request_DRKeyHost `protobuf:"bytes,7,opt,name=dst_host,json=dstHost,proto3" json:"dst_host,omitempty"`
	Misc     []byte                      `protobuf:"bytes,8,opt,name=misc,proto3" json:"misc,omitempty"`
}

func (x *DRKeyLvl2Request) Reset() {
	*x = DRKeyLvl2Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DRKeyLvl2Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DRKeyLvl2Request) ProtoMessage() {}

func (x *DRKeyLvl2Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DRKeyLvl2Request.ProtoReflect.Descriptor instead.
func (*DRKeyLvl2Request) Descriptor() ([]byte, []int) {
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP(), []int{2}
}

func (x *DRKeyLvl2Request) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *DRKeyLvl2Request) GetReqType() uint32 {
	if x != nil {
		return x.ReqType
	}
	return 0
}

func (x *DRKeyLvl2Request) GetValTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ValTime
	}
	return nil
}

func (x *DRKeyLvl2Request) GetSrcIa() uint64 {
	if x != nil {
		return x.SrcIa
	}
	return 0
}

func (x *DRKeyLvl2Request) GetDstIa() uint64 {
	if x != nil {
		return x.DstIa
	}
	return 0
}

func (x *DRKeyLvl2Request) GetSrcHost() *DRKeyLvl2Request_DRKeyHost {
	if x != nil {
		return x.SrcHost
	}
	return nil
}

func (x *DRKeyLvl2Request) GetDstHost() *DRKeyLvl2Request_DRKeyHost {
	if x != nil {
		return x.DstHost
	}
	return nil
}

func (x *DRKeyLvl2Request) GetMisc() []byte {
	if x != nil {
		return x.Misc
	}
	return nil
}

type DRKeyLvl2Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp  *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Drkey      []byte                 `protobuf:"bytes,2,opt,name=drkey,proto3" json:"drkey,omitempty"`
	EpochBegin *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=epoch_begin,json=epochBegin,proto3" json:"epoch_begin,omitempty"`
	EpochEnd   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=epoch_end,json=epochEnd,proto3" json:"epoch_end,omitempty"`
	Misc       []byte                 `protobuf:"bytes,5,opt,name=misc,proto3" json:"misc,omitempty"`
}

func (x *DRKeyLvl2Response) Reset() {
	*x = DRKeyLvl2Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DRKeyLvl2Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DRKeyLvl2Response) ProtoMessage() {}

func (x *DRKeyLvl2Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DRKeyLvl2Response.ProtoReflect.Descriptor instead.
func (*DRKeyLvl2Response) Descriptor() ([]byte, []int) {
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP(), []int{3}
}

func (x *DRKeyLvl2Response) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *DRKeyLvl2Response) GetDrkey() []byte {
	if x != nil {
		return x.Drkey
	}
	return nil
}

func (x *DRKeyLvl2Response) GetEpochBegin() *timestamppb.Timestamp {
	if x != nil {
		return x.EpochBegin
	}
	return nil
}

func (x *DRKeyLvl2Response) GetEpochEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.EpochEnd
	}
	return nil
}

func (x *DRKeyLvl2Response) GetMisc() []byte {
	if x != nil {
		return x.Misc
	}
	return nil
}

type DRKeyLvl2Request_DRKeyHost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type uint32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Host []byte `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *DRKeyLvl2Request_DRKeyHost) Reset() {
	*x = DRKeyLvl2Request_DRKeyHost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DRKeyLvl2Request_DRKeyHost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DRKeyLvl2Request_DRKeyHost) ProtoMessage() {}

func (x *DRKeyLvl2Request_DRKeyHost) ProtoReflect() protoreflect.Message {
	mi := &file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DRKeyLvl2Request_DRKeyHost.ProtoReflect.Descriptor instead.
func (*DRKeyLvl2Request_DRKeyHost) Descriptor() ([]byte, []int) {
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP(), []int{2, 0}
}

func (x *DRKeyLvl2Request_DRKeyHost) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *DRKeyLvl2Request_DRKeyHost) GetHost() []byte {
	if x != nil {
		return x.Host
	}
	return nil
}

var File_proto_drkey_mgmt_v1_mgmt_proto protoreflect.FileDescriptor

var file_proto_drkey_mgmt_v1_mgmt_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x72, 0x6b, 0x65, 0x79, 0x2f, 0x6d, 0x67,
	0x6d, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x64, 0x72, 0x6b, 0x65, 0x79, 0x2e, 0x6d, 0x67,
	0x6d, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83, 0x01, 0x0a, 0x10, 0x44, 0x52, 0x4b, 0x65, 0x79,
	0x4c, 0x76, 0x6c, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x76,
	0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x76, 0x61, 0x6c, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xd9, 0x01, 0x0a,
	0x11, 0x44, 0x52, 0x4b, 0x65, 0x79, 0x4c, 0x76, 0x6c, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x5f, 0x62, 0x65, 0x67, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0a, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x12,
	0x37, 0x0a, 0x09, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x45, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x72, 0x6b, 0x65,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x64, 0x72, 0x6b, 0x65, 0x79, 0x12, 0x38,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x8f, 0x03, 0x0a, 0x10, 0x44, 0x52, 0x4b,
	0x65, 0x79, 0x4c, 0x76, 0x6c, 0x32, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x71,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x72, 0x65, 0x71,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x07, 0x76, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x73,
	0x72, 0x63, 0x5f, 0x69, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x73, 0x72, 0x63,
	0x49, 0x61, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x73, 0x74, 0x5f, 0x69, 0x61, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x64, 0x73, 0x74, 0x49, 0x61, 0x12, 0x4a, 0x0a, 0x08, 0x73, 0x72, 0x63,
	0x5f, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x64, 0x72, 0x6b, 0x65, 0x79, 0x2e, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x52, 0x4b, 0x65, 0x79, 0x4c, 0x76, 0x6c, 0x32, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x44, 0x52, 0x4b, 0x65, 0x79, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x07, 0x73, 0x72,
	0x63, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x4a, 0x0a, 0x08, 0x64, 0x73, 0x74, 0x5f, 0x68, 0x6f, 0x73,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x64, 0x72, 0x6b, 0x65, 0x79, 0x2e, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x52,
	0x4b, 0x65, 0x79, 0x4c, 0x76, 0x6c, 0x32, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x44,
	0x52, 0x4b, 0x65, 0x79, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x07, 0x64, 0x73, 0x74, 0x48, 0x6f, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x69, 0x73, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x6d, 0x69, 0x73, 0x63, 0x1a, 0x33, 0x0a, 0x09, 0x44, 0x52, 0x4b, 0x65, 0x79, 0x48, 0x6f,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0xed, 0x01, 0x0a, 0x11, 0x44,
	0x52, 0x4b, 0x65, 0x79, 0x4c, 0x76, 0x6c, 0x32, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x72,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x64, 0x72, 0x6b, 0x65, 0x79,
	0x12, 0x3b, 0x0a, 0x0b, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x5f, 0x62, 0x65, 0x67, 0x69, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x12, 0x37, 0x0a,
	0x09, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x65, 0x70,
	0x6f, 0x63, 0x68, 0x45, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x69, 0x73, 0x63, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x6d, 0x69, 0x73, 0x63, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x63, 0x69, 0x6f, 0x6e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x63, 0x69, 0x6f, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x72, 0x6b, 0x65, 0x79, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_drkey_mgmt_v1_mgmt_proto_rawDescOnce sync.Once
	file_proto_drkey_mgmt_v1_mgmt_proto_rawDescData = file_proto_drkey_mgmt_v1_mgmt_proto_rawDesc
)

func file_proto_drkey_mgmt_v1_mgmt_proto_rawDescGZIP() []byte {
	file_proto_drkey_mgmt_v1_mgmt_proto_rawDescOnce.Do(func() {
		file_proto_drkey_mgmt_v1_mgmt_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_drkey_mgmt_v1_mgmt_proto_rawDescData)
	})
	return file_proto_drkey_mgmt_v1_mgmt_proto_rawDescData
}

var file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_drkey_mgmt_v1_mgmt_proto_goTypes = []interface{}{
	(*DRKeyLvl1Request)(nil),           // 0: proto.drkey.mgmt.v1.DRKeyLvl1Request
	(*DRKeyLvl1Response)(nil),          // 1: proto.drkey.mgmt.v1.DRKeyLvl1Response
	(*DRKeyLvl2Request)(nil),           // 2: proto.drkey.mgmt.v1.DRKeyLvl2Request
	(*DRKeyLvl2Response)(nil),          // 3: proto.drkey.mgmt.v1.DRKeyLvl2Response
	(*DRKeyLvl2Request_DRKeyHost)(nil), // 4: proto.drkey.mgmt.v1.DRKeyLvl2Request.DRKeyHost
	(*timestamppb.Timestamp)(nil),      // 5: google.protobuf.Timestamp
}
var file_proto_drkey_mgmt_v1_mgmt_proto_depIdxs = []int32{
	5,  // 0: proto.drkey.mgmt.v1.DRKeyLvl1Request.val_time:type_name -> google.protobuf.Timestamp
	5,  // 1: proto.drkey.mgmt.v1.DRKeyLvl1Request.timestamp:type_name -> google.protobuf.Timestamp
	5,  // 2: proto.drkey.mgmt.v1.DRKeyLvl1Response.epoch_begin:type_name -> google.protobuf.Timestamp
	5,  // 3: proto.drkey.mgmt.v1.DRKeyLvl1Response.epoch_end:type_name -> google.protobuf.Timestamp
	5,  // 4: proto.drkey.mgmt.v1.DRKeyLvl1Response.timestamp:type_name -> google.protobuf.Timestamp
	5,  // 5: proto.drkey.mgmt.v1.DRKeyLvl2Request.val_time:type_name -> google.protobuf.Timestamp
	4,  // 6: proto.drkey.mgmt.v1.DRKeyLvl2Request.src_host:type_name -> proto.drkey.mgmt.v1.DRKeyLvl2Request.DRKeyHost
	4,  // 7: proto.drkey.mgmt.v1.DRKeyLvl2Request.dst_host:type_name -> proto.drkey.mgmt.v1.DRKeyLvl2Request.DRKeyHost
	5,  // 8: proto.drkey.mgmt.v1.DRKeyLvl2Response.timestamp:type_name -> google.protobuf.Timestamp
	5,  // 9: proto.drkey.mgmt.v1.DRKeyLvl2Response.epoch_begin:type_name -> google.protobuf.Timestamp
	5,  // 10: proto.drkey.mgmt.v1.DRKeyLvl2Response.epoch_end:type_name -> google.protobuf.Timestamp
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_proto_drkey_mgmt_v1_mgmt_proto_init() }
func file_proto_drkey_mgmt_v1_mgmt_proto_init() {
	if File_proto_drkey_mgmt_v1_mgmt_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DRKeyLvl1Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DRKeyLvl1Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DRKeyLvl2Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DRKeyLvl2Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DRKeyLvl2Request_DRKeyHost); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_drkey_mgmt_v1_mgmt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_drkey_mgmt_v1_mgmt_proto_goTypes,
		DependencyIndexes: file_proto_drkey_mgmt_v1_mgmt_proto_depIdxs,
		MessageInfos:      file_proto_drkey_mgmt_v1_mgmt_proto_msgTypes,
	}.Build()
	File_proto_drkey_mgmt_v1_mgmt_proto = out.File
	file_proto_drkey_mgmt_v1_mgmt_proto_rawDesc = nil
	file_proto_drkey_mgmt_v1_mgmt_proto_goTypes = nil
	file_proto_drkey_mgmt_v1_mgmt_proto_depIdxs = nil
}
