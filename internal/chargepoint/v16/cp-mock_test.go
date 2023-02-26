package v16

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

type chargePointMock struct {
	mock.Mock
}

func (c *chargePointMock) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *core.BootNotificationRequest)) (*core.BootNotificationConfirmation, error) {
	args := c.Called(chargePointModel, chargePointVendor)

	if args.Get(0) != nil {
		return args.Get(0).(*core.BootNotificationConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) Authorize(idTag string, props ...func(request *core.AuthorizeRequest)) (*core.AuthorizeConfirmation, error) {
	args := c.Called(idTag)

	if args.Get(0) != nil {
		return args.Get(0).(*core.AuthorizeConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) DataTransfer(vendorId string, props ...func(request *core.DataTransferRequest)) (*core.DataTransferConfirmation, error) {
	args := c.Called(vendorId)

	if args.Get(0) != nil {
		return args.Get(0).(*core.DataTransferConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) Heartbeat(props ...func(request *core.HeartbeatRequest)) (*core.HeartbeatConfirmation, error) {
	args := c.Called()

	if args.Get(0) != nil {
		return args.Get(0).(*core.HeartbeatConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) MeterValues(connectorId int, meterValues []types.MeterValue, props ...func(request *core.MeterValuesRequest)) (*core.MeterValuesConfirmation, error) {
	args := c.Called(connectorId, meterValues)

	if args.Get(0) != nil {
		return args.Get(0).(*core.MeterValuesConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *types.DateTime, props ...func(request *core.StartTransactionRequest)) (*core.StartTransactionConfirmation, error) {
	args := c.Called(connectorId, idTag, meterStart)

	if args.Get(0) != nil {
		return args.Get(0).(*core.StartTransactionConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) StopTransaction(meterStop int, timestamp *types.DateTime, transactionId int, props ...func(request *core.StopTransactionRequest)) (*core.StopTransactionConfirmation, error) {
	args := c.Called(transactionId, meterStop)

	if args.Get(0) != nil {
		return args.Get(0).(*core.StopTransactionConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) StatusNotification(connectorId int, errorCode core.ChargePointErrorCode, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) (*core.StatusNotificationConfirmation, error) {
	args := c.Called(connectorId, status)

	if args.Get(0) != nil {
		return args.Get(0).(*core.StatusNotificationConfirmation), args.Error(1)
	}

	return nil, args.Error(1)
}

func (c *chargePointMock) DiagnosticsStatusNotification(status firmware.DiagnosticsStatus, props ...func(request *firmware.DiagnosticsStatusNotificationRequest)) (*firmware.DiagnosticsStatusNotificationConfirmation, error) {
	panic("implement me")
}

func (c *chargePointMock) FirmwareStatusNotification(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationConfirmation, error) {
	panic("implement me")
}

func (c *chargePointMock) SetCoreHandler(listener core.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SetLocalAuthListHandler(listener localauth.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SetFirmwareManagementHandler(listener firmware.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SetReservationHandler(listener reservation.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SetRemoteTriggerHandler(listener remotetrigger.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SetSmartChargingHandler(listener smartcharging.ChargePointHandler) {
	c.Called()
}

func (c *chargePointMock) SendRequest(request ocpp.Request) (ocpp.Response, error) {
	args := c.Called(request)
	return args.Get(0).(ocpp.Response), args.Error(1)
}

func (c *chargePointMock) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	args := c.Called(request)

	go func() {
		time.Sleep(time.Millisecond * 100)
		callback(args.Get(0).(ocpp.Response), args.Error(1))
	}()

	return args.Error(2)
}

func (c *chargePointMock) Start(centralSystemUrl string) error {
	return c.Called(centralSystemUrl).Error(0)
}

func (c *chargePointMock) Stop() {
	c.Called()
}

func (c *chargePointMock) Errors() <-chan error {
	c.Called()
	return nil
}
