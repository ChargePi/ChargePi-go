package settings

type (
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
		Consumption          float64 `fig:"Consumption"`
		ShuntOffset          float64 `fig:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
	}

	PowerMeters struct {
		MinPower int `fig:"MinPower" default:"20"`
		Retries  int `fig:"Retries" default:"3"`
	}
)
