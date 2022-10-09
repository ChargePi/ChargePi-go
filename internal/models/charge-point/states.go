package chargePoint

import "github.com/xBlaz3kx/ChargePi-go/pkg/models/evcc"

type (
	StateNotification struct {
		State evcc.CarState
		Error string
	}
)

func NewStateNotification(state evcc.CarState, error string) StateNotification {
	return StateNotification{State: state, Error: error}
}
