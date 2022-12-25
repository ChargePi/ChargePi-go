package test

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

type PowerMeterMock struct {
	mock.Mock
}

/*---------------------- Power Meter Mock ----------------------*/

func (p *PowerMeterMock) Init(ctx context.Context) error {
	return p.Called().Error(0)
}

func (p *PowerMeterMock) Reset() {
	p.Called()
}

func (p *PowerMeterMock) GetEnergy() types.SampledValue {
	args := p.Called()
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetPower() types.SampledValue {
	args := p.Called()
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetCurrent(phase int) types.SampledValue {
	args := p.Called(phase)
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetVoltage(phase int) types.SampledValue {
	args := p.Called(phase)
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetRMSCurrent(phase int) types.SampledValue {
	args := p.Called(phase)
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetRMSVoltage(phase int) types.SampledValue {
	args := p.Called(phase)
	return args.Get(0).(types.SampledValue)
}

func (p *PowerMeterMock) GetType() string {
	return ""
}
