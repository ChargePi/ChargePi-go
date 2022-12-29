package test

import (
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
)

type IndicatorMock struct {
	mock.Mock
}

func (i *IndicatorMock) DisplayColor(index int, colorHex indicator.Color) error {
	args := i.Called(index, colorHex)
	return args.Error(0)
}

func (i *IndicatorMock) Blink(index int, times int, colorHex indicator.Color) error {
	args := i.Called(index, times, colorHex)
	return args.Error(0)
}

func (i *IndicatorMock) Cleanup() {
	i.Called()
}

func (i *IndicatorMock) GetType() string {
	return i.Called().String(0)
}
