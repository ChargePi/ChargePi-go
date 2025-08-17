package grpc

import (
	"context"

	"github.com/samber/lo"

	"github.com/ChargePi/ocpp-manager/ocpp_v16"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"google.golang.org/protobuf/types/known/emptypb"

	configurationv1 "github.com/ChargePi/ChargePi-go/gen/proto/configuration/v1"
	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
)

type ConfigurationHandler struct {
	configurationv1.UnimplementedConfigurationServiceServer
	configurationv1.UnimplementedDisplayServiceServer
	configurationv1.UnimplementedIndicatorServiceServer
	configurationv1.UnimplementedTagReaderServiceServer
	settingsManager manager.Manager
	point           chargePoint.ChargePoint
}

func NewConfigurationHandler(settingsManager manager.Manager) *ConfigurationHandler {
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

	displaySettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Display = &configurationv1.Display{
		Type:     displaySettings.Hardware.Display.Driver,
		Enabled:  displaySettings.Hardware.Display.IsEnabled,
		Language: &displaySettings.Hardware.Display.Language,
		// I2C:      i2cSettings,
	}

	return response, nil
}

func (s *ConfigurationHandler) SetReaderSettings(ctx context.Context, request *configurationv1.SetReaderSettingsRequest) (*configurationv1.SetReaderSettingsResponse, error) {
	response := &configurationv1.SetReaderSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ConfigurationHandler) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*configurationv1.GetReaderSettingsResponse, error) {
	response := &configurationv1.GetReaderSettingsResponse{}

	readerSettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Reader = &configurationv1.Reader{
		Type:    readerSettings.Hardware.TagReader.ReaderModel,
		Enabled: readerSettings.Hardware.TagReader.IsEnabled,
		// DeviceAddress: readerSettings.Device,
	}

	return response, nil
}

func (s *ConfigurationHandler) SetIndicatorSettings(ctx context.Context, request *configurationv1.SetIndicatorSettingsRequest) (*configurationv1.SetIndicatorSettingsResponse, error) {
	response := &configurationv1.SetIndicatorSettingsResponse{
		Status: "Failed",
	}

	return response, nil
}

func (s *ConfigurationHandler) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*configurationv1.GetIndicatorSettingsResponse, error) {
	response := &configurationv1.GetIndicatorSettingsResponse{}

	indicatorSettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Indicator = &configurationv1.Indicator{
		Type:             indicatorSettings.Hardware.Indicator.Type,
		Enabled:          indicatorSettings.Hardware.Indicator.Enabled,
		IndicateCardRead: &indicatorSettings.Hardware.Indicator.IndicateCardRead,
		// Invert:           indicatorSettings.Invert,
	}

	return response, nil
}

func (s *ConfigurationHandler) GetVariables(ctx context.Context, e *emptypb.Empty) (*configurationv1.GetVariablesResponse, error) {
	response := &configurationv1.GetVariablesResponse{}

	configuration, err := s.settingsManager.GetConfiguration()
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

		err := s.settingsManager.UpdateKey(ocpp_v16.Key(variable.GetKey()), lo.ToPtr(variable.GetValue()))
		if err == nil {
			status = "Success"
		}

		response.Statuses = append(response.Statuses, status)
	}

	return response, nil
}

func (s *ConfigurationHandler) GetVariable(ctx context.Context, request *configurationv1.GetVariableRequest) (*configurationv1.GetVariableResponse, error) {
	value, err := s.settingsManager.GetConfigurationValue(ocpp_v16.Key(request.GetKey()))
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

func toDisplay(d *configurationv1.Display) display.Settings {
	return display.Settings{
		IsEnabled: false,
		Driver:    d.GetType(),
		Language:  d.GetLanguage(),
		// I2C:       nil,
	}
}
