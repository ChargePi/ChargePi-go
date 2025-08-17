package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	grpc "github.com/ChargePi/ChargePi-go/gen/proto/charge_point/v1"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
)

type ChargePointHandler struct {
	grpc.UnimplementedChargePointServiceServer
	point chargepoint.ChargePoint
}

func NewChargePointService(point chargepoint.ChargePoint) *ChargePointHandler {
	return &ChargePointHandler{
		point: point,
	}
}

func (s *ChargePointHandler) Restart(ctx context.Context, request *grpc.RestartRequest) (*empty.Empty, error) {
	err := s.point.Reset(ctx, request.Type)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ChargePointHandler) ChangeChargePointDetails(ctx context.Context, request *grpc.ChangeChargePointDetailsRequest) (*grpc.ChangeChargePointDetailsResponse, error) {
	response := &grpc.ChangeChargePointDetailsResponse{}

	return response, nil
}

func (s *ChargePointHandler) GetVersion(ctx context.Context, e *empty.Empty) (*grpc.GetVersionResponse, error) {
	return &grpc.GetVersionResponse{
		Version: s.point.GetVersion(),
	}, nil
}

func (s *ChargePointHandler) GetStatus(ctx context.Context, e *empty.Empty) (*grpc.GetStatusResponse, error) {
	return &grpc.GetStatusResponse{
		Connected: s.point.IsConnected(),
		Status:    s.point.GetStatus(),
	}, nil
}

func (s *ChargePointHandler) mustEmbedUnimplementedChargePointServer() {
}
