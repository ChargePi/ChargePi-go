package evcc

const (
	// Car is unplugged
	StateA1 = CarState("A1")
	StateA2 = CarState("A2")

	// Car connected
	StateB1 = CarState("B1")
	StateB2 = CarState("B2")

	// Car wants to charge
	StateC1 = CarState("C1")
	// EVCC allowed charging
	StateC2 = CarState("C2")

	// Charging with ventilation
	StateD1 = CarState("D1")
	StateD2 = CarState("D2")

	// Error state
	StateE = CarState("E")

	// Fault state
	StateF = CarState("F")
)

type CarState string

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
