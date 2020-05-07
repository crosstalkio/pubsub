// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.4
// source: api/pubsub/pubsub.proto

package pubsub

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*Message_Control
	//	*Message_Data
	Type isMessage_Type `protobuf_oneof:"type"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{0}
}

func (m *Message) GetType() isMessage_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Message) GetControl() *Control {
	if x, ok := x.GetType().(*Message_Control); ok {
		return x.Control
	}
	return nil
}

func (x *Message) GetData() *Data {
	if x, ok := x.GetType().(*Message_Data); ok {
		return x.Data
	}
	return nil
}

type isMessage_Type interface {
	isMessage_Type()
}

type Message_Control struct {
	Control *Control `protobuf:"bytes,1,opt,name=control,proto3,oneof"`
}

type Message_Data struct {
	Data *Data `protobuf:"bytes,2,opt,name=data,proto3,oneof"`
}

func (*Message_Control) isMessage_Type() {}

func (*Message_Data) isMessage_Type() {}

type Control struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*Control_Request
	//	*Control_Response
	Type isControl_Type `protobuf_oneof:"type"`
}

func (x *Control) Reset() {
	*x = Control{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Control) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Control) ProtoMessage() {}

func (x *Control) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Control.ProtoReflect.Descriptor instead.
func (*Control) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{1}
}

func (m *Control) GetType() isControl_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Control) GetRequest() *Request {
	if x, ok := x.GetType().(*Control_Request); ok {
		return x.Request
	}
	return nil
}

func (x *Control) GetResponse() *Response {
	if x, ok := x.GetType().(*Control_Response); ok {
		return x.Response
	}
	return nil
}

type isControl_Type interface {
	isControl_Type()
}

type Control_Request struct {
	Request *Request `protobuf:"bytes,1,opt,name=request,proto3,oneof"`
}

type Control_Response struct {
	Response *Response `protobuf:"bytes,2,opt,name=response,proto3,oneof"`
}

func (*Control_Request) isControl_Type() {}

func (*Control_Response) isControl_Type() {}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NanoTime int64  `protobuf:"varint,1,opt,name=nano_time,json=nanoTime,proto3" json:"nano_time,omitempty"`
	Channel  string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	From     string `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	// Types that are assignable to Type:
	//	*Data_Text
	//	*Data_Binary
	Type isData_Type `protobuf_oneof:"type"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{2}
}

func (x *Data) GetNanoTime() int64 {
	if x != nil {
		return x.NanoTime
	}
	return 0
}

func (x *Data) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *Data) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (m *Data) GetType() isData_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Data) GetText() string {
	if x, ok := x.GetType().(*Data_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Data) GetBinary() []byte {
	if x, ok := x.GetType().(*Data_Binary); ok {
		return x.Binary
	}
	return nil
}

type isData_Type interface {
	isData_Type()
}

type Data_Text struct {
	Text string `protobuf:"bytes,4,opt,name=text,proto3,oneof"`
}

type Data_Binary struct {
	Binary []byte `protobuf:"bytes,5,opt,name=binary,proto3,oneof"`
}

func (*Data_Text) isData_Type() {}

func (*Data_Binary) isData_Type() {}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Type:
	//	*Request_Subscribe
	//	*Request_Unsubscribe
	//	*Request_Publish
	Type isRequest_Type `protobuf_oneof:"type"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{3}
}

func (x *Request) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *Request) GetType() isRequest_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Request) GetSubscribe() *Subscribe {
	if x, ok := x.GetType().(*Request_Subscribe); ok {
		return x.Subscribe
	}
	return nil
}

func (x *Request) GetUnsubscribe() *Unsubscribe {
	if x, ok := x.GetType().(*Request_Unsubscribe); ok {
		return x.Unsubscribe
	}
	return nil
}

func (x *Request) GetPublish() *Publish {
	if x, ok := x.GetType().(*Request_Publish); ok {
		return x.Publish
	}
	return nil
}

type isRequest_Type interface {
	isRequest_Type()
}

type Request_Subscribe struct {
	Subscribe *Subscribe `protobuf:"bytes,2,opt,name=subscribe,proto3,oneof"`
}

type Request_Unsubscribe struct {
	Unsubscribe *Unsubscribe `protobuf:"bytes,3,opt,name=unsubscribe,proto3,oneof"`
}

type Request_Publish struct {
	Publish *Publish `protobuf:"bytes,4,opt,name=publish,proto3,oneof"`
}

func (*Request_Subscribe) isRequest_Type() {}

