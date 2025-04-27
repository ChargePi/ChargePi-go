package grpc

import (
	"context"

	"github.com/ChargePi/ChargePi-go/gen/proto/connection/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConnectivityHandler struct {
	connectionv1.UnimplementedConnectionServiceServer
}

func NewConnectivityHandler() *ConnectivityHandler {
	return &ConnectivityHandler{}
}

func (c *ConnectivityHandler) GetConnectionDetails(ctx context.Context, empty *emptypb.Empty) (*connectionv1.GetConnectionDetailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConnectivityHandler) ChangeConnectionDetails(ctx context.Context, request *connectionv1.ChangeConnectionDetailsRequest) (*connectionv1.ChangeConnectionDetailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConnectivityHandler) mustEmbedUnimplementedConnectionServiceServer() {
}
