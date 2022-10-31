package test

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
)

type (
	EvseMock struct {
		mock.Mock
	}
)

/*------------------ EVSE mock ------------------*/

func (m *EvseMock) Init(ctx context.Context) error {
	return m.Called().Error(0)
}

func (m *EvseMock) StopCharging(reason core.Reason) error {
	args := m.Called(reason)
	return args.Error(0)
}

func (m *EvseMock) StartCharging(transactionId, tagId string, connectorId *int) error {
	args := m.Called(transactionId, tagId, connectorId)
	return args.Error(0)
}

func (m *EvseMock) ResumeCharging(session session.Session) (*int, error) {
	args := m.Called(session)
	if args.Get(0) != nil {
		return args.Get(0).(*int), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *EvseMock) GetConnectors() []evse.Connector {
	return m.Called().Get(0).([]evse.Connector)
}

func (m *EvseMock) ReserveEvse(reservationId int, tagId string) error {
	return m.Called(reservationId, tagId, tagId).Error(0)
}

func (m *EvseMock) SetMeterValuesChannel(notificationChannel chan<- chargePoint.MeterValueNotification) {
	m.Called(notificationChannel)
}

func (m *EvseMock) SetNotificationChannel(notificationChannel chan<- chargePoint.StatusNotification) {
	m.Called(notificationChannel)
}

func (m *EvseMock) ReserveConnector(reservationId int, tagId string) error {
	args := m.Called(reservationId, tagId)
	return args.Error(0)
}

func (m *EvseMock) RemoveReservation() error {
	args := m.Called()
	return args.Error(0)
}

func (m *EvseMock) GetReservationId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *EvseMock) GetTagId() string {
	args := m.Called()
	return args.String(0)
}

func (m *EvseMock) GetTransactionId() string {
	args := m.Called()
	return args.String(0)
}

func (m *EvseMock) GetConnectorId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *EvseMock) GetEvseId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *EvseMock) CalculateSessionAvgEnergyConsumption() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *EvseMock) SamplePowerMeter(measurands []types.Measurand) {
	m.Called(measurands)
}

func (m *EvseMock) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	m.Called(status, errCode)
}

func (m *EvseMock) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	args := m.Called()
	return core.ChargePointStatus(args.String(0)), core.ChargePointErrorCode(args.String(1))
}

func (m *EvseMock) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *EvseMock) IsPreparing() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *EvseMock) IsCharging() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *EvseMock) IsReserved() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *EvseMock) IsUnavailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *EvseMock) GetPowerMeter() powerMeter.PowerMeter {
	args := m.Called()
	return args.Get(0).(powerMeter.PowerMeter)
}

func (m *EvseMock) GetMaxChargingTime() *int {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*int)
	}

	return nil
}

func (m *EvseMock) SetAvailability(isAvailable bool) {
	m.Called(isAvailable)
}

func (m *EvseMock) Lock() {
	m.Called()
}

func (m *EvseMock) Unlock() {
	m.Called()
}
