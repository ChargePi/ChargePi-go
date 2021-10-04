package power_meter

import (
	"fmt"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
)

const (
	TypeC5460A = "cs5460a"
)

// PowerMeter is an abstraction for measurement hardware.
type PowerMeter interface {
	Reset()
	GetEnergy() float64
	GetPower() float64
	GetCurrent() float64
	GetVoltage() float64
	GetRMSCurrent() float64
	GetRMSVoltage() float64
}

// NewPowerMeter creates a new power meter based on the connector settings.
func NewPowerMeter(connector *settings.Connector) (PowerMeter, error) {
	if connector.PowerMeter.Enabled {
		log.Println("Creating a new power meter:", connector.PowerMeter.Type)
		switch connector.PowerMeter.Type {
		case TypeC5460A:
			return NewCS5460PowerMeter(
				connector.PowerMeter.PowerMeterPin,
				connector.PowerMeter.SpiBus,
				connector.PowerMeter.ShuntOffset,
				connector.PowerMeter.VoltageDividerOffset,
			)
		default:
			return nil, fmt.Errorf("power meter type not supported")
		}
	}
	return nil, fmt.Errorf("power meter not enabled")
}
