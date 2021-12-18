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

	Settings struct {
		ChargePoint ChargePoint
	}

	ChargePoint struct {
		Info struct {
			Id              string   `fig:"Id" validate:"required"`
			ProtocolVersion string   `fig:"ProtocolVersion" default:"1.6"`
			ServerUri       string   `fig:"ServerUri" validate:"required"`
			LogServer       string   `fig:"LogServer" validate:"required"`
			MaxChargingTime int      `fig:"MaxChargingTime" default:"180"`
			OCPPInfo        OCPPInfo `fig:"ocpp"`
		}
		TLS      TLS      `fig:"TLS"`
		Hardware Hardware `fig:"Hardware"`
	}

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

	TLS struct {
		IsEnabled             bool   `fig:"isEnabled"`
		CACertificatePath     string `fig:"CACertificatePath"`
		ClientCertificatePath string `fig:"ClientCertificatePath"`
		ClientKeyPath         string `fig:"ClientKeyPath"`
	}

	/* ------------- Hardware structs ------------*/

	Hardware struct {
		Lcd          Lcd          `fig:"lcd"`
		TagReader    TagReader    `fig:"tagReader"`
		LedIndicator LedIndicator `fig:"ledIndicator"`
	}

	Relay struct {
		RelayPin     int  `fig:"RelayPin" validate:"required"`
		InverseLogic bool `fig:"InverseLogic"`
	}

	LedIndicator struct {
		Enabled          bool   `fig:"Enabled"`
		DataPin          int    `fig:"DataPin"`
		IndicateCardRead bool   `fig:"IndicateCardRead"`
		Type             string `fig:"Type"`
		Invert           bool   `fig:"Invert"`
	}

	TagReader struct {
		IsEnabled   bool   `fig:"IsEnabled"`
		ReaderModel string `fig:"ReaderModel"`
		Device      string `fig:"Device"`
		ResetPin    int    `fig:"ResetPin"`
	}

	Lcd struct {
		IsEnabled  bool   `fig:"IsEnabled"`
		Driver     string `fig:"Driver"`
		Language   string `fig:"Language"`
		I2CAddress string `fig:"I2CAddress"`
		I2CBus     int    `fig:"I2CBus"`
	}

	PowerMeter struct {
		Enabled              bool    `fig:"Enabled"`
		Type                 string  `fig:"Type"`
		PowerMeterPin        int     `fig:"PowerMeterPin"`
		SpiBus               int     `fig:"SpiBus" default:"0"`
		PowerUnits           string  `fig:"PowerUnits"`
		Consumption          float64 `fig:"Consumption"`
		ShuntOffset          float64 `fig:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
	}

	PowerMeters struct {
		MinPower int `fig:"MinPower" default:"20"`
		Retries  int `fig:"Retries" default:"3"`
	}
)
