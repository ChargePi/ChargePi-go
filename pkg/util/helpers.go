package util

import (
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// IsNilInterfaceOrPointer check if the variable is nil or if the pointer's value is nil.
func IsNilInterfaceOrPointer(sth interface{}) bool {
	return sth == nil || (reflect.ValueOf(sth).Kind() == reflect.Ptr && reflect.ValueOf(sth).IsNil())
}

// GetTypesToSample get the measurands to sample from the OCPP configuration.
func GetTypesToSample() []types.Measurand {
	var (
		measurands []types.Measurand
		// Get the types to sample
		measurandsString, err = ocppConfigManager.GetConfigurationValue(v16.MeterValuesSampledData.String())
	)

	if err != nil {
		return measurands
	}

	for _, measurand := range strings.Split(measurandsString, ",") {
		measurands = append(measurands, types.Measurand(measurand))
	}

	return measurands
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