func (*Request_Unsubscribe) isRequest_Type() {}

func (*Request_Publish) isRequest_Type() {}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Type:
	//	*Response_Success
	//	*Response_Error
	Type isResponse_Type `protobuf_oneof:"type"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{4}
}

func (x *Response) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *Response) GetType() isResponse_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Response) GetSuccess() *Success {
	if x, ok := x.GetType().(*Response_Success); ok {
		return x.Success
	}
	return nil
}

func (x *Response) GetError() *Error {
	if x, ok := x.GetType().(*Response_Error); ok {
		return x.Error
	}
	return nil
}

type isResponse_Type interface {
	isResponse_Type()
}

type Response_Success struct {
	Success *Success `protobuf:"bytes,2,opt,name=success,proto3,oneof"`
}

type Response_Error struct {
	Error *Error `protobuf:"bytes,3,opt,name=error,proto3,oneof"`
}

func (*Response_Success) isResponse_Type() {}

func (*Response_Error) isResponse_Type() {}

type Success struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Success) Reset() {
	*x = Success{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Success) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Success) ProtoMessage() {}

func (x *Success) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Success.ProtoReflect.Descriptor instead.
func (*Success) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{5}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Reason string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{6}
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type Subscribe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
}

func (x *Subscribe) Reset() {
	*x = Subscribe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscribe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscribe) ProtoMessage() {}

func (x *Subscribe) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscribe.ProtoReflect.Descriptor instead.
func (*Subscribe) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{7}
}

func (x *Subscribe) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

type Unsubscribe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
}

func (x *Unsubscribe) Reset() {
	*x = Unsubscribe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unsubscribe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unsubscribe) ProtoMessage() {}

func (x *Unsubscribe) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unsubscribe.ProtoReflect.Descriptor instead.
func (*Unsubscribe) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{8}
}

func (x *Unsubscribe) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

type Publish struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channel string `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	// Types that are assignable to Type:
	//	*Publish_Text
	//	*Publish_Binary
	Type isPublish_Type `protobuf_oneof:"type"`
}

func (x *Publish) Reset() {
	*x = Publish{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pubsub_pubsub_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Publish) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Publish) ProtoMessage() {}

func (x *Publish) ProtoReflect() protoreflect.Message {
	mi := &file_api_pubsub_pubsub_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Publish.ProtoReflect.Descriptor instead.
func (*Publish) Descriptor() ([]byte, []int) {
	return file_api_pubsub_pubsub_proto_rawDescGZIP(), []int{9}
}

func (x *Publish) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (m *Publish) GetType() isPublish_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Publish) GetText() string {
	if x, ok := x.GetType().(*Publish_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Publish) GetBinary() []byte {
	if x, ok := x.GetType().(*Publish_Binary); ok {
		return x.Binary
	}
	return nil
}

type isPublish_Type interface {
	isPublish_Type()
}

type Publish_Text struct {
	Text string `protobuf:"bytes,2,opt,name=text,proto3,oneof"`
}

type Publish_Binary struct {
	Binary []byte `protobuf:"bytes,3,opt,name=binary,proto3,oneof"`
}

func (*Publish_Text) isPublish_Type() {}

func (*Publish_Binary) isPublish_Type() {}

var File_api_pubsub_pubsub_proto protoreflect.FileDescriptor

var file_api_pubsub_pubsub_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2f, 0x70, 0x75, 0x62,
	0x73, 0x75, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x69, 0x6f, 0x2e, 0x63, 0x72,
	0x6f, 0x73, 0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61,
	0x70, 0x69, 0x22, 0x84, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3c,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70,
	0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x33, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x6f, 0x2e,
	0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x90, 0x01, 0x0a, 0x07, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x12, 0x3c, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73,
	0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x89, 0x01, 0x0a,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x61, 0x6e, 0x6f, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6e, 0x61, 0x6e, 0x6f, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d,
	0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79,
	0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xed, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f,
	0x73, 0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x48, 0x00, 0x52, 0x09, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x48, 0x0a, 0x0b, 0x75, 0x6e, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62,
	0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x12, 0x3c, 0x0a, 0x07, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x61,
	0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x48, 0x00, 0x52, 0x07, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x98, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x74, 0x61, 0x6c, 0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x48, 0x00, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x36, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x61, 0x6c,
	0x6b, 0x2e, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x22, 0x09, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x33,
	0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x22, 0x25, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0x27, 0x0a, 0x0b, 0x55, 0x6e,
	0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x22, 0x5b, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18,
	0x0a, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00,
	0x52, 0x06, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x42, 0x15, 0x50, 0x01, 0x5a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62,
	0x3b, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_pubsub_pubsub_proto_rawDescOnce sync.Once
	file_api_pubsub_pubsub_proto_rawDescData = file_api_pubsub_pubsub_proto_rawDesc
)

