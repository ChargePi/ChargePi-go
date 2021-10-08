package settings

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type Settings struct {
	ChargePoint struct {
		Info struct {
			Vendor               string `fig:"Vendor" default:"UL FE"`
			Model                string `fig:"Model" default:"ChargePi"`
			Id                   string `fig:"Id" validate:"required"`
			ProtocolVersion      string `fig:"ProtocolVersion" default:"1.6"`
			CurrentClientVersion string `fig:"CurrentClientVersion" default:"1.0"`
			TargetClientVersion  string `fig:"TargetClientVersion" default:"1.0"`
			ServerUri            string `fig:"ServerUri" validate:"required"`
			LogServer            string `fig:"LogServer" validate:"required"`
			MaxChargingTime      int    `fig:"MaxChargingTime" default:"180"`
			TLS                  struct {
				IsEnabled             bool   `fig:"isEnabled"`
				CACertificatePath     string `fig:"CACertificatePath"`
				ClientCertificatePath string `fig:"ClientCertificatePath"`
				ClientKeyPath         string `fig:"ClientKeyPath"`
			}
		}
		Hardware struct {
			Lcd struct {
				IsSupported bool   `fig:"IsSupported"`
				Driver      string `fig:"Driver"`
				Language    string `fig:"Language"`
				I2CAddress  string `fig:"I2CAddress"`
				I2CBus      int    `fig:"I2CBus"`
			}
			TagReader struct {
				IsSupported bool   `fig:"IsSupported"`
				ReaderModel string `fig:"ReaderModel"`
				Device      string `fig:"Device"`
				ResetPin    int    `fig:"ResetPin"`
			}
			LedIndicator struct {
				Enabled          bool   `fig:"Enabled"`
				DataPin          int    `fig:"DataPin"`
				IndicateCardRead bool   `fig:"IndicateCardRead"`
				Type             string `fig:"Type"`
				Invert           bool   `fig:"Invert"`
			}
			PowerMeters struct {
				MinPower int `fig:"MinPower" default:"20"`
				Retries  int `fig:"Retries" default:"3"`
			}
		}
	}
}

type Connector struct {
	EvseId      int    `fig:"EvseId" validate:"required"`
	ConnectorId int    `fig:"ConnectorId" validate:"required"`
	Type        string `fig:"Type" validate:"required"`
	Status      string `fig:"Status" validation:"required"`
	Session     struct {
		IsActive      bool   `fig:"IsActive"`
		TransactionId string `fig:"TransactionId" default:""`
		TagId         string `fig:"TagId" default:""`
		Started       string `fig:"Started" default:""`
		Consumption   []types.MeterValue
	} `fig:"Session"`
	Relay struct {
		RelayPin     int  `fig:"RelayPin" validate:"required"`
		InverseLogic bool `fig:"InverseLogic"`
	} `fig:"Relay"`
	PowerMeter struct {
		Enabled              bool    `fig:"Enabled"`
		Type                 string  `fig:"Type"`
		PowerMeterPin        int     `fig:"PowerMeterPin"`
		SpiBus               int     `fig:"SpiBus" default:"0"`
		PowerUnits           string  `fig:"PowerUnits"`
		Consumption          float64 `fig:"Consumption"`
		ShuntOffset          float64 `fig:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
	} `fig:"PowerMeter"`
}

type Session struct {
	IsActive      bool
	TransactionId string
	TagId         string
	Started       string
	Consumption   []types.MeterValue
}

type OCPPConfig struct {
	Version int
	Keys    []core.ConfigurationKey
}
