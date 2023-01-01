package evse

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

func (evse *Impl) IsAvailable() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusAvailable && evse.availability == core.AvailabilityTypeOperative
}

func (evse *Impl) IsCharging() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusCharging
}

func (evse *Impl) IsPreparing() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusPreparing
}

func (evse *Impl) IsReserved() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusReserved
}

func (evse *Impl) IsUnavailable() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusUnavailable
}

func (evse *Impl) SetAvailability(isAvailable bool) {
	if isAvailable {
		evse.availability = core.AvailabilityTypeOperative
		return
	}

	evse.availability = core.AvailabilityTypeInoperative
}

func (evse *Impl) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})
	logInfo.Debugf("Setting evse status %s with err %s", status, errCode)

	evse.mu.Lock()
	defer evse.mu.Unlock()

	evse.status = status
	evse.errorCode = errCode
	settings.UpdateEVSEStatus(evse.evseId, status)

	if evse.notificationChannel != nil {
		evse.notificationChannel <- notifications.NewStatusNotification(evse.evseId, string(status), string(errCode))
	}
}

func (evse *Impl) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	return evse.status, evse.errorCode
}
