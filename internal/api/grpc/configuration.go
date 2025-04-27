package grpc

import (
	"context"

	commonv1 "github.com/ChargePi/ChargePi-go/gen/proto/common/v1"
	configurationv1 "github.com/ChargePi/ChargePi-go/gen/proto/configuration/v1"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/ChargePi/ChargePi-go/pkg/display"
	commonSettings "github.com/ChargePi/ChargePi-go/pkg/models/settings"
	settings2 "github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConfigurationHandler struct {
	configurationv1.UnimplementedConfigurationServiceServer
	settingsManager cfg.Manager
	point           chargePoint.ChargePoint
}

func NewConfigurationHandler(settingsManager cfg.Manager) *ConfigurationHandler {
	return &ConfigurationHandler{
		settingsManager: settingsManager,
	}
}

func (s *ConfigurationHandler) mustEmbedUnimplementedConfigurationServer() {

}

func (s *ConfigurationHandler) SetDisplaySettings(ctx context.Context, request *configurationv1.SetDisplaySettingsRequest) (*configurationv1.SetDisplaySettingsResponse, error) {
	response := &configurationv1.SetDisplaySettingsResponse{
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

func (s *ConfigurationHandler) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*configurationv1.GetDisplaySettingsResponse, error) {
	response := &configurationv1.GetDisplaySettingsResponse{}

	displaySettings := s.settingsManager.GetChargePointSettings().Hardware.Display

	response.Display = &configurationv1.Display{
		Type:     displaySettings.Driver,
		Enabled:  displaySettings.IsEnabled,
		Language: &displaySettings.Language,
		// I2C:      i2cSettings,
	}

	return response, nil
}

func (s *ChargePointHandler) SetReaderSettings(ctx context.Context, request *configurationv1.SetReaderSettingsRequest) (*configurationv1.SetReaderSettingsResponse, error) {
	response := &configurationv1.SetReaderSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ChargePointHandler) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*configurationv1.GetReaderSettingsResponse, error) {
	response := &configurationv1.GetReaderSettingsResponse{}

	readerSettings := s.settingsManager.GetChargePointSettings().Hardware.TagReader

	response.Reader = &configurationv1.Reader{
		Type:    readerSettings.ReaderModel,
		Enabled: readerSettings.IsEnabled,
		// DeviceAddress: readerSettings.Device,
	}

	return response, nil
}

func (s *ChargePointHandler) SetIndicatorSettings(ctx context.Context, request *configurationv1.SetIndicatorSettingsRequest) (*configurationv1.SetIndicatorSettingsResponse, error) {
	response := &configurationv1.SetIndicatorSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ChargePointHandler) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*configurationv1.GetIndicatorSettingsResponse, error) {
	response := &configurationv1.GetIndicatorSettingsResponse{}

	indicatorSettings := s.settingsManager.GetChargePointSettings().Hardware.Indicator

	response.Indicator = &configurationv1.Indicator{
		Type:             indicatorSettings.Type,
		Enabled:          indicatorSettings.Enabled,
		IndicateCardRead: &indicatorSettings.IndicateCardRead,
		// Invert:           indicatorSettings.Invert,
	}

	return response, nil
}

func (s *ConfigurationHandler) GetVariables(ctx context.Context, e *emptypb.Empty) (*configurationv1.GetVariablesResponse, error) {
	response := &configurationv1.GetVariablesResponse{}

	// todo get the manager depending on the charge point ocpp version
	configuration, err := s.settingsManager.GetOcppV16Manager().GetConfiguration()
	if err != nil {
		return nil, err
	}

	for _, config := range configuration {
		response.Variables = append(response.Variables, toConfiguration(config))
	}

	return response, nil
}

func (s *ConfigurationHandler) SetVariables(ctx context.Context, request *configurationv1.SetVariablesRequest) (*configurationv1.SetVariablesResponse, error) {
	response := &configurationv1.SetVariablesResponse{}

	for _, variable := range request.GetVariables() {
		status := "Failed"

		// todo get the manager depending on the charge point ocpp version
		err := s.settingsManager.GetOcppV16Manager().UpdateKey(ocpp_v16.Key(variable.Key), variable.Value)
		if err == nil {
			status = "Success"
		}

		response.Statuses = append(response.Statuses, status)
	}

	return response, nil
}

func (s *ConfigurationHandler) GetVariable(ctx context.Context, request *configurationv1.GetVariableRequest) (*configurationv1.GetVariableResponse, error) {
	value, err := s.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.Key(request.GetKey()))
	if err != nil {
		return nil, err
	}

	response := &configurationv1.GetVariableResponse{
		Key:      request.GetKey(),
		Value:    value,
		ReadOnly: false,
	}
	return response, nil
}

func (s *ConfigurationHandler) mustEmbedUnimplementedConfigurationServiceServer() {
}

func toConfiguration(key core.ConfigurationKey) *configurationv1.OcppVariable {
	return &configurationv1.OcppVariable{
		Key:      key.Key,
		ReadOnly: key.Readonly,
		Value:    key.Value,
	}
}

func toDisplay(display *configurationv1.Display) settings2.Display {

	return settings2.Display{
		IsEnabled: false,
		Driver:    display.Type,
		Language:  *display.Language,
		// I2C:       nil,
	}
}

func toI2c(i2c commonSettings.I2C) *commonv1.I2C {
	return &commonv1.I2C{
		Address: i2c.Address,
		Bus:     int32(i2c.Bus),
	}
}
