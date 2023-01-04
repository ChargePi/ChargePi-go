package settings

type (
	OCPPDetails struct {
		Vendor                  string `json:"vendor" yaml:"vendor" mapstructure:"vendor" validate:"required"`
		Model                   string `json:"model" yaml:"model" mapstructure:"model" validate:"required"`
		ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" yaml:"chargeBoxSerialNumber,omitempty" mapstructure:"chargeBoxSerialNumber,omitempty"`
		ChargePointSerialNumber string `json:"pointSerialNumber" yaml:"pointSerialNumber" mapstructure:"pointSerialNumber"`
		Iccid                   string `json:"iccid,omitempty" yaml:"iccid,omitempty" mapstructure:"iccid"`
		Imsi                    string `json:"imsi,omitempty" yaml:"imsi,omitempty" mapstructure:"imsi"`
	}
)
