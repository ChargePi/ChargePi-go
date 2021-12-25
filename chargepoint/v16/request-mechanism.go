package v16

import (
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
	"strconv"
	"time"
)

// SendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func (handler *ChargePointHandler) SendRequest(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	var (
		maxMessageAttempts, attemptErr  = conf_manager.GetConfigurationValue("TransactionMessageAttempts")
		retryIntervalValue, intervalErr = conf_manager.GetConfigurationValue("TransactionMessageRetryInterval")
		maxRetries, convError           = strconv.Atoi(maxMessageAttempts)
		retryInterval, convError2       = strconv.Atoi(retryIntervalValue)
	)

	if attemptErr != nil || convError != nil {
		maxRetries = 5
	}

	if intervalErr != nil || convError2 != nil {
		retryInterval = 30
	}

	return retry.Do(
		func() error {
			return handler.chargePoint.SendRequestAsync(
				request,
				callback,
			)
		},
		retry.Attempts(uint(maxRetries)),
		retry.Delay(time.Duration(retryInterval)),
	)
}
