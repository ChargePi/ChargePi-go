package integration_tests

import (
	"context"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display"
	powerMeter "github.com/xBlaz3kx/ChargePi-go/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/data/session"
)

type (
	displayMock struct {
		mock.Mock
	}

	readerMock struct {
		mock.Mock
	}

	centralSystemMock struct {
		mock.Mock
	}

	connectorMock struct {
		mock.Mock
		connector.Connector
	}
)

/*------------------ CentralSystem mock ------------------*/

func (c *centralSystemMock) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (confirmation *core.AuthorizeConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.AuthorizeConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (confirmation *core.BootNotificationConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.BootNotificationConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnDataTransfer(chargePointId string, request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.DataTransferConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnHeartbeat(chargePointId string, request *core.HeartbeatRequest) (confirmation *core.HeartbeatConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.HeartbeatConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnMeterValues(chargePointId string, request *core.MeterValuesRequest) (confirmation *core.MeterValuesConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.MeterValuesConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnStatusNotification(chargePointId string, request *core.StatusNotificationRequest) (confirmation *core.StatusNotificationConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StatusNotificationConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnStartTransaction(chargePointId string, request *core.StartTransactionRequest) (confirmation *core.StartTransactionConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StartTransactionConfirmation), args.Error(1)
}

func (c *centralSystemMock) OnStopTransaction(chargePointId string, request *core.StopTransactionRequest) (confirmation *core.StopTransactionConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StopTransactionConfirmation), args.Error(1)
}

/*------------------ Display mock ------------------*/

func (l *displayMock) DisplayMessage(message display.LCDMessage) {
	l.Called(message)
}

func (l *displayMock) ListenForMessages(ctx context.Context) {
	l.Called()
}

func (l *displayMock) Cleanup() {
	l.Called()
}

func (l *displayMock) Clear() {
	l.Called()
}

func (l *displayMock) GetLcdChannel() chan<- display.LCDMessage {
	return l.Called().Get(0).(chan display.LCDMessage)
}

/*------------------ Reader mock ------------------*/

func (r *readerMock) ListenForTags(ctx context.Context) {
	r.Called()
}

func (r *readerMock) Cleanup() {
	r.Called()
}

func (r *readerMock) Reset() {
	r.Called()
}

func (r *readerMock) GetTagChannel() <-chan string {
	return r.Called().Get(0).(chan string)
}

/*------------------ Connector mock ------------------*/

func (m *connectorMock) StartCharging(transactionId string, tagId string) error {
	args := m.Called(transactionId, tagId)
	return args.Error(0)
}

func (m *connectorMock) ResumeCharging(session session.Session) (error, int) {
	args := m.Called(session)
	return args.Error(0), args.Int(1)
}

func (m *connectorMock) StopCharging(reason core.Reason) error {
	args := m.Called(reason)
	return args.Error(0)
}

func (m *connectorMock) SetNotificationChannel(notificationChannel chan<- rxgo.Item) {
	m.Called(notificationChannel)
}

func (m *connectorMock) ReserveConnector(reservationId int) error {
	args := m.Called(reservationId)
	return args.Error(0)
}

func (m *connectorMock) RemoveReservation() error {
	args := m.Called()
	return args.Error(0)
}

func (m *connectorMock) GetReservationId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) GetTagId() string {
	args := m.Called()
	return args.String(0)
}

func (m *connectorMock) GetTransactionId() string {
	args := m.Called()
	return args.String(0)
}

func (m *connectorMock) GetConnectorId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) GetEvseId() int {
	args := m.Called()
	return args.Int(0)
}

func (m *connectorMock) CalculateSessionAvgEnergyConsumption() float32 {
	args := m.Called()
	return args.Get(0).(float32)
}

func (m *connectorMock) SamplePowerMeter(measurands []types.Measurand) {
	m.Called(measurands)
}

func (m *connectorMock) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	m.Called(status, errCode)
}

func (m *connectorMock) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	args := m.Called()
	return core.ChargePointStatus(args.String(0)), core.ChargePointErrorCode(args.String(1))
}

func (m *connectorMock) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsPreparing() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsCharging() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsReserved() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) IsUnavailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *connectorMock) GetPowerMeter() powerMeter.PowerMeter {
	args := m.Called()
	return args.Get(0).(powerMeter.PowerMeter)
}

func (m *connectorMock) GetMaxChargingTime() int {
	args := m.Called()
	return args.Int(0)
}
