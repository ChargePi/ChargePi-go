package ocpp

type DataTransferChargePointInfo struct {

	// The Charge Point type. Permitted: AC or DC
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// Max charging power rating (in kW)
	MaxPower float32 `json:"maxPower" yaml:"maxPower" mapstructure:"maxPower"`
}

func NewChargePointInfo(chargePointType string, maxPower float32) DataTransferChargePointInfo {
	return DataTransferChargePointInfo{
		Type:     chargePointType,
		MaxPower: maxPower,
	}
}
