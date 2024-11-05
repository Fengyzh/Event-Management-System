// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: EventService.proto

package eventservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EventService_CreateEvent_FullMethodName      = "/eventservice.EventService/CreateEvent"
	EventService_GetAllEvents_FullMethodName     = "/eventservice.EventService/GetAllEvents"
	EventService_GetEvent_FullMethodName         = "/eventservice.EventService/GetEvent"
	EventService_UpdateEvent_FullMethodName      = "/eventservice.EventService/UpdateEvent"
	EventService_DeleteEvent_FullMethodName      = "/eventservice.EventService/DeleteEvent"
	EventService_OrderEventTicket_FullMethodName = "/eventservice.EventService/OrderEventTicket"
)

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	CreateEvent(ctx context.Context, in *EventCreateRequest, opts ...grpc.CallOption) (*EventResponse, error)
	GetAllEvents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*EventList, error)
	GetEvent(ctx context.Context, in *EventId, opts ...grpc.CallOption) (*Event, error)
	UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*EventResponse, error)
	DeleteEvent(ctx context.Context, in *EventId, opts ...grpc.CallOption) (*EventResponse, error)
	OrderEventTicket(ctx context.Context, in *EventOrderRequest, opts ...grpc.CallOption) (*EventResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *EventCreateRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_CreateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetAllEvents(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*EventList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventList)
	err := c.cc.Invoke(ctx, EventService_GetAllEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetEvent(ctx context.Context, in *EventId, opts ...grpc.CallOption) (*Event, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Event)
	err := c.cc.Invoke(ctx, EventService_GetEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_UpdateEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *EventId, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_DeleteEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) OrderEventTicket(ctx context.Context, in *EventOrderRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EventService_OrderEventTicket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility.
type EventServiceServer interface {
	CreateEvent(context.Context, *EventCreateRequest) (*EventResponse, error)
	GetAllEvents(context.Context, *Empty) (*EventList, error)
	GetEvent(context.Context, *EventId) (*Event, error)
	UpdateEvent(context.Context, *Event) (*EventResponse, error)
	DeleteEvent(context.Context, *EventId) (*EventResponse, error)
	OrderEventTicket(context.Context, *EventOrderRequest) (*EventResponse, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventServiceServer struct{}

func (UnimplementedEventServiceServer) CreateEvent(context.Context, *EventCreateRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventServiceServer) GetAllEvents(context.Context, *Empty) (*EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllEvents not implemented")
}
func (UnimplementedEventServiceServer) GetEvent(context.Context, *EventId) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedEventServiceServer) UpdateEvent(context.Context, *Event) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventServiceServer) DeleteEvent(context.Context, *EventId) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventServiceServer) OrderEventTicket(context.Context, *EventOrderRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderEventTicket not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}
func (UnimplementedEventServiceServer) testEmbeddedByValue()                      {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	// If the following call pancis, it indicates UnimplementedEventServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_CreateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*EventCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetAllEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetAllEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_GetAllEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetAllEvents(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_GetEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEvent(ctx, req.(*EventId))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_UpdateEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_DeleteEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*EventId))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_OrderEventTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).OrderEventTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_OrderEventTicket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).OrderEventTicket(ctx, req.(*EventOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "eventservice.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "GetAllEvents",
			Handler:    _EventService_GetAllEvents_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _EventService_GetEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
		{
			MethodName: "OrderEventTicket",
			Handler:    _EventService_OrderEventTicket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
