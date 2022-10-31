package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
)

type ManagerMock struct {
	mock.Mock
}

/*------------------ Manager mock ------------------*/

func (o *ManagerMock) GetEVSEs() []evse.EVSE {
	return o.Called().Get(0).([]evse.EVSE)
}

func (o *ManagerMock) FindEVSE(evseId int) (evse.EVSE, error) {
	args := o.Called(evseId)
	if args.Get(0) != nil {
		return args.Get(0).(evse.EVSE), args.Error(1)
	}

	return nil, args.Error(1)
}

func (o *ManagerMock) FindAvailableEVSE() (evse.EVSE, error) {
	args := o.Called()
	if args.Get(0) != nil {
		return args.Get(0).(evse.EVSE), args.Error(1)
	}

	return nil, args.Error(1)
}

func (o *ManagerMock) FindEVSEWithTagId(tagId string) (evse.EVSE, error) {
	args := o.Called(tagId)
	if args.Get(0) != nil {
		return args.Get(0).(evse.EVSE), args.Error(1)
	}

	return nil, args.Error(1)
}

func (o *ManagerMock) FindEVSEWithTransactionId(transactionId string) (evse.EVSE, error) {
	args := o.Called(transactionId)
	if args.Get(0) != nil {
		return args.Get(0).(evse.EVSE), args.Error(1)
	}

	return nil, args.Error(1)
}

func (o *ManagerMock) FindEVSEWithReservationId(reservationId int) (evse.EVSE, error) {
	args := o.Called(reservationId)
	if args.Get(0) != nil {
		return args.Get(0).(evse.EVSE), args.Error(1)
	}

	return nil, args.Error(1)
}

func (o *ManagerMock) StartCharging(evseId int, tagId, transactionId string) error {
	return o.Called(evseId, tagId, transactionId).Error(0)
}

func (o *ManagerMock) StopCharging(tagId, transactionId string, reason core.Reason) error {
	return o.Called(tagId, transactionId, reason).Error(0)
}

func (o *ManagerMock) StopAllEVSEs(reason core.Reason) error {
	return o.Called(reason).Error(0)
}

func (o *ManagerMock) AddEVSE(c evse.EVSE) error {
	return o.Called(c).Error(0)
}

func (o *ManagerMock) AddEVSEFromSettings(maxChargingTime *int, c *settings.EVSE) error {
	return o.Called(maxChargingTime, c).Error(0)
}

func (o *ManagerMock) AddEVSEsFromSettings(maxChargingTime *int, c []*settings.EVSE) error {
	return o.Called(maxChargingTime, c).Error(0)
}

func (o *ManagerMock) RestoreEVSEStatus(s *settings.EVSE) error {
	return o.Called(s).Error(0)
}

func (o *ManagerMock) SetNotificationChannel(notificationChannel chan chargePoint.StatusNotification) {
	o.Called()
}

func (o *ManagerMock) SetMeterValuesChannel(notificationChannel chan chargePoint.MeterValueNotification) {
	o.Called()
}
