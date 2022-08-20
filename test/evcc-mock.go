package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/evcc"
)

type EvccMock struct {
	mock.Mock
}

func (e *EvccMock) Init(ctx context.Context) error {
	return e.Called().Error(0)
}

func (e *EvccMock) EnableCharging() error {
	return e.Called().Error(0)
}

func (e *EvccMock) DisableCharging() {
	e.Called()
}

func (e *EvccMock) SetMaxChargingCurrent(value float64) error {
	return e.Called().Error(0)
}

func (e *EvccMock) GetMaxChargingCurrent() float64 {
	//TODO implement me
	panic("implement me")
}

func (e *EvccMock) Lock() {
	e.Called()
}

func (e *EvccMock) Unlock() {
	e.Called()
}

func (e *EvccMock) GetState() evcc.CarState {
	return evcc.CarState(e.Called().String(0))
}

func (e *EvccMock) Cleanup() error {
	return e.Called().Error(0)
}

func (e *EvccMock) GetType() string {
	return e.Called().String(0)
}

func (e *EvccMock) GetError() string {
	return e.Called().String(0)
}

func (e *EvccMock) GetStatusChangeChannel() <-chan evcc.StateNotification {
	return e.Called().Get(0).(chan evcc.StateNotification)
}
