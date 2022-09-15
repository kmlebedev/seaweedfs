// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package master_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SeaweedClient is the client API for Seaweed service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeaweedClient interface {
	SendHeartbeat(ctx context.Context, opts ...grpc.CallOption) (Seaweed_SendHeartbeatClient, error)
	KeepConnected(ctx context.Context, opts ...grpc.CallOption) (Seaweed_KeepConnectedClient, error)
	LookupVolume(ctx context.Context, in *LookupVolumeRequest, opts ...grpc.CallOption) (*LookupVolumeResponse, error)
	Assign(ctx context.Context, in *AssignRequest, opts ...grpc.CallOption) (*AssignResponse, error)
	Statistics(ctx context.Context, in *StatisticsRequest, opts ...grpc.CallOption) (*StatisticsResponse, error)
	CollectionList(ctx context.Context, in *CollectionListRequest, opts ...grpc.CallOption) (*CollectionListResponse, error)
	CollectionDelete(ctx context.Context, in *CollectionDeleteRequest, opts ...grpc.CallOption) (*CollectionDeleteResponse, error)
	VolumeList(ctx context.Context, in *VolumeListRequest, opts ...grpc.CallOption) (*VolumeListResponse, error)
	LookupEcVolume(ctx context.Context, in *LookupEcVolumeRequest, opts ...grpc.CallOption) (*LookupEcVolumeResponse, error)
	VacuumVolume(ctx context.Context, in *VacuumVolumeRequest, opts ...grpc.CallOption) (*VacuumVolumeResponse, error)
	GetMasterConfiguration(ctx context.Context, in *GetMasterConfigurationRequest, opts ...grpc.CallOption) (*GetMasterConfigurationResponse, error)
	ListClusterNodes(ctx context.Context, in *ListClusterNodesRequest, opts ...grpc.CallOption) (*ListClusterNodesResponse, error)
	LeaseAdminToken(ctx context.Context, in *LeaseAdminTokenRequest, opts ...grpc.CallOption) (*LeaseAdminTokenResponse, error)
	ReleaseAdminToken(ctx context.Context, in *ReleaseAdminTokenRequest, opts ...grpc.CallOption) (*ReleaseAdminTokenResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	RaftListClusterServers(ctx context.Context, in *RaftListClusterServersRequest, opts ...grpc.CallOption) (*RaftListClusterServersResponse, error)
	RaftAddServer(ctx context.Context, in *RaftAddServerRequest, opts ...grpc.CallOption) (*RaftAddServerResponse, error)
	RaftRemoveServer(ctx context.Context, in *RaftRemoveServerRequest, opts ...grpc.CallOption) (*RaftRemoveServerResponse, error)
	VolumeMarkReadonly(ctx context.Context, in *VolumeMarkReadonlyRequest, opts ...grpc.CallOption) (*VolumeMarkReadonlyResponse, error)
	VolumeMarkWritable(ctx context.Context, in *VolumeMarkWritableRequest, opts ...grpc.CallOption) (*VolumeMarkWritableResponse, error)
}

type seaweedClient struct {
	cc grpc.ClientConnInterface
}

func NewSeaweedClient(cc grpc.ClientConnInterface) SeaweedClient {
	return &seaweedClient{cc}
}

