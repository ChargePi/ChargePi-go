## ⚡ EVCC

All EVCCs must implement the `EVCC` interface. Every time a new EVCC is added, the Init method is called. All the
communication must be established in that method.

Add a constant for each implementation. The current state should be persisted in the struct and the GetType should
return the exact type.

```golang
package evcc

import (
	"context"

	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

const (
	PhoenixEMCPPPETH = "EM-CP-PP-ETH"
	Relay            = "Relay"
	Western          = "Western"
)

type EVCC interface {
	Init(ctx context.Context) error
	EnableCharging() error
	DisableCharging()
	SetMaxChargingCurrent(value float64) error
	GetMaxChargingCurrent() float64
	Lock()
	Unlock()
	GetState() CarState
	GetError() string
	Cleanup() error
	GetType() string
	GetStatusChangeChannel() <-chan StateNotification
	Reset()
}

```