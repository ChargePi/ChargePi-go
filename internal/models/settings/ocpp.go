package settings

type (
	OCPPInfo struct {
		Vendor                  string `fig:"Vendor" default:"UL FE" json:"vendor,omitempty" yaml:"vendor" mapstructure:"vendor"`
		Model                   string `fig:"Model" default:"ChargePi" json:"model,omitempty" yaml:"model" mapstructure:"model"`
		ChargeBoxSerialNumber   string `fig:"ChargeBoxSerialNumber" default:"" json:"charge_box_serial_number,omitempty" yaml:"charge_box_serial_number" mapstructure:"charge_box_serial_number"`
		ChargePointSerialNumber string `fig:"ChargePointSerialNumber" default:"" json:"charge_point_serial_number,omitempty" yaml:"charge_point_serial_number" mapstructure:"charge_point_serial_number"`
		Iccid                   string `fig:"Iccid" default:"" json:"iccid,omitempty" yaml:"iccid" mapstructure:"iccid"`
		Imsi                    string `fig:"Imsi" default:"" json:"imsi,omitempty" yaml:"imsi" mapstructure:"imsi"`
	}
)
