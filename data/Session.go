package data

import (
	strUtil "github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"time"
)

type Session struct {
	IsActive      bool
	TransactionId string
	TagId         string
	Started       string
	Consumption   []types.MeterValue
}

// StartSession Starts the Session, storing the transactionId and tagId of the user.
// Checks if transaction and tagId are valid strings.
func (session *Session) StartSession(transactionId string, tagId string) bool {
	if !session.IsActive && strUtil.IsAlphanumeric(transactionId) && strUtil.IsAlphanumeric(tagId) {
		session.TransactionId = transactionId
		session.TagId = tagId
		session.IsActive = true
		session.Started = time.Now().Format(time.RFC3339)
		session.Consumption = []types.MeterValue{}
		return true
	}
	return false
}

//EndSession End the Session if one is active. Reset the attributes, except the measurands.
func (session *Session) EndSession() {
	if session.IsActive {
		session.TransactionId = ""
		session.TagId = ""
		session.IsActive = false
		session.Started = ""
	}
}

// AddSampledValue Add all the samples taken to the Session.
func (session *Session) AddSampledValue(samples []types.SampledValue) {
	if session.IsActive {
		session.Consumption = append(session.Consumption, types.MeterValue{SampledValue: samples})
	}
}
