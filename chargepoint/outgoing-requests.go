package chargepoint

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"os"
	"strconv"
	"time"
)

// bootNotification Notify the central system that the charging point is online. Set the setHeartbeat interval and call restoreState.
// If the central system does not accept the charge point, exit the client.
func (handler *ChargePointHandler) bootNotification() {
	request := core.BootNotificationRequest{
		ChargePointModel:  handler.Settings.ChargePoint.Info.Model,
		ChargePointVendor: handler.Settings.ChargePoint.Info.Vendor,
	}
	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)
		if bootConf.Status == core.RegistrationStatusAccepted {
			log.Printf("Notified and accepted from the central system")
			handler.setHeartbeat(bootConf.Interval)
			handler.restoreState()
		} else {
			log.Printf("Denied by the central system.")
			os.Exit(1)
		}
	}
	handler.SendRequest(request, callback)
}

func (handler *ChargePointHandler) setHeartbeat(interval int) {
	heartBeatInterval, _ := settings.GetConfigurationValue("HeartbeatInterval")
	if interval > 0 {
		heartBeatInterval = fmt.Sprintf("%d", interval)
	}
	heartBeatInterval = fmt.Sprintf("%ss", heartBeatInterval)
	_, err := scheduler.GetScheduler().Every(heartBeatInterval).Tag("heartbeat").Do(handler.sendHeartBeat)
	if err != nil {
		fmt.Println(err)
	}
}

// SendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func (handler *ChargePointHandler) SendRequest(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	var (
		retryIntervalValue string
		maxMessageAttempts string
		maxRetries         int
		retryInterval      int
		attemptErr         error
		intervalErr        error
	)

	maxMessageAttempts, attemptErr = settings.GetConfigurationValue("TransactionMessageAttempts")
	retryIntervalValue, intervalErr = settings.GetConfigurationValue("TransactionMessageRetryInterval")

	maxRetries, convError := strconv.Atoi(maxMessageAttempts)
	if attemptErr != nil || convError != nil {
		maxRetries = 5
	}

	retryInterval, convError = strconv.Atoi(retryIntervalValue)
	if intervalErr != nil || convError != nil {
		retryInterval = 30
	}

	err := retry.Do(
		func() error {
			return handler.chargePoint.SendRequestAsync(
				request,
				callback,
			)
		},
		retry.Attempts(uint(maxRetries)),
		retry.Delay(time.Duration(retryInterval)),
	)
	if err != nil {
		log.Println("max retries reached")
	}
	return err
}

// sendHeartBeat Send a setHeartbeat to the central system.
func (handler *ChargePointHandler) sendHeartBeat() error {
	return handler.SendRequest(
		core.HeartbeatRequest{},
		func(confirmation ocpp.Response, protoError error) {
			log.Printf("Sent heartbeat")
		})
}

// startCharging Start charging on the first available Connector. If there is no available Connector, reject the request.
func (handler *ChargePointHandler) startCharging(tagId string) error {
	var connector = handler.FindAvailableConnector()
	if connector != nil {
		return handler.startChargingConnector(connector, tagId)
	}
	return errors.New("no available connectors")
}

// startChargingConnector Start charging a connector with the specified ID. Send the request to the central system, turn on the Connector,
// update the status of the Connector, and start the maxChargingTime timer and sample the PowerMeter, if it's enabled.
func (handler *ChargePointHandler) startChargingConnector(connector *Connector, tagId string) error {
	if connector == nil {
		return fmt.Errorf("connector is nil")
	}
	if !connector.IsAvailable() || !handler.IsAvailable {
		return errors.New("connector or cp unavailable")
	}

	if handler.isTagAuthorized(tagId) {
		request := core.StartTransactionRequest{
			ConnectorId: connector.ConnectorId,
			IdTag:       tagId,
			Timestamp:   &types2.DateTime{Time: time.Now()},
			MeterStart:  0,
		}
		callback := func(confirmation ocpp.Response, protoError error) {
			startTransactionConf := confirmation.(*core.StartTransactionConfirmation)
			if startTransactionConf.TransactionId > 0 && startTransactionConf.IdTagInfo.Status == types2.AuthorizationStatusAccepted {
				err := connector.StartCharging(strconv.Itoa(startTransactionConf.TransactionId), tagId)
				if err != nil {
					log.Printf("Unable to start charging connector %d: %s", connector.ConnectorId, err)
					return
				}

				log.Printf("Started charging connector %d at %s", connector.ConnectorId, time.Now())
				if connector.PowerMeterEnabled {
					sampleError := connector.preparePowerMeterAtConnector()
					if sampleError != nil {
						log.Printf("Cannot sample connector %d; %v \n", connector.ConnectorId, err)
					}
				}

				// schedule timer to stop the transaction at the time limit
				_, err = scheduler.GetScheduler().Every(connector.MaxChargingTime).Minutes().LimitRunsTo(1).
					Tag(fmt.Sprintf("connector%dTimer", connector.ConnectorId)).Do(handler.stopChargingConnector, connector, core.ReasonOther)
				if err != nil {
					fmt.Println("cannot schedule stop charging:", err)
				}
			} else {
				log.Printf("Transaction unauthorized at connector %d", connector.ConnectorId)
			}
		}
		return handler.SendRequest(request, callback)
	}
	return errors.New("card unauthorized")
}

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the Connector and calculate the energy consumed.
func (handler *ChargePointHandler) stopChargingConnector(connector *Connector, reason core.Reason) error {
	if connector == nil {
		return fmt.Errorf("connector pointer is nil")
	}

	if connector.IsCharging() || connector.IsPreparing() {
		stopTransactionOnEVDisconnect, err := settings.GetConfigurationValue("StopTransactionOnEVSideDisconnect")
		if err != nil {
			return err
		}

		if stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
			err := connector.StopCharging(reason)
			return err
		}

		transactionId, err := strconv.Atoi(connector.session.TransactionId)
		if err != nil {
			return err
		}

		request := core.StopTransactionRequest{
			TransactionId: transactionId,
			IdTag:         "",
			MeterStop:     int(connector.session.CalculateEnergyConsumptionWithAvgPower()),
			Timestamp:     &types2.DateTime{Time: time.Now()},
			Reason:        reason,
		}
		var callback = func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				log.Printf("Server responded with error for stopping a transaction at %d: %s", connector.ConnectorId, err)
				return
			}

			log.Println("Stopping transaction at ", connector.ConnectorId)
			err := connector.StopCharging(reason)
			if err != nil {
				log.Printf("Unable to stop charging connector %d: %s", connector.ConnectorId, err)
				return
			}

			_ = scheduler.GetScheduler().RemoveByTag(fmt.Sprintf("connector%dSampling", connector.ConnectorId))
			err = scheduler.GetScheduler().RemoveByTag(fmt.Sprintf("connector%dTimer", connector.ConnectorId))
			log.Printf("Stopped charging connector %d at %s", connector.ConnectorId, time.Now())
		}
		return handler.SendRequest(request, callback)
	}
	return errors.New("connector not charging")
}

// stopChargingConnectorWithTagId Search for a Connector that contains the tagId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var connector = handler.FindConnectorWithTagId(tagId)
	if connector != nil {
		return handler.stopChargingConnector(connector, reason)
	}
	return errors.New("no connector with tag id")
}

// stopChargingConnectorWithTransactionId Search for a Connector that contains the transactionId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTransactionId(transactionId string) error {
	var connector = handler.FindConnectorWithTransactionId(transactionId)
	if connector != nil {
		return handler.stopChargingConnector(connector, core.ReasonRemote)
	}
	return errors.New("no connector with transaction id")
}
