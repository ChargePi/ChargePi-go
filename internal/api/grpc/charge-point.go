package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
	commonSettings "github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
)

type ChargePointService struct {
	grpc.UnimplementedChargePointServer
	point           chargePoint.ChargePoint
	ocppManager     ocppConfigManager.Manager
	settingsManager cfg.Manager
}

func NewChargePointService(point chargePoint.ChargePoint, ocppManager ocppConfigManager.Manager) *ChargePointService {
	return &ChargePointService{
		point:           point,
		ocppManager:     ocppManager,
		settingsManager: cfg.GetManager(),
	}
}

func (s *ChargePointService) SetDisplaySettings(ctx context.Context, request *grpc.SetDisplaySettingsRequest) (*grpc.SetDisplaySettingsResponse, error) {
	response := &grpc.SetDisplaySettingsResponse{
		Status: "Failed",
	}

	displaySettings := toDisplay(request.GetDisplay())

	newDisplay, err := display.NewDisplay(displaySettings)
	if err != nil {
		return response, nil
	}

	err = s.point.SetDisplay(newDisplay)
	if err != nil {
		return response, nil
	}

	// todo set the display settings in the manager

	response.Status = "Success"
	return response, nil
}

func (s *ChargePointService) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*grpc.GetDisplaySettingsResponse, error) {
	response := &grpc.GetDisplaySettingsResponse{}

	displaySettings := s.settingsManager.GetChargePointSettings().Hardware.Display

	response.Display = &grpc.Display{
		Type:     displaySettings.Driver,
		Enabled:  displaySettings.IsEnabled,
		Language: &displaySettings.Language,
		// I2C:      i2cSettings,
	}

	return response, nil
}

func (s *ChargePointService) SetReaderSettings(ctx context.Context, request *grpc.SetReaderSettingsRequest) (*grpc.SetReaderSettingsResponse, error) {
	response := &grpc.SetReaderSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ChargePointService) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetReaderSettingsResponse, error) {
	response := &grpc.GetReaderSettingsResponse{}

	readerSettings := s.settingsManager.GetChargePointSettings().Hardware.TagReader

	response.Reader = &grpc.TagReader{
		Type:    readerSettings.ReaderModel,
		Enabled: readerSettings.IsEnabled,
		// DeviceAddress: readerSettings.Device,
	}

	return response, nil
}

func (s *ChargePointService) SetIndicatorSettings(ctx context.Context, request *grpc.SetIndicatorSettingsRequest) (*grpc.SetIndicatorSettingsResponse, error) {
	response := &grpc.SetIndicatorSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ChargePointService) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetIndicatorSettingsResponse, error) {
	response := &grpc.GetIndicatorSettingsResponse{}

	indicatorSettings := s.settingsManager.GetChargePointSettings().Hardware.Indicator

	response.Indicator = &grpc.Indicator{
		Type:             indicatorSettings.Type,
		Enabled:          indicatorSettings.Enabled,
		IndicateCardRead: &indicatorSettings.IndicateCardRead,
		// Invert:           indicatorSettings.Invert,
	}

	return response, nil
}

func (s *ChargePointService) Restart(ctx context.Context, request *grpc.RestartRequest) (*empty.Empty, error) {
	err := s.point.Reset(request.Type)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ChargePointService) ChangeConnectionDetails(ctx context.Context, request *grpc.ChangeConnectionDetailsRequest) (*grpc.ChangeConnectionDetailsResponse, error) {
	response := &grpc.ChangeConnectionDetailsResponse{}

	return response, nil
}

func (s *ChargePointService) ChangeChargePointDetails(ctx context.Context, request *grpc.ChangeChargePointDetailsRequest) (*grpc.ChangeChargePointDetailsResponse, error) {
	response := &grpc.ChangeChargePointDetailsResponse{}

	return response, nil
}

func (s *ChargePointService) GetOCPPVariables(ctx context.Context, e *empty.Empty) (*grpc.GetVariablesResponse, error) {
	response := &grpc.GetVariablesResponse{}

	configuration, err := s.ocppManager.GetConfiguration()
	if err != nil {
		return nil, err
	}

	for _, config := range configuration {
		response.Variables = append(response.Variables, toConfiguration(config))
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
			status = "Success"
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
		// I2C:       nil,
	}
}

func toI2c(i2c commonSettings.I2C) *grpc.I2C {
	return &grpc.I2C{
		Address: i2c.Address,
		Bus:     int32(i2c.Bus),
	}
}
