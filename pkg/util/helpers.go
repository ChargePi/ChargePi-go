package util

import (
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	log "github.com/sirupsen/logrus"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"reflect"
	"strconv"
	"time"
)

func IsNilInterfaceOrPointer(sth interface{}) bool {
	return sth == nil || (reflect.ValueOf(sth).Kind() == reflect.Ptr && reflect.ValueOf(sth).IsNil())
}

func HandleRequestErr(err error, text string) {
	if err != nil {
		log.WithError(err).Errorf(text)
	}
}

// SendRequest is a middleware function that implements a retry mechanism for sending requests. If the max attempts is reached, return an error
func SendRequest(chargePoint ocpp16.ChargePoint, request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error {
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
