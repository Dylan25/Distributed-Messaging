// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chatpb/chat.proto

package chatpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Letter struct {
	User                 string   `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Letter) Reset()         { *m = Letter{} }
func (m *Letter) String() string { return proto.CompactTextString(m) }
func (*Letter) ProtoMessage()    {}
func (*Letter) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{0}
}

func (m *Letter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Letter.Unmarshal(m, b)
}
func (m *Letter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Letter.Marshal(b, m, deterministic)
}
func (m *Letter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Letter.Merge(m, src)
}
func (m *Letter) XXX_Size() int {
	return xxx_messageInfo_Letter.Size(m)
}
func (m *Letter) XXX_DiscardUnknown() {
	xxx_messageInfo_Letter.DiscardUnknown(m)
}

var xxx_messageInfo_Letter proto.InternalMessageInfo

func (m *Letter) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Letter) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type ChatRequest struct {
	Msg                  *Letter  `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChatRequest) Reset()         { *m = ChatRequest{} }
func (m *ChatRequest) String() string { return proto.CompactTextString(m) }
func (*ChatRequest) ProtoMessage()    {}
func (*ChatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{1}
}

func (m *ChatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatRequest.Unmarshal(m, b)
}
func (m *ChatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatRequest.Marshal(b, m, deterministic)
}
func (m *ChatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatRequest.Merge(m, src)
}
func (m *ChatRequest) XXX_Size() int {
	return xxx_messageInfo_ChatRequest.Size(m)
}
func (m *ChatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChatRequest proto.InternalMessageInfo

func (m *ChatRequest) GetMsg() *Letter {
	if m != nil {
		return m.Msg
	}
	return nil
}

type ChatResponse struct {
	Msg                  *Letter  `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChatResponse) Reset()         { *m = ChatResponse{} }
func (m *ChatResponse) String() string { return proto.CompactTextString(m) }
func (*ChatResponse) ProtoMessage()    {}
func (*ChatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{2}
}

func (m *ChatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatResponse.Unmarshal(m, b)
}
func (m *ChatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatResponse.Marshal(b, m, deterministic)
}
func (m *ChatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatResponse.Merge(m, src)
}
func (m *ChatResponse) XXX_Size() int {
	return xxx_messageInfo_ChatResponse.Size(m)
}
func (m *ChatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChatResponse proto.InternalMessageInfo

func (m *ChatResponse) GetMsg() *Letter {
	if m != nil {
		return m.Msg
	}
	return nil
}

type ConnectRequest struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectRequest) Reset()         { *m = ConnectRequest{} }
func (m *ConnectRequest) String() string { return proto.CompactTextString(m) }
func (*ConnectRequest) ProtoMessage()    {}
func (*ConnectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{3}
}

func (m *ConnectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectRequest.Unmarshal(m, b)
}
func (m *ConnectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectRequest.Marshal(b, m, deterministic)
}
func (m *ConnectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectRequest.Merge(m, src)
}
func (m *ConnectRequest) XXX_Size() int {
	return xxx_messageInfo_ConnectRequest.Size(m)
}
func (m *ConnectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectRequest proto.InternalMessageInfo

func (m *ConnectRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type ConnectResponse struct {
	Ack                  string   `protobuf:"bytes,1,opt,name=ack,proto3" json:"ack,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectResponse) Reset()         { *m = ConnectResponse{} }
func (m *ConnectResponse) String() string { return proto.CompactTextString(m) }
func (*ConnectResponse) ProtoMessage()    {}
func (*ConnectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{4}
}

func (m *ConnectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectResponse.Unmarshal(m, b)
}
func (m *ConnectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectResponse.Marshal(b, m, deterministic)
}
func (m *ConnectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectResponse.Merge(m, src)
}
func (m *ConnectResponse) XXX_Size() int {
	return xxx_messageInfo_ConnectResponse.Size(m)
}
func (m *ConnectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectResponse proto.InternalMessageInfo

func (m *ConnectResponse) GetAck() string {
	if m != nil {
		return m.Ack
	}
	return ""
}

type ListenRequest struct {
	User                 string   `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListenRequest) Reset()         { *m = ListenRequest{} }
func (m *ListenRequest) String() string { return proto.CompactTextString(m) }
func (*ListenRequest) ProtoMessage()    {}
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{5}
}

func (m *ListenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListenRequest.Unmarshal(m, b)
}
func (m *ListenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListenRequest.Marshal(b, m, deterministic)
}
func (m *ListenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListenRequest.Merge(m, src)
}
func (m *ListenRequest) XXX_Size() int {
	return xxx_messageInfo_ListenRequest.Size(m)
}
func (m *ListenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListenRequest proto.InternalMessageInfo

func (m *ListenRequest) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

type ListenResponse struct {
	Msg                  *Letter  `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListenResponse) Reset()         { *m = ListenResponse{} }
func (m *ListenResponse) String() string { return proto.CompactTextString(m) }
func (*ListenResponse) ProtoMessage()    {}
func (*ListenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_41cf34dba7f8bdf6, []int{6}
}

func (m *ListenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListenResponse.Unmarshal(m, b)
}
func (m *ListenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListenResponse.Marshal(b, m, deterministic)
}
func (m *ListenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListenResponse.Merge(m, src)
}
func (m *ListenResponse) XXX_Size() int {
	return xxx_messageInfo_ListenResponse.Size(m)
}
func (m *ListenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListenResponse proto.InternalMessageInfo

func (m *ListenResponse) GetMsg() *Letter {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*Letter)(nil), "chat.Letter")
	proto.RegisterType((*ChatRequest)(nil), "chat.ChatRequest")
	proto.RegisterType((*ChatResponse)(nil), "chat.ChatResponse")
	proto.RegisterType((*ConnectRequest)(nil), "chat.ConnectRequest")
	proto.RegisterType((*ConnectResponse)(nil), "chat.ConnectResponse")
	proto.RegisterType((*ListenRequest)(nil), "chat.ListenRequest")
	proto.RegisterType((*ListenResponse)(nil), "chat.ListenResponse")
}

func init() { proto.RegisterFile("chatpb/chat.proto", fileDescriptor_41cf34dba7f8bdf6) }

var fileDescriptor_41cf34dba7f8bdf6 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0xc5, 0x69, 0x14, 0xc1, 0xb5, 0x04, 0x7a, 0x30, 0x54, 0x1d, 0x50, 0xe5, 0x2e, 0x5d, 0x08,
	0x51, 0x3b, 0xb0, 0xd3, 0xb5, 0x53, 0xd8, 0xd8, 0xd2, 0xe8, 0x44, 0x23, 0x44, 0x62, 0xe2, 0x2b,
	0x42, 0xe2, 0xe7, 0x91, 0x73, 0x4e, 0x14, 0x24, 0x24, 0x3a, 0xf9, 0xf9, 0xf9, 0xdd, 0x7b, 0x77,
	0x3e, 0x98, 0x16, 0x87, 0x9c, 0xcd, 0xfe, 0xc1, 0x1d, 0x89, 0x69, 0x6a, 0xae, 0x31, 0x74, 0x58,
	0xa7, 0x10, 0xed, 0x88, 0x99, 0x1a, 0x44, 0x08, 0x8f, 0x96, 0x9a, 0x99, 0x5a, 0xa8, 0xd5, 0x45,
	0xd6, 0x62, 0xc7, 0x31, 0x7d, 0xf1, 0x2c, 0x10, 0xce, 0x61, 0x7d, 0x0f, 0xe3, 0xed, 0x21, 0xe7,
	0x8c, 0x3e, 0x8e, 0x64, 0x19, 0xef, 0x60, 0xf4, 0x6e, 0x5f, 0xdb, 0xaa, 0xf1, 0x7a, 0x92, 0xb4,
	0x01, 0xe2, 0x98, 0xb9, 0x07, 0x9d, 0xc0, 0x44, 0xe4, 0xd6, 0xd4, 0x95, 0xa5, 0x7f, 0xf5, 0x0b,
	0x88, 0xb7, 0x75, 0x55, 0x51, 0xd1, 0x27, 0xc4, 0x10, 0x94, 0xc6, 0xb7, 0x15, 0x94, 0x46, 0x2f,
	0xe1, 0xaa, 0x57, 0x78, 0xd3, 0x6b, 0x18, 0xe5, 0xc5, 0x9b, 0xd7, 0x38, 0xa8, 0x97, 0x70, 0xb9,
	0x2b, 0x2d, 0x53, 0xd5, 0xb9, 0xfc, 0x31, 0x9e, 0x4e, 0x21, 0xee, 0x44, 0xa7, 0x75, 0xb7, 0xfe,
	0x96, 0xe1, 0x9f, 0xa9, 0xf9, 0x2c, 0x0b, 0xc2, 0x0d, 0x84, 0xee, 0x8a, 0x53, 0x51, 0x0e, 0xfe,
	0x65, 0x8e, 0x43, 0x4a, 0xdc, 0xf5, 0xd9, 0x4a, 0xa5, 0x0a, 0x1f, 0x21, 0x92, 0x54, 0xbc, 0xf1,
	0x01, 0xc3, 0x46, 0xe7, 0xb7, 0xbf, 0xc9, 0xae, 0x34, 0x55, 0x4f, 0xe7, 0x2f, 0x91, 0xac, 0x71,
	0x1f, 0xb5, 0x2b, 0xdc, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0xbc, 0x3e, 0xd8, 0x89, 0xd7, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatServiceClient interface {
	//rpc Connect(ConnectRequest) returns (ConnectResponse) {};
	Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error)
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (ChatService_ListenClient, error)
}

type chatServiceClient struct {
	cc *grpc.ClientConn
}

func NewChatServiceClient(cc *grpc.ClientConn) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (ChatService_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChatService_serviceDesc.Streams[0], "/chat.ChatService/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceChatClient{stream}
	return x, nil
}

type ChatService_ChatClient interface {
	Send(*ChatRequest) error
	Recv() (*ChatResponse, error)
	grpc.ClientStream
}

type chatServiceChatClient struct {
	grpc.ClientStream
}

func (x *chatServiceChatClient) Send(m *ChatRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceChatClient) Recv() (*ChatResponse, error) {
	m := new(ChatResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (ChatService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ChatService_serviceDesc.Streams[1], "/chat.ChatService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_ListenClient interface {
	Recv() (*ListenResponse, error)
	grpc.ClientStream
}

type chatServiceListenClient struct {
	grpc.ClientStream
}

func (x *chatServiceListenClient) Recv() (*ListenResponse, error) {
	m := new(ListenResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServiceServer is the server API for ChatService service.
type ChatServiceServer interface {
	//rpc Connect(ConnectRequest) returns (ConnectResponse) {};
	Chat(ChatService_ChatServer) error
	Listen(*ListenRequest, ChatService_ListenServer) error
}

func RegisterChatServiceServer(s *grpc.Server, srv ChatServiceServer) {
	s.RegisterService(&_ChatService_serviceDesc, srv)
}

func _ChatService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Chat(&chatServiceChatServer{stream})
}

type ChatService_ChatServer interface {
	Send(*ChatResponse) error
	Recv() (*ChatRequest, error)
	grpc.ServerStream
}

type chatServiceChatServer struct {
	grpc.ServerStream
}

func (x *chatServiceChatServer) Send(m *ChatResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceChatServer) Recv() (*ChatRequest, error) {
	m := new(ChatRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ChatService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).Listen(m, &chatServiceListenServer{stream})
}

type ChatService_ListenServer interface {
	Send(*ListenResponse) error
	grpc.ServerStream
}

type chatServiceListenServer struct {
	grpc.ServerStream
}

func (x *chatServiceListenServer) Send(m *ListenResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _ChatService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _ChatService_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Listen",
			Handler:       _ChatService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chatpb/chat.proto",
}
