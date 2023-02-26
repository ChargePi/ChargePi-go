## ⚡ Power meters

All power meters must implement the `PowerMeter` interface. They should be able to, in one form or another, provide the
basic functionality of a power meter: read the current, voltage and energy from the Power Meter hardware.

Various communication protocols are used to get the data from the reader. All communication should be initialized in the
`Init` method. The init method gets called before the program calls any of other methods.

When adding a new power meter, declare a constant in the `power-meter` file.

```golang
package powerMeter

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// Supported power meters
const (
	TypeC5460A = "cs5460a"
)

// PowerMeter is an abstraction for measurement hardware.
type PowerMeter interface {
	Init(ctx context.Context) error
	Reset()
	GetEnergy() types.SampledValue
	GetPower() types.SampledValue
	GetCurrent(phase int) types.SampledValue
	GetVoltage(phase int) types.SampledValue
	GetRMSCurrent(phase int) types.SampledValue
	GetRMSVoltage(phase int) types.SampledValue
	GetType() string
}

```