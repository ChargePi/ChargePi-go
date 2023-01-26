package settings

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type (
	EVSE struct {
		EvseId      int                 `json:"evseId,omitempty" yaml:"evseId" mapstructure:"evseId" validate:"required"`
		MaxPower    float32             `json:"maxPower" yaml:"maxPower" mapstructure:"maxPower" validate:"gt=0"`
		EVCC        settings.EVCC       `json:"evcc" yaml:"evcc" mapstructure:"evcc"`
		PowerMeter  settings.PowerMeter `json:"powerMeter" yaml:"powerMeter" mapstructure:"powerMeter"`
		Session     Session             `json:"session" yaml:"session" mapstructure:"session"`
		Status      string              `json:"status,omitempty" yaml:"status" mapstructure:"status"`
		Reservation *int                `json:"reservation,omitempty" yaml:"reservation" mapstructure:"reservation"`
		Connectors  []Connector         `json:"connectors" yaml:"connectors" mapstructure:"connectors"`
	}

	Connector struct {
		ConnectorId int    `json:"connectorId,omitempty" yaml:"connectorId" mapstructure:"connectorId"`
		Type        string `json:"type,omitempty" yaml:"type" mapstructure:"type"`
		Status      string `json:"status,omitempty" yaml:"status" mapstructure:"status"`
	}

	Session struct {
		IsActive      bool               `json:"active" yaml:"active" mapstructure:"active"`
		TransactionId string             `json:"transactionId,omitempty" yaml:"transactionId" mapstructure:"transactionId"`
		TagId         string             `json:"tagId,omitempty" yaml:"tagId" mapstructure:"tagId"`
		Started       string             `json:"started,omitempty" yaml:"started" mapstructure:"started"`
		Consumption   []types.MeterValue `json:"consumption,omitempty" yaml:"consumption" mapstructure:"consumption"`
	}
)
