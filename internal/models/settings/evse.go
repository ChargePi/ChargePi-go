package settings

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

type (
	EVSE struct {
		EvseId     int         `fig:"evseId" validate:"required" json:"evseId,omitempty" yaml:"evseId" mapstructure:"evseId"`
		MaxPower   float32     `fig:"maxPower" json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		EVCC       EVCC        `fig:"evcc" json:"evcc" yaml:"evcc" mapstructure:"evcc"`
		PowerMeter PowerMeter  `fig:"PowerMeter" json:"PowerMeter" yaml:"PowerMeter" mapstructure:"PowerMeter"`
		Session    Session     `fig:"Session" json:"session" yaml:"session" mapstructure:"session"`
		Status     string      `fig:"Status" validation:"required" json:"status,omitempty" yaml:"status" mapstructure:"status"`
		Connectors []Connector `fig:"connectors" json:"connectors" yaml:"connectors" mapstructure:"connectors"`
	}

	Connector struct {
		ConnectorId int    `fig:"ConnectorId" json:"ConnectorId,omitempty" yaml:"ConnectorId" mapstructure:"ConnectorId"`
		Type        string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		Status      string `fig:"Status" json:"status,omitempty" yaml:"status" mapstructure:"status"`
	}

	Session struct {
		IsActive      bool               `fig:"IsActive" json:"IsActive,omitempty" yaml:"IsActive" mapstructure:"IsActive"`
		TransactionId string             `fig:"TransactionId" default:"" json:"TransactionId,omitempty" yaml:"TransactionId" mapstructure:"TransactionId"`
		TagId         string             `fig:"TagId" default:"" json:"TagId,omitempty" yaml:"TagId" mapstructure:"TagId"`
		Started       string             `fig:"Started" default:"" json:"started,omitempty" yaml:"started" mapstructure:"started"`
		Consumption   []types.MeterValue `fig:"Consumption" json:"consumption,omitempty" yaml:"consumption" mapstructure:"consumption"`
	}
)