func file_api_pubsub_pubsub_proto_rawDescGZIP() []byte {
	file_api_pubsub_pubsub_proto_rawDescOnce.Do(func() {
		file_api_pubsub_pubsub_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_pubsub_pubsub_proto_rawDescData)
	})
	return file_api_pubsub_pubsub_proto_rawDescData
}

var file_api_pubsub_pubsub_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_pubsub_pubsub_proto_goTypes = []interface{}{
	(*Message)(nil),     // 0: io.crosstalk.pubsub.api.Message
	(*Control)(nil),     // 1: io.crosstalk.pubsub.api.Control
	(*Data)(nil),        // 2: io.crosstalk.pubsub.api.Data
	(*Request)(nil),     // 3: io.crosstalk.pubsub.api.Request
	(*Response)(nil),    // 4: io.crosstalk.pubsub.api.Response
	(*Success)(nil),     // 5: io.crosstalk.pubsub.api.Success
	(*Error)(nil),       // 6: io.crosstalk.pubsub.api.Error
	(*Subscribe)(nil),   // 7: io.crosstalk.pubsub.api.Subscribe
	(*Unsubscribe)(nil), // 8: io.crosstalk.pubsub.api.Unsubscribe
	(*Publish)(nil),     // 9: io.crosstalk.pubsub.api.Publish
}
var file_api_pubsub_pubsub_proto_depIdxs = []int32{
	1, // 0: io.crosstalk.pubsub.api.Message.control:type_name -> io.crosstalk.pubsub.api.Control
	2, // 1: io.crosstalk.pubsub.api.Message.data:type_name -> io.crosstalk.pubsub.api.Data
	3, // 2: io.crosstalk.pubsub.api.Control.request:type_name -> io.crosstalk.pubsub.api.Request
	4, // 3: io.crosstalk.pubsub.api.Control.response:type_name -> io.crosstalk.pubsub.api.Response
	7, // 4: io.crosstalk.pubsub.api.Request.subscribe:type_name -> io.crosstalk.pubsub.api.Subscribe
	8, // 5: io.crosstalk.pubsub.api.Request.unsubscribe:type_name -> io.crosstalk.pubsub.api.Unsubscribe
	9, // 6: io.crosstalk.pubsub.api.Request.publish:type_name -> io.crosstalk.pubsub.api.Publish
	5, // 7: io.crosstalk.pubsub.api.Response.success:type_name -> io.crosstalk.pubsub.api.Success
	6, // 8: io.crosstalk.pubsub.api.Response.error:type_name -> io.crosstalk.pubsub.api.Error
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_api_pubsub_pubsub_proto_init() }
func file_api_pubsub_pubsub_proto_init() {
	if File_api_pubsub_pubsub_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_pubsub_pubsub_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Control); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Success); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscribe); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Unsubscribe); i {
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
		file_api_pubsub_pubsub_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Publish); i {
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
	file_api_pubsub_pubsub_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Message_Control)(nil),
		(*Message_Data)(nil),
	}
	file_api_pubsub_pubsub_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Control_Request)(nil),
		(*Control_Response)(nil),
	}
	file_api_pubsub_pubsub_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Data_Text)(nil),
		(*Data_Binary)(nil),
	}
	file_api_pubsub_pubsub_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Request_Subscribe)(nil),
		(*Request_Unsubscribe)(nil),
		(*Request_Publish)(nil),
	}
	file_api_pubsub_pubsub_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Response_Success)(nil),
		(*Response_Error)(nil),
	}
	file_api_pubsub_pubsub_proto_msgTypes[9].OneofWrappers = []interface{}{
		(*Publish_Text)(nil),
		(*Publish_Binary)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_pubsub_pubsub_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_pubsub_pubsub_proto_goTypes,
		DependencyIndexes: file_api_pubsub_pubsub_proto_depIdxs,
		MessageInfos:      file_api_pubsub_pubsub_proto_msgTypes,
	}.Build()
	File_api_pubsub_pubsub_proto = out.File
	file_api_pubsub_pubsub_proto_rawDesc = nil
	file_api_pubsub_pubsub_proto_goTypes = nil
	file_api_pubsub_pubsub_proto_depIdxs = nil
}