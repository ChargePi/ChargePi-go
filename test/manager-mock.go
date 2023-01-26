package test

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
)

type ManagerMock struct {
	mock.Mock
}

func (o *ManagerMock) InitAll(ctx context.Context) error {
	return o.Called().Error(0)
}

func (o *ManagerMock) UpdateEVSE(ctx context.Context, c evse.EVSE) error {
	return o.Called(c).Error(0)
}

func (o *ManagerMock) RemoveEVSE(evseId int) error {
	return o.Called(evseId).Error(0)
}

func (o *ManagerMock) SetMaxChargingTime(maxChargingTime int) error {
	return o.Called(maxChargingTime).Error(0)
}

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

func (o *ManagerMock) AddEVSE(ctx context.Context, c evse.EVSE) error {
	return o.Called(c).Error(0)
}

func (o *ManagerMock) RestoreEVSEs() error {
	return o.Called().Error(0)
}
func (o *ManagerMock) SetNotificationChannel(notificationChannel chan notifications.StatusNotification) {
	o.Called()
}

func (o *ManagerMock) SetMeterValuesChannel(notificationChannel chan notifications.MeterValueNotification) {
	o.Called()
}

func (o *ManagerMock) GetNotificationChannel() chan notifications.StatusNotification {
	args := o.Called()
	return args.Get(0).(chan notifications.StatusNotification)
}
