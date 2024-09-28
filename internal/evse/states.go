package evse

type State interface {
	StartCharging(evse *V1) error
	StopCharging(evse *V1) error
	Fault(evse *V1) error
	Suspend(evse *V1) error
	Resume(evse *V1) error
}

type AvailableState struct{}

type ChargingState struct{}

type FaultState struct{}

type SuspendState struct{}

type ResumeState struct{}
