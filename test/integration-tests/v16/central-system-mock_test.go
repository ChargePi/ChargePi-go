package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
)

type (
	centralSystemV16Mock struct {
		mock.Mock
	}
)

/*------------------ CentralSystem mock ------------------*/

func (c *centralSystemV16Mock) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (confirmation *core.AuthorizeConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.AuthorizeConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (confirmation *core.BootNotificationConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.BootNotificationConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnDataTransfer(chargePointId string, request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.DataTransferConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnHeartbeat(chargePointId string, request *core.HeartbeatRequest) (confirmation *core.HeartbeatConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.HeartbeatConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnMeterValues(chargePointId string, request *core.MeterValuesRequest) (confirmation *core.MeterValuesConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.MeterValuesConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnStatusNotification(chargePointId string, request *core.StatusNotificationRequest) (confirmation *core.StatusNotificationConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StatusNotificationConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnStartTransaction(chargePointId string, request *core.StartTransactionRequest) (confirmation *core.StartTransactionConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StartTransactionConfirmation), args.Error(1)
}

func (c *centralSystemV16Mock) OnStopTransaction(chargePointId string, request *core.StopTransactionRequest) (confirmation *core.StopTransactionConfirmation, err error) {
	args := c.Called(chargePointId, request)
	return args.Get(0).(*core.StopTransactionConfirmation), args.Error(1)
}
