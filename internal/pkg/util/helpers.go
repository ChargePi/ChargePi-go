package util

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

// IsNilInterfaceOrPointer check if the variable is nil or if the pointer's value is nil.
func IsNilInterfaceOrPointer(sth interface{}) bool {
	return sth == nil || (reflect.ValueOf(sth).Kind() == reflect.Ptr && reflect.ValueOf(sth).IsNil())
}

func HandleRequestErr(err error, text string) {
	if err != nil {
		log.WithError(err).Errorf(text)
	}
}
