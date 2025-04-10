// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: internal/pkg/tunnel/tunnel.proto

package tunnel

import (
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

type Step int32

const (
	Step_REQUEST        Step = 0
	Step_RESPONSE       Step = 1
	Step_SERVER_ELAPSED Step = 2
)

// Enum value maps for Step.
var (
	Step_name = map[int32]string{
		0: "REQUEST",
		1: "RESPONSE",
		2: "SERVER_ELAPSED",
	}
	Step_value = map[string]int32{
		"REQUEST":        0,
		"RESPONSE":       1,
		"SERVER_ELAPSED": 2,
	}
)

func (x Step) Enum() *Step {
	p := new(Step)
	*p = x
	return p
}

func (x Step) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Step) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_pkg_tunnel_tunnel_proto_enumTypes[0].Descriptor()
}

func (Step) Type() protoreflect.EnumType {
	return &file_internal_pkg_tunnel_tunnel_proto_enumTypes[0]
}

func (x Step) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Step.Descriptor instead.
func (Step) EnumDescriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{0}
}

type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record  *Record  `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
	Session *Session `protobuf:"bytes,2,opt,name=session,proto3" json:"session,omitempty"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{0}
}

func (x *ServerMessage) GetRecord() *Record {
	if x != nil {
		return x.Record
	}
	return nil
}

func (x *ServerMessage) GetSession() *Session {
	if x != nil {
		return x.Session
	}
	return nil
}

type ClientMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Record    *Record    `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
	Handshake *Handshake `protobuf:"bytes,2,opt,name=handshake,proto3" json:"handshake,omitempty"`
}

func (x *ClientMessage) Reset() {
	*x = ClientMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessage) ProtoMessage() {}

func (x *ClientMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessage.ProtoReflect.Descriptor instead.
func (*ClientMessage) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{1}
}

func (x *ClientMessage) GetRecord() *Record {
	if x != nil {
		return x.Record
	}
	return nil
}

func (x *ClientMessage) GetHandshake() *Handshake {
	if x != nil {
		return x.Handshake
	}
	return nil
}

type Handshake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	ClientKey    []byte `protobuf:"bytes,2,opt,name=clientKey,proto3" json:"clientKey,omitempty"`
	AllowReaders bool   `protobuf:"varint,3,opt,name=allowReaders,proto3" json:"allowReaders,omitempty"`
	Bearer       []byte `protobuf:"bytes,4,opt,name=bearer,proto3" json:"bearer,omitempty"`
	PublicKey    []byte `protobuf:"bytes,5,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
}

func (x *Handshake) Reset() {
	*x = Handshake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Handshake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handshake) ProtoMessage() {}

func (x *Handshake) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Handshake.ProtoReflect.Descriptor instead.
func (*Handshake) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{2}
}

func (x *Handshake) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Handshake) GetClientKey() []byte {
	if x != nil {
		return x.ClientKey
	}
	return nil
}

func (x *Handshake) GetAllowReaders() bool {
	if x != nil {
		return x.AllowReaders
	}
	return false
}

func (x *Handshake) GetBearer() []byte {
	if x != nil {
		return x.Bearer
	}
	return nil
}

func (x *Handshake) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

type Session struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId        string `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	Hostname        string `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	HttpPort        string `protobuf:"bytes,3,opt,name=httpPort,proto3" json:"httpPort,omitempty"`
	HttpsPort       string `protobuf:"bytes,4,opt,name=httpsPort,proto3" json:"httpsPort,omitempty"`
	Bearer          []byte `protobuf:"bytes,5,opt,name=bearer,proto3" json:"bearer,omitempty"`
	MaxRequestSize  int32  `protobuf:"varint,6,opt,name=maxRequestSize,proto3" json:"maxRequestSize,omitempty"`
	MaxResponseSize int32  `protobuf:"varint,7,opt,name=maxResponseSize,proto3" json:"maxResponseSize,omitempty"`
	ResponseTimeout int64  `protobuf:"varint,8,opt,name=responseTimeout,proto3" json:"responseTimeout,omitempty"`
	ExpireAt        int64  `protobuf:"varint,9,opt,name=expireAt,proto3" json:"expireAt,omitempty"`
}

func (x *Session) Reset() {
	*x = Session{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{3}
}

func (x *Session) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Session) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Session) GetHttpPort() string {
	if x != nil {
		return x.HttpPort
	}
	return ""
}

