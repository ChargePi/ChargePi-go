package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	connectionv1 "github.com/ChargePi/ChargePi-go/gen/proto/connection/v1"
)

type ConnectivityHandler struct {
	connectionv1.UnimplementedConnectionServiceServer
}

func NewConnectivityHandler() *ConnectivityHandler {
	return &ConnectivityHandler{}
}

func (c *ConnectivityHandler) GetConnectionDetails(ctx context.Context, empty *emptypb.Empty) (*connectionv1.GetConnectionDetailsResponse, error) {
	return &connectionv1.GetConnectionDetailsResponse{}, nil
}

func (c *ConnectivityHandler) ChangeConnectionDetails(ctx context.Context, request *connectionv1.ChangeConnectionDetailsRequest) (*connectionv1.ChangeConnectionDetailsResponse, error) {
	return &connectionv1.ChangeConnectionDetailsResponse{}, nil
}

func (c *ConnectivityHandler) mustEmbedUnimplementedConnectionServiceServer() {
}
