package ocpp

type DataTransferChargePointInfo struct {
	// AC or DC
	Type string `fig:"type" json:"type" yaml:"type" mapstructure:"type"`
	// in kW
	MaxPower float32 `fig:"maxPower" default:"180" json:"maxPower" yaml:"maxPower" mapstructure:"maxPower"`
}

func NewChargePointInfo(chargePointType string, maxPower float32) DataTransferChargePointInfo {
	return DataTransferChargePointInfo{
		Type:     chargePointType,
		MaxPower: maxPower,
	}
}
