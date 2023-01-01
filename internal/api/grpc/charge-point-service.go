package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/users"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

type Service struct {
	UnimplementedChargePointServer
	UnimplementedEvseServer
	UnimplementedLogServer
	UnimplementedTagServer
	UnimplementedUserServer
	point       chargePoint.ChargePoint
	tagManager  auth.TagManager
	evseManager evse.Manager
	ocppManager ocppConfigManager.Manager
	userService users.Service
}

func NewGrpcService(
	point chargePoint.ChargePoint,
	tagManager auth.TagManager,
	manager evse.Manager,
	configurationManager ocppConfigManager.Manager,
	userService users.Service,
) *Service {
	return &Service{
		point:       point,
		tagManager:  tagManager,
		evseManager: manager,
		ocppManager: configurationManager,
		userService: userService,
	}
}

func (s *Service) GetEVSEs(ctx context.Context, empty *empty.Empty) (*GetEvsesResponse, error) {
	s.evseManager.GetEVSEs()
	return nil, nil
}

func (s *Service) AddEVSE(ctx context.Context, request *SetEVCCRequest) (*AddEvseResponse, error) {
	return nil, nil
}

func (s *Service) GetEVSE(ctx context.Context, request *GetEvseRequest) (*GetEvseResponse, error) {
	return nil, nil
}

func (s *Service) SetEVCC(ctx context.Context, request *SetEVCCRequest) (*SetEvccResponse, error) {
	return nil, nil
}

func (s *Service) SetPowerMeter(ctx context.Context, request *SetPowerMeterRequest) (*SetPowerMeterDetails, error) {
	return nil, nil
}

func (s *Service) GetUsageForEVSE(request *GetUsageForEVSERequest, server Evse_GetUsageForEVSEServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedEvseServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetLogs(e *empty.Empty, server Log_GetLogsServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedLogServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetDisplaySettings(ctx context.Context, request *SetDisplaySettingsRequest) (*SetDisplaySettingsResponse, error) {
	return nil, nil
}

func (s *Service) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*GetDisplaySettingsResponse, error) {
	return nil, nil
}

func (s *Service) SetReaderSettings(ctx context.Context, request *SetReaderSettingsRequest) (*SetReaderSettingsResponse, error) {
	return nil, nil
}

func (s *Service) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*GetReaderSettingsResponse, error) {
	return nil, nil
}

func (s *Service) SetIndicatorSettings(ctx context.Context, request *SetIndicatorSettingsRequest) (*SetIndicatorSettingsResponse, error) {
	return nil, nil
}

func (s *Service) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*GetIndicatorSettingsResponse, error) {
	return nil, nil
}

func (s *Service) GetAuthorizedCards(ctx context.Context, empty *empty.Empty) (*GetAuthorizedCardsResponse, error) {
	return nil, nil
}

func (s *Service) AddAuthorizedCards(ctx context.Context, request *AddAuthorizedCardsRequest) (*AddAuthorizedCardsResponse, error) {
	return nil, nil
}

func (s *Service) RemoveAuthorizedCard(ctx context.Context, request *AddAuthorizedCardsRequest) (*AddAuthorizedCardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedTagServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Restart(ctx context.Context, request *RestartRequest) (*empty.Empty, error) {
	return nil, nil
}

func (s *Service) ChangeConnectionDetails(ctx context.Context, request *ChangeConnectionDetailsRequest) (*ChangeConnectionDetailsResponse, error) {
	return nil, nil
}

func (s *Service) ChangeChargePointDetails(context.Context, *ChangeChargePointDetailsRequest) (*ChangeChargePointDetailsResponse, error) {
	return nil, nil
}

func (s *Service) mustEmbedUnimplementedChargePointServer() {
}

func (s *Service) AddUser(ctx context.Context, request *AddUserRequest) (*AddUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUsers(ctx context.Context, e *empty.Empty) (*GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetUser(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) RemoveUser(ctx context.Context, request *RemoveUserRequest) (*RemoveUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}
