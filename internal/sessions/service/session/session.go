package session

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models/session"
)

type (
	Manager interface {
		StartSession(evseId int, connectorId *int, tagId, transactionId string) error
		StopSession(transactionId string) error
		UpdateMeterValues(transactionId string, values ...types.SampledValue) error
		GetSession(evseId int, connectorId *int) (*session.Session, error)
		GetSessionWithTransactionId(transactionId string) (*session.Session, error)
		GetSessionWithTagId(tagId string) (*session.Session, error)
	}

	Impl struct {
		sessionRepository database.SessionRepository
	}
)

func NewSessionManager(sessionRepository database.SessionRepository) *Impl {
	return &Impl{
		sessionRepository: sessionRepository,
	}
}

func (i *Impl) StartSession(evseId int, connectorId *int, tagId, transactionId string) error {
	// Check if a session already exists
	_, err := i.GetSession(evseId, connectorId)
	if err == nil {
		return nil
	}

	// Create a new session
	newSession := session.NewEmptySession()
	err = newSession.StartSession(transactionId, tagId)
	if err != nil {
		return err
	}

	return i.sessionRepository.CreateSession(newSession)
}

func (i *Impl) StopSession(transactionId string) error {
	return i.sessionRepository.StopSession(transactionId)
}

func (i *Impl) UpdateMeterValues(transactionId string, values ...types.SampledValue) error {
	sesh, err := i.GetSessionWithTransactionId(transactionId)
	if err != nil {
		return nil
	}

	sesh.AddSampledValue(values)

	return i.sessionRepository.UpdateSession(sesh)
}

func (i *Impl) GetSession(evseId int, connectorId *int) (*session.Session, error) {
	return i.sessionRepository.GetSession(evseId, connectorId)
}

func (i *Impl) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	return i.sessionRepository.GetSessionWithTransactionId(transactionId)
}

func (i *Impl) GetSessionWithTagId(tagId string) (*session.Session, error) {
	return i.sessionRepository.GetSessionWithTagId(tagId)
}
