package evcc

const (
	StateA1 = CarState("A1")
	StateA2 = CarState("A2")
	StateB1 = CarState("B1")
	StateB2 = CarState("B2")
	StateC1 = CarState("C1")
	StateC2 = CarState("C2")
	StateD1 = CarState("D1")
	StateD2 = CarState("D2")
	StateE  = CarState("E")
	StateF  = CarState("F")
)

type (
	CarState string

	StateNotification struct {
		State CarState
		Error string
	}
)

func NewStateNotification(state CarState, error string) StateNotification {
	return StateNotification{State: state, Error: error}
}

func IsStateValid(state CarState) bool {
	switch state {
	case StateA1,
		StateA2,
		StateB1,
		StateB2,
		StateC1,
		StateC2,
		StateD1,
		StateD2,
		StateE,
		StateF:
		return true
	default:
		return false
	}
}