func (x *Session) GetHttpsPort() string {
	if x != nil {
		return x.HttpsPort
	}
	return ""
}

func (x *Session) GetBearer() []byte {
	if x != nil {
		return x.Bearer
	}
	return nil
}

func (x *Session) GetMaxRequestSize() int32 {
	if x != nil {
		return x.MaxRequestSize
	}
	return 0
}

func (x *Session) GetMaxResponseSize() int32 {
	if x != nil {
		return x.MaxResponseSize
	}
	return 0
}

func (x *Session) GetResponseTimeout() int64 {
	if x != nil {
		return x.ResponseTimeout
	}
	return 0
}

func (x *Session) GetExpireAt() int64 {
	if x != nil {
		return x.ExpireAt
	}
	return 0
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proto      string            `protobuf:"bytes,1,opt,name=proto,proto3" json:"proto,omitempty"`
	Method     string            `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Url        string            `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Path       string            `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Header     map[string]string `protobuf:"bytes,5,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body       []byte            `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	RemoteAddr string            `protobuf:"bytes,7,opt,name=remoteAddr,proto3" json:"remoteAddr,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[4]
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
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{4}
}

func (x *Request) GetProto() string {
	if x != nil {
		return x.Proto
	}
	return ""
}

func (x *Request) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Request) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Request) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Request) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Request) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Request) GetRemoteAddr() string {
	if x != nil {
		return x.RemoteAddr
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Proto         string            `protobuf:"bytes,1,opt,name=proto,proto3" json:"proto,omitempty"`
	Status        string            `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Code          int32             `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
	Header        map[string]string `protobuf:"bytes,4,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body          []byte            `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Elapsed       int64             `protobuf:"varint,6,opt,name=elapsed,proto3" json:"elapsed,omitempty"`
	ServerElapsed int64             `protobuf:"varint,7,opt,name=serverElapsed,proto3" json:"serverElapsed,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[5]
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
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{5}
}

func (x *Response) GetProto() string {
	if x != nil {
		return x.Proto
	}
	return ""
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Response) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Response) GetElapsed() int64 {
	if x != nil {
		return x.Elapsed
	}
	return 0
}

func (x *Response) GetServerElapsed() int64 {
	if x != nil {
		return x.ServerElapsed
	}
	return 0
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Request  *Request  `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	Response *Response `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`
	Step     Step      `protobuf:"varint,4,opt,name=step,proto3,enum=tunnel.Step" json:"step,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{6}
}

func (x *Record) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Record) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *Record) GetResponse() *Response {
	if x != nil {
		return x.Response
	}
	return nil
}

func (x *Record) GetStep() Step {
	if x != nil {
		return x.Step
	}
	return Step_REQUEST
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pkg_tunnel_tunnel_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP(), []int{7}
}

var File_internal_pkg_tunnel_tunnel_proto protoreflect.FileDescriptor

var file_internal_pkg_tunnel_tunnel_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74,
	0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x22, 0x62, 0x0a, 0x0d, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x75,
	0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x06, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x29, 0x0a, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x68,
	0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x26, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52,
	0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x2f, 0x0a, 0x09, 0x68, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x09, 0x68,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x22, 0x9f, 0x01, 0x0a, 0x09, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79,
	0x12, 0x22, 0x0a, 0x0c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0xab, 0x02, 0x0a, 0x07, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x68, 0x74, 0x74, 0x70, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x68, 0x74, 0x74, 0x70, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x68, 0x74,
	0x74, 0x70, 0x73, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x68,
	0x74, 0x74, 0x70, 0x73, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x65, 0x61, 0x72,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72,
	0x12, 0x26, 0x0a, 0x0e, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x6d, 0x61, 0x78, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0f, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x41, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x41, 0x74, 0x22, 0x81, 0x02, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x33, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x41, 0x64, 0x64,
	0x72, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x91, 0x02, 0x0a,
	0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x75,
	0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x12,
	0x24, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6c,
	0x61, 0x70, 0x73, 0x65, 0x64, 0x1a, 0x39, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x93, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x29, 0x0a, 0x07, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x74,
	0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x73, 0x74, 0x65, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x74, 0x65, 0x70,
	0x52, 0x04, 0x73, 0x74, 0x65, 0x70, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x2a,
	0x35, 0x0a, 0x04, 0x53, 0x74, 0x65, 0x70, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x51, 0x55, 0x45,
	0x53, 0x54, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45,
	0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x45, 0x4c, 0x41,
	0x50, 0x53, 0x45, 0x44, 0x10, 0x02, 0x32, 0x72, 0x0a, 0x06, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c,
	0x12, 0x3c, 0x0a, 0x06, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x15, 0x2e, 0x74, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x15, 0x2e, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x2a,
	0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x0d, 0x2e, 0x74, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x74, 0x75, 0x6e, 0x6e,
	0x65, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_pkg_tunnel_tunnel_proto_rawDescOnce sync.Once
	file_internal_pkg_tunnel_tunnel_proto_rawDescData = file_internal_pkg_tunnel_tunnel_proto_rawDesc
)

func file_internal_pkg_tunnel_tunnel_proto_rawDescGZIP() []byte {
	file_internal_pkg_tunnel_tunnel_proto_rawDescOnce.Do(func() {
		file_internal_pkg_tunnel_tunnel_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_pkg_tunnel_tunnel_proto_rawDescData)
	})
	return file_internal_pkg_tunnel_tunnel_proto_rawDescData
}

var file_internal_pkg_tunnel_tunnel_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_pkg_tunnel_tunnel_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_internal_pkg_tunnel_tunnel_proto_goTypes = []interface{}{
	(Step)(0),             // 0: tunnel.Step
	(*ServerMessage)(nil), // 1: tunnel.ServerMessage
	(*ClientMessage)(nil), // 2: tunnel.ClientMessage
	(*Handshake)(nil),     // 3: tunnel.Handshake
	(*Session)(nil),       // 4: tunnel.Session
	(*Request)(nil),       // 5: tunnel.Request
	(*Response)(nil),      // 6: tunnel.Response
	(*Record)(nil),        // 7: tunnel.Record
	(*Empty)(nil),         // 8: tunnel.Empty
	nil,                   // 9: tunnel.Request.HeaderEntry
	nil,                   // 10: tunnel.Response.HeaderEntry
}
var file_internal_pkg_tunnel_tunnel_proto_depIdxs = []int32{
	7,  // 0: tunnel.ServerMessage.record:type_name -> tunnel.Record
	4,  // 1: tunnel.ServerMessage.session:type_name -> tunnel.Session
	7,  // 2: tunnel.ClientMessage.record:type_name -> tunnel.Record
	3,  // 3: tunnel.ClientMessage.handshake:type_name -> tunnel.Handshake
	9,  // 4: tunnel.Request.header:type_name -> tunnel.Request.HeaderEntry
	10, // 5: tunnel.Response.header:type_name -> tunnel.Response.HeaderEntry
	5,  // 6: tunnel.Record.request:type_name -> tunnel.Request
	6,  // 7: tunnel.Record.response:type_name -> tunnel.Response
	0,  // 8: tunnel.Record.step:type_name -> tunnel.Step
	2,  // 9: tunnel.Tunnel.Tunnel:input_type -> tunnel.ClientMessage
	8,  // 10: tunnel.Tunnel.TestConn:input_type -> tunnel.Empty
	1,  // 11: tunnel.Tunnel.Tunnel:output_type -> tunnel.ServerMessage
	8,  // 12: tunnel.Tunnel.TestConn:output_type -> tunnel.Empty
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_internal_pkg_tunnel_tunnel_proto_init() }
func file_internal_pkg_tunnel_tunnel_proto_init() {
	if File_internal_pkg_tunnel_tunnel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessage); i {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Handshake); i {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Session); i {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_internal_pkg_tunnel_tunnel_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_internal_pkg_tunnel_tunnel_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pkg_tunnel_tunnel_proto_goTypes,
		DependencyIndexes: file_internal_pkg_tunnel_tunnel_proto_depIdxs,
		EnumInfos:         file_internal_pkg_tunnel_tunnel_proto_enumTypes,
		MessageInfos:      file_internal_pkg_tunnel_tunnel_proto_msgTypes,
	}.Build()
	File_internal_pkg_tunnel_tunnel_proto = out.File
	file_internal_pkg_tunnel_tunnel_proto_rawDesc = nil
	file_internal_pkg_tunnel_tunnel_proto_goTypes = nil
	file_internal_pkg_tunnel_tunnel_proto_depIdxs = nil
}
