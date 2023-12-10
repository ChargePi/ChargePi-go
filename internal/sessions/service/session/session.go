package session

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/database"
	session "github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models"
)

type (
	Manager interface {
		StartSession(evseId int, connectorId *int, tagId, transactionId string) error
		StopSession(transactionId string) error
		UpdateMeterValues(transactionId string, values ...types.SampledValue) error
		GetSession(evseId int, connectorId *int) (*session.Session, error)
		// GetSessions(evseId *int) ([]session.Session, error)
		GetSessionWithTransactionId(transactionId string) (*session.Session, error)
		GetSessionWithTagId(tagId string) (*session.Session, error)
	}

	Impl struct {
		sessionRepository database.SessionRepository
		logger            log.FieldLogger
	}
)

func NewSessionManager(sessionRepository database.SessionRepository) *Impl {
	return &Impl{
		sessionRepository: sessionRepository,
		logger:            log.WithField("component", "session-service"),
	}
}

func (i *Impl) StartSession(evseId int, connectorId *int, tagId, transactionId string) error {
	i.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
		"tagId":       tagId,
		"transaction": transactionId,
	}).Info("Starting session")

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
	i.logger.WithFields(log.Fields{"transaction": transactionId}).Info("Stopping a session")
	return i.sessionRepository.StopSession(transactionId)
}

func (i *Impl) UpdateMeterValues(transactionId string, values ...types.SampledValue) error {
	i.logger.WithField("transaction", transactionId).Info("Updating meter values")

	sessionWithTransactionId, err := i.GetSessionWithTransactionId(transactionId)
	if err != nil {
		return err
	}

	sessionWithTransactionId.AddSampledValue(values)

	return i.sessionRepository.UpdateSession(sessionWithTransactionId)
}

func (i *Impl) GetSession(evseId int, connectorId *int) (*session.Session, error) {
	i.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
	}).Info("Getting session")

	return i.sessionRepository.GetSession(evseId, connectorId)
}

func (i *Impl) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	i.logger.WithFields(log.Fields{
		"transactionId": transactionId,
	}).Info("Getting session with transaction id")

	return i.sessionRepository.GetSessionWithTransactionId(transactionId)
}

func (i *Impl) GetSessionWithTagId(tagId string) (*session.Session, error) {
	i.logger.WithFields(log.Fields{
		"tagId": tagId,
	}).Info("Getting session with tag id")
	return i.sessionRepository.GetSessionWithTagId(tagId)
}
