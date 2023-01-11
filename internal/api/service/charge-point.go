package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

type ChargePointService struct {
	grpc.UnimplementedChargePointServer
	point       chargePoint.ChargePoint
	ocppManager ocppConfigManager.Manager
}

func NewChargePointService(point chargePoint.ChargePoint, ocppManager ocppConfigManager.Manager) *ChargePointService {
	return &ChargePointService{
		point:       point,
		ocppManager: ocppManager,
	}
}

func (s *ChargePointService) SetDisplaySettings(ctx context.Context, request *grpc.SetDisplaySettingsRequest) (*grpc.SetDisplaySettingsResponse, error) {
	response := &grpc.SetDisplaySettingsResponse{}

	newDisplay, err := display.NewDisplay(toDisplay(request.GetDisplay()))
	if err != nil {
		return nil, err
	}

	err = s.point.SetDisplay(newDisplay)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ChargePointService) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*grpc.GetDisplaySettingsResponse, error) {
	response := &grpc.GetDisplaySettingsResponse{}

	return response, nil
}

func (s *ChargePointService) SetReaderSettings(ctx context.Context, request *grpc.SetReaderSettingsRequest) (*grpc.SetReaderSettingsResponse, error) {
	return nil, nil
}

func (s *ChargePointService) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetReaderSettingsResponse, error) {

	return nil, nil
}

func (s *ChargePointService) SetIndicatorSettings(ctx context.Context, request *grpc.SetIndicatorSettingsRequest) (*grpc.SetIndicatorSettingsResponse, error) {

	return nil, nil
}

func (s *ChargePointService) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetIndicatorSettingsResponse, error) {
	return nil, nil
}

func (s *ChargePointService) Restart(ctx context.Context, request *grpc.RestartRequest) (*empty.Empty, error) {
	err := s.point.Reset(request.Type)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ChargePointService) ChangeConnectionDetails(ctx context.Context, request *grpc.ChangeConnectionDetailsRequest) (*grpc.ChangeConnectionDetailsResponse, error) {
	return nil, nil
}

func (s *ChargePointService) ChangeChargePointDetails(ctx context.Context, request *grpc.ChangeChargePointDetailsRequest) (*grpc.ChangeChargePointDetailsResponse, error) {
	return nil, nil
}

func (s *ChargePointService) GetOCPPVariables(ctx context.Context, e *empty.Empty) (*grpc.GetVariablesResponse, error) {
	response := &grpc.GetVariablesResponse{}

	configuration, err := s.ocppManager.GetConfiguration()
	if err != nil {
		return nil, err
	}

	for _, cfg := range configuration {
		response.Variables = append(response.Variables, toConfiguration(cfg))
	}

	return response, nil
}

func (s *ChargePointService) GetVersion(ctx context.Context, e *empty.Empty) (*grpc.GetVersionResponse, error) {
	return &grpc.GetVersionResponse{
		Version: s.point.GetVersion(),
	}, nil
}

func (s *ChargePointService) GetStatus(ctx context.Context, e *empty.Empty) (*grpc.GetStatusResponse, error) {
	return &grpc.GetStatusResponse{
		Connected: s.point.IsConnected(),
		Status:    s.point.GetStatus(),
	}, nil
}

func (s *ChargePointService) SetOCPPVariables(ctx context.Context, request *grpc.SetVariablesRequest) (*grpc.SetVariablesResponse, error) {
	response := &grpc.SetVariablesResponse{}

	for _, variable := range request.GetVariables() {
		status := "Failed"

		err := s.ocppManager.UpdateKey(variable.Key, variable.Value)
		if err == nil {
			status = "Sucess"
		}

		response.Statuses = append(response.Statuses, status)
	}

	return response, nil
}

func (s *ChargePointService) GetOCPPVariable(ctx context.Context, request *grpc.GetVariableRequest) (*grpc.OcppVariable, error) {
	value, err := s.ocppManager.GetConfigurationValue(request.GetKey())
	if err != nil {
		return nil, err
	}

	return toConfiguration(core.ConfigurationKey{
		Key:      request.Key,
		Readonly: false,
		Value:    value,
	}), nil
}

func (s *ChargePointService) mustEmbedUnimplementedChargePointServer() {
}

func toConfiguration(key core.ConfigurationKey) *grpc.OcppVariable {
	return &grpc.OcppVariable{
		Key:      key.Key,
		Readonly: key.Readonly,
		Value:    key.Value,
	}
}

func toDisplay(display *grpc.Display) settings.Display {
	return settings.Display{
		IsEnabled: false,
		Driver:    display.Type,
		Language:  *display.Language,
		I2C:       nil,
	}
}