func (c *seaweedClient) SendHeartbeat(ctx context.Context, opts ...grpc.CallOption) (Seaweed_SendHeartbeatClient, error) {
	stream, err := c.cc.NewStream(ctx, &Seaweed_ServiceDesc.Streams[0], "/master_pb.Seaweed/SendHeartbeat", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedSendHeartbeatClient{stream}
	return x, nil
}

type Seaweed_SendHeartbeatClient interface {
	Send(*Heartbeat) error
	Recv() (*HeartbeatResponse, error)
	grpc.ClientStream
}

type seaweedSendHeartbeatClient struct {
	grpc.ClientStream
}

func (x *seaweedSendHeartbeatClient) Send(m *Heartbeat) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedSendHeartbeatClient) Recv() (*HeartbeatResponse, error) {
	m := new(HeartbeatResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedClient) KeepConnected(ctx context.Context, opts ...grpc.CallOption) (Seaweed_KeepConnectedClient, error) {
	stream, err := c.cc.NewStream(ctx, &Seaweed_ServiceDesc.Streams[1], "/master_pb.Seaweed/KeepConnected", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedKeepConnectedClient{stream}
	return x, nil
}

type Seaweed_KeepConnectedClient interface {
	Send(*KeepConnectedRequest) error
	Recv() (*KeepConnectedResponse, error)
	grpc.ClientStream
}

type seaweedKeepConnectedClient struct {
	grpc.ClientStream
}

func (x *seaweedKeepConnectedClient) Send(m *KeepConnectedRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedKeepConnectedClient) Recv() (*KeepConnectedResponse, error) {
	m := new(KeepConnectedResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedClient) LookupVolume(ctx context.Context, in *LookupVolumeRequest, opts ...grpc.CallOption) (*LookupVolumeResponse, error) {
	out := new(LookupVolumeResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/LookupVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) Assign(ctx context.Context, in *AssignRequest, opts ...grpc.CallOption) (*AssignResponse, error) {
	out := new(AssignResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/Assign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) Statistics(ctx context.Context, in *StatisticsRequest, opts ...grpc.CallOption) (*StatisticsResponse, error) {
	out := new(StatisticsResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/Statistics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) CollectionList(ctx context.Context, in *CollectionListRequest, opts ...grpc.CallOption) (*CollectionListResponse, error) {
	out := new(CollectionListResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/CollectionList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) CollectionDelete(ctx context.Context, in *CollectionDeleteRequest, opts ...grpc.CallOption) (*CollectionDeleteResponse, error) {
	out := new(CollectionDeleteResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/CollectionDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) VolumeList(ctx context.Context, in *VolumeListRequest, opts ...grpc.CallOption) (*VolumeListResponse, error) {
	out := new(VolumeListResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/VolumeList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) LookupEcVolume(ctx context.Context, in *LookupEcVolumeRequest, opts ...grpc.CallOption) (*LookupEcVolumeResponse, error) {
	out := new(LookupEcVolumeResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/LookupEcVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) VacuumVolume(ctx context.Context, in *VacuumVolumeRequest, opts ...grpc.CallOption) (*VacuumVolumeResponse, error) {
	out := new(VacuumVolumeResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/VacuumVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) GetMasterConfiguration(ctx context.Context, in *GetMasterConfigurationRequest, opts ...grpc.CallOption) (*GetMasterConfigurationResponse, error) {
	out := new(GetMasterConfigurationResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/GetMasterConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) ListClusterNodes(ctx context.Context, in *ListClusterNodesRequest, opts ...grpc.CallOption) (*ListClusterNodesResponse, error) {
	out := new(ListClusterNodesResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/ListClusterNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) LeaseAdminToken(ctx context.Context, in *LeaseAdminTokenRequest, opts ...grpc.CallOption) (*LeaseAdminTokenResponse, error) {
	out := new(LeaseAdminTokenResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/LeaseAdminToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) ReleaseAdminToken(ctx context.Context, in *ReleaseAdminTokenRequest, opts ...grpc.CallOption) (*ReleaseAdminTokenResponse, error) {
	out := new(ReleaseAdminTokenResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/ReleaseAdminToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) RaftListClusterServers(ctx context.Context, in *RaftListClusterServersRequest, opts ...grpc.CallOption) (*RaftListClusterServersResponse, error) {
	out := new(RaftListClusterServersResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/RaftListClusterServers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) RaftAddServer(ctx context.Context, in *RaftAddServerRequest, opts ...grpc.CallOption) (*RaftAddServerResponse, error) {
	out := new(RaftAddServerResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/RaftAddServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) RaftRemoveServer(ctx context.Context, in *RaftRemoveServerRequest, opts ...grpc.CallOption) (*RaftRemoveServerResponse, error) {
	out := new(RaftRemoveServerResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/RaftRemoveServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) VolumeMarkReadonly(ctx context.Context, in *VolumeMarkReadonlyRequest, opts ...grpc.CallOption) (*VolumeMarkReadonlyResponse, error) {
	out := new(VolumeMarkReadonlyResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/VolumeMarkReadonly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedClient) VolumeMarkWritable(ctx context.Context, in *VolumeMarkWritableRequest, opts ...grpc.CallOption) (*VolumeMarkWritableResponse, error) {
	out := new(VolumeMarkWritableResponse)
	err := c.cc.Invoke(ctx, "/master_pb.Seaweed/VolumeMarkWritable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SeaweedServer is the server API for Seaweed service.
// All implementations must embed UnimplementedSeaweedServer
// for forward compatibility
type SeaweedServer interface {
	SendHeartbeat(Seaweed_SendHeartbeatServer) error
	KeepConnected(Seaweed_KeepConnectedServer) error
	LookupVolume(context.Context, *LookupVolumeRequest) (*LookupVolumeResponse, error)
	Assign(context.Context, *AssignRequest) (*AssignResponse, error)
	Statistics(context.Context, *StatisticsRequest) (*StatisticsResponse, error)
	CollectionList(context.Context, *CollectionListRequest) (*CollectionListResponse, error)
	CollectionDelete(context.Context, *CollectionDeleteRequest) (*CollectionDeleteResponse, error)
	VolumeList(context.Context, *VolumeListRequest) (*VolumeListResponse, error)
	LookupEcVolume(context.Context, *LookupEcVolumeRequest) (*LookupEcVolumeResponse, error)
	VacuumVolume(context.Context, *VacuumVolumeRequest) (*VacuumVolumeResponse, error)
	GetMasterConfiguration(context.Context, *GetMasterConfigurationRequest) (*GetMasterConfigurationResponse, error)
	ListClusterNodes(context.Context, *ListClusterNodesRequest) (*ListClusterNodesResponse, error)
	LeaseAdminToken(context.Context, *LeaseAdminTokenRequest) (*LeaseAdminTokenResponse, error)
	ReleaseAdminToken(context.Context, *ReleaseAdminTokenRequest) (*ReleaseAdminTokenResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	RaftListClusterServers(context.Context, *RaftListClusterServersRequest) (*RaftListClusterServersResponse, error)
	RaftAddServer(context.Context, *RaftAddServerRequest) (*RaftAddServerResponse, error)
	RaftRemoveServer(context.Context, *RaftRemoveServerRequest) (*RaftRemoveServerResponse, error)
	VolumeMarkReadonly(context.Context, *VolumeMarkReadonlyRequest) (*VolumeMarkReadonlyResponse, error)
	VolumeMarkWritable(context.Context, *VolumeMarkWritableRequest) (*VolumeMarkWritableResponse, error)
	mustEmbedUnimplementedSeaweedServer()
}

// UnimplementedSeaweedServer must be embedded to have forward compatible implementations.
type UnimplementedSeaweedServer struct {
}

func (UnimplementedSeaweedServer) SendHeartbeat(Seaweed_SendHeartbeatServer) error {
	return status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
func (UnimplementedSeaweedServer) KeepConnected(Seaweed_KeepConnectedServer) error {
	return status.Errorf(codes.Unimplemented, "method KeepConnected not implemented")
}
func (UnimplementedSeaweedServer) LookupVolume(context.Context, *LookupVolumeRequest) (*LookupVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LookupVolume not implemented")
}
func (UnimplementedSeaweedServer) Assign(context.Context, *AssignRequest) (*AssignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Assign not implemented")
}
func (UnimplementedSeaweedServer) Statistics(context.Context, *StatisticsRequest) (*StatisticsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Statistics not implemented")
}
func (UnimplementedSeaweedServer) CollectionList(context.Context, *CollectionListRequest) (*CollectionListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionList not implemented")
}
func (UnimplementedSeaweedServer) CollectionDelete(context.Context, *CollectionDeleteRequest) (*CollectionDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectionDelete not implemented")
}
func (UnimplementedSeaweedServer) VolumeList(context.Context, *VolumeListRequest) (*VolumeListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeList not implemented")
}
func (UnimplementedSeaweedServer) LookupEcVolume(context.Context, *LookupEcVolumeRequest) (*LookupEcVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LookupEcVolume not implemented")
}
func (UnimplementedSeaweedServer) VacuumVolume(context.Context, *VacuumVolumeRequest) (*VacuumVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VacuumVolume not implemented")
}
func (UnimplementedSeaweedServer) GetMasterConfiguration(context.Context, *GetMasterConfigurationRequest) (*GetMasterConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMasterConfiguration not implemented")
}
func (UnimplementedSeaweedServer) ListClusterNodes(context.Context, *ListClusterNodesRequest) (*ListClusterNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClusterNodes not implemented")
}
func (UnimplementedSeaweedServer) LeaseAdminToken(context.Context, *LeaseAdminTokenRequest) (*LeaseAdminTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaseAdminToken not implemented")
}
func (UnimplementedSeaweedServer) ReleaseAdminToken(context.Context, *ReleaseAdminTokenRequest) (*ReleaseAdminTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseAdminToken not implemented")
}
func (UnimplementedSeaweedServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSeaweedServer) RaftListClusterServers(context.Context, *RaftListClusterServersRequest) (*RaftListClusterServersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RaftListClusterServers not implemented")
}
func (UnimplementedSeaweedServer) RaftAddServer(context.Context, *RaftAddServerRequest) (*RaftAddServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RaftAddServer not implemented")
}
func (UnimplementedSeaweedServer) RaftRemoveServer(context.Context, *RaftRemoveServerRequest) (*RaftRemoveServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RaftRemoveServer not implemented")
}
func (UnimplementedSeaweedServer) VolumeMarkReadonly(context.Context, *VolumeMarkReadonlyRequest) (*VolumeMarkReadonlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeMarkReadonly not implemented")
}
func (UnimplementedSeaweedServer) VolumeMarkWritable(context.Context, *VolumeMarkWritableRequest) (*VolumeMarkWritableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolumeMarkWritable not implemented")
}
func (UnimplementedSeaweedServer) mustEmbedUnimplementedSeaweedServer() {}

// UnsafeSeaweedServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeaweedServer will
// result in compilation errors.
type UnsafeSeaweedServer interface {
	mustEmbedUnimplementedSeaweedServer()
}

func RegisterSeaweedServer(s grpc.ServiceRegistrar, srv SeaweedServer) {
	s.RegisterService(&Seaweed_ServiceDesc, srv)
}

func _Seaweed_SendHeartbeat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedServer).SendHeartbeat(&seaweedSendHeartbeatServer{stream})
}

type Seaweed_SendHeartbeatServer interface {
	Send(*HeartbeatResponse) error
	Recv() (*Heartbeat, error)
	grpc.ServerStream
}

type seaweedSendHeartbeatServer struct {
	grpc.ServerStream
}

func (x *seaweedSendHeartbeatServer) Send(m *HeartbeatResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedSendHeartbeatServer) Recv() (*Heartbeat, error) {
	m := new(Heartbeat)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Seaweed_KeepConnected_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedServer).KeepConnected(&seaweedKeepConnectedServer{stream})
}

type Seaweed_KeepConnectedServer interface {
	Send(*KeepConnectedResponse) error
	Recv() (*KeepConnectedRequest, error)
	grpc.ServerStream
}

type seaweedKeepConnectedServer struct {
	grpc.ServerStream
}

func (x *seaweedKeepConnectedServer) Send(m *KeepConnectedResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedKeepConnectedServer) Recv() (*KeepConnectedRequest, error) {
	m := new(KeepConnectedRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Seaweed_LookupVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).LookupVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/LookupVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).LookupVolume(ctx, req.(*LookupVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_Assign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).Assign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/Assign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).Assign(ctx, req.(*AssignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_Statistics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatisticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).Statistics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/Statistics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).Statistics(ctx, req.(*StatisticsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_CollectionList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).CollectionList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/CollectionList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).CollectionList(ctx, req.(*CollectionListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_CollectionDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectionDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).CollectionDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/CollectionDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).CollectionDelete(ctx, req.(*CollectionDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_VolumeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).VolumeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/VolumeList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).VolumeList(ctx, req.(*VolumeListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_LookupEcVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupEcVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).LookupEcVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/LookupEcVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).LookupEcVolume(ctx, req.(*LookupEcVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_VacuumVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VacuumVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).VacuumVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/VacuumVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).VacuumVolume(ctx, req.(*VacuumVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_GetMasterConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMasterConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).GetMasterConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/GetMasterConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).GetMasterConfiguration(ctx, req.(*GetMasterConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_ListClusterNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClusterNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).ListClusterNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/ListClusterNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).ListClusterNodes(ctx, req.(*ListClusterNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_LeaseAdminToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaseAdminTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).LeaseAdminToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/LeaseAdminToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).LeaseAdminToken(ctx, req.(*LeaseAdminTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_ReleaseAdminToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseAdminTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).ReleaseAdminToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/ReleaseAdminToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).ReleaseAdminToken(ctx, req.(*ReleaseAdminTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_RaftListClusterServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RaftListClusterServersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).RaftListClusterServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/RaftListClusterServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).RaftListClusterServers(ctx, req.(*RaftListClusterServersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_RaftAddServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RaftAddServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).RaftAddServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/RaftAddServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).RaftAddServer(ctx, req.(*RaftAddServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_RaftRemoveServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RaftRemoveServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).RaftRemoveServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/RaftRemoveServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).RaftRemoveServer(ctx, req.(*RaftRemoveServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_VolumeMarkReadonly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeMarkReadonlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).VolumeMarkReadonly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/VolumeMarkReadonly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).VolumeMarkReadonly(ctx, req.(*VolumeMarkReadonlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Seaweed_VolumeMarkWritable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VolumeMarkWritableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedServer).VolumeMarkWritable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/master_pb.Seaweed/VolumeMarkWritable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedServer).VolumeMarkWritable(ctx, req.(*VolumeMarkWritableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Seaweed_ServiceDesc is the grpc.ServiceDesc for Seaweed service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Seaweed_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "master_pb.Seaweed",
	HandlerType: (*SeaweedServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LookupVolume",
			Handler:    _Seaweed_LookupVolume_Handler,
		},
		{
			MethodName: "Assign",
			Handler:    _Seaweed_Assign_Handler,
		},
		{
			MethodName: "Statistics",
			Handler:    _Seaweed_Statistics_Handler,
		},
		{
			MethodName: "CollectionList",
			Handler:    _Seaweed_CollectionList_Handler,
		},
		{
			MethodName: "CollectionDelete",
			Handler:    _Seaweed_CollectionDelete_Handler,
		},
		{
			MethodName: "VolumeList",
			Handler:    _Seaweed_VolumeList_Handler,
		},
		{
			MethodName: "LookupEcVolume",
			Handler:    _Seaweed_LookupEcVolume_Handler,
		},
		{
			MethodName: "VacuumVolume",
			Handler:    _Seaweed_VacuumVolume_Handler,
		},
		{
			MethodName: "GetMasterConfiguration",
			Handler:    _Seaweed_GetMasterConfiguration_Handler,
		},
		{
			MethodName: "ListClusterNodes",
			Handler:    _Seaweed_ListClusterNodes_Handler,
		},
		{
			MethodName: "LeaseAdminToken",
			Handler:    _Seaweed_LeaseAdminToken_Handler,
		},
		{
			MethodName: "ReleaseAdminToken",
			Handler:    _Seaweed_ReleaseAdminToken_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Seaweed_Ping_Handler,
		},
		{
			MethodName: "RaftListClusterServers",
			Handler:    _Seaweed_RaftListClusterServers_Handler,
		},
		{
			MethodName: "RaftAddServer",
			Handler:    _Seaweed_RaftAddServer_Handler,
		},
		{
			MethodName: "RaftRemoveServer",
			Handler:    _Seaweed_RaftRemoveServer_Handler,
		},
		{
			MethodName: "VolumeMarkReadonly",
			Handler:    _Seaweed_VolumeMarkReadonly_Handler,
		},
		{
			MethodName: "VolumeMarkWritable",
			Handler:    _Seaweed_VolumeMarkWritable_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendHeartbeat",
			Handler:       _Seaweed_SendHeartbeat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "KeepConnected",
			Handler:       _Seaweed_KeepConnected_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "master.proto",
}
