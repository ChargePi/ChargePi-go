package session

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models/session"
)

type (
	Manager interface {
		StartSession(evseId int, connectorId *int, tagId, transactionId string) error
		StopSession(evseId int, connectorId *int, reason core.Reason) error
		UpdateMeterValues(evseId int, connectorId *int, values ...types.SampledValue) error
		GetSession(evseId int, connectorId *int) (*session.Session, error)
		GetSessionWithTransactionId(transactionId string) (*session.Session, error)
		GetSessionWithTagId(tagId string) (*session.Session, error)
	}

	Impl struct {
		database.SessionRepository
	}
)

func (i *Impl) StartSession(evseId int, connectorId *int, tagId, transactionId string) error {
	// TODO implement me
	panic("implement me")
}

func (i *Impl) StopSession(evseId int, connectorId *int, reason core.Reason) error {
	// TODO implement me
	panic("implement me")
}

func (i *Impl) UpdateMeterValues(evseId int, connectorId *int, values ...types.SampledValue) error {
	// TODO implement me
	panic("implement me")
}

func (i *Impl) GetSession(evseId int, connectorId *int) (*session.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (i *Impl) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	// TODO implement me
	panic("implement me")
}

func (i *Impl) GetSessionWithTagId(tagId string) (*session.Session, error) {
	// TODO implement me
	panic("implement me")
}
