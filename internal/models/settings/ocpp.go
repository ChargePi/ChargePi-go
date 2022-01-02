package settings

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const (
	OCPP16  = ProtocolVersion("1.6")
	OCPP201 = ProtocolVersion("2.0.1")
)

type (
	ProtocolVersion string

	OCPPInfo struct {
		Vendor                  string `fig:"Vendor" default:"UL FE"`
		Model                   string `fig:"Model" default:"ChargePi"`
		ChargeBoxSerialNumber   string `fig:"ChargeBoxSerialNumber" default:""`
		ChargePointSerialNumber string `fig:"ChargePointSerialNumber" default:""`
		Iccid                   string `fig:"Iccid" default:""`
		Imsi                    string `fig:"Imsi" default:""`
	}

	Connector struct {
		EvseId      int        `fig:"EvseId" validate:"required"`
		ConnectorId int        `fig:"ConnectorId" validate:"required"`
		Type        string     `fig:"Type" validate:"required"`
		Status      string     `fig:"Status" validation:"required"`
		Session     Session    `fig:"Session"`
		Relay       Relay      `fig:"Relay"`
		PowerMeter  PowerMeter `fig:"PowerMeter"`
	}

	Session struct {
		IsActive      bool               `fig:"IsActive"`
		TransactionId string             `fig:"TransactionId" default:""`
		TagId         string             `fig:"TagId" default:""`
		Started       string             `fig:"Started" default:""`
		Consumption   []types.MeterValue `fig:"Consumption"`
	}
)
