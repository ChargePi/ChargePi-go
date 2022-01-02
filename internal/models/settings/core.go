package settings

type (
	Settings struct {
		ChargePoint ChargePoint
	}

	ChargePoint struct {
		Info struct {
			Id              string   `fig:"Id" validate:"required"`
			ProtocolVersion string   `fig:"ProtocolVersion" default:"1.6"`
			ServerUri       string   `fig:"ServerUri" validate:"required"`
			MaxChargingTime int      `fig:"MaxChargingTime" default:"180"`
			OCPPInfo        OCPPInfo `fig:"ocpp"`
		}
		Logging  Logging  `fig:"Logging"`
		TLS      TLS      `fig:"TLS"`
		Hardware Hardware `fig:"Hardware"`
	}

	TLS struct {
		IsEnabled             bool   `fig:"isEnabled"`
		CACertificatePath     string `fig:"CACertificatePath"`
		ClientCertificatePath string `fig:"ClientCertificatePath"`
		ClientKeyPath         string `fig:"ClientKeyPath"`
	}

	Logging struct {
		Type   []string `fig:"type" validate:"required"` // file, remote, console
		Format string   `fig:"format" default:"syslog"`  // gelf, syslog, json, etc
		Host   string   `fig:"host"`
		Port   int      `fig:"port" default:"1514"`
	}
)
