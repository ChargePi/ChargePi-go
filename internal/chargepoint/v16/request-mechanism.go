package v16

import (
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	log "github.com/sirupsen/logrus"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"strconv"
	"time"
)

// sendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func sendRequest(chargePoint ocpp16.ChargePoint, request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
	var (
		maxMessageAttempts, attemptErr  = ocppConfigManager.GetConfigurationValue(v16.TransactionMessageAttempts.String())
		retryIntervalValue, intervalErr = ocppConfigManager.GetConfigurationValue(v16.TransactionMessageRetryInterval.String())
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
			return chargePoint.SendRequestAsync(
				request,
				callback,
			)
		},
		retry.Attempts(uint(maxRetries)),
		retry.Delay(time.Duration(retryInterval)),
	)
}

func handleRequestErr(err error, text string) {
	if err != nil {
		log.WithError(err).Errorf(text)
	}
}
