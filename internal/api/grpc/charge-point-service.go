package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
)

type Service struct {
	UnimplementedChargePointServer
	point      chargePoint.ChargePoint
	tagManager auth.TagManager
}

func NewChargePointService(point chargePoint.ChargePoint, tagManager auth.TagManager) *Service {
	return &Service{
		point:      point,
		tagManager: tagManager,
	}
}

func (s *Service) GetEVSEs(ctx context.Context, empty *empty.Empty) (*GetEvsesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddEVSE(ctx context.Context, request *SetEVCCRequest) (*AddEvseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetEVSE(ctx context.Context, request *GetEvseRequest) (*GetEvseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetEVCC(ctx context.Context, request *SetEVCCRequest) (*SetEvccResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetPowerMeter(ctx context.Context, request *SetPowerMeterRequest) (*SetPowerMeterDetails, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUsageForEVSE(request *GetUsageForEVSERequest, server ChargePoint_GetUsageForEVSEServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetDisplaySettings(ctx context.Context, request *SetDisplaySettingsRequest) (*SetDisplaySettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*GetDisplaySettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetReaderSettings(ctx context.Context, request *SetReaderSettingsRequest) (*SetReaderSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*GetReaderSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetIndicatorSettings(ctx context.Context, request *SetIndicatorSettingsRequest) (*SetIndicatorSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*GetIndicatorSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetAuthorizedCards(ctx context.Context, empty *empty.Empty) (*GetAuthorizedCardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) AddAuthorizedCards(ctx context.Context, request *AddAuthorizedCardsRequest) (*AddAuthorizedCardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Restart(ctx context.Context, request *RestartRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ChangeConnectionDetails(ctx context.Context, request *ChangeConnectionDetailsRequest) (*ChangeConnectionDetailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ChangeChargePointDetails(ctx context.Context, request *ChangeConnectionDetailsRequest) (*ChangeConnectionDetailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedChargePointServer() {
}
