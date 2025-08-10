package session

import (
	"github.com/ChargePi/ChargePi-go/internal/sessions/pkg/database"
	session "github.com/ChargePi/ChargePi-go/internal/sessions/pkg/models"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.uber.org/zap"
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
		logger            *zap.Logger
	}
)

func NewSessionManager(logger *zap.Logger, sessionRepository database.SessionRepository) *Impl {
	return &Impl{
		sessionRepository: sessionRepository,
		logger:            logger.Named("session_service"),
	}
}

func (i *Impl) StartSession(evseId int, connectorId *int, tagId, transactionId string) error {
	i.logger.With(
		zap.Int("evseId", evseId),
		zap.Intp("connectorId", connectorId),
		zap.String("tagId", tagId),
		zap.String("transactionId", transactionId),
	).Info("Starting session")

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
	i.logger.With(zap.String("transaction_id", transactionId)).Info("Stopping a session")
	return i.sessionRepository.StopSession(transactionId)
}

func (i *Impl) UpdateMeterValues(transactionId string, values ...types.SampledValue) error {
	i.logger.With(zap.String("transaction", transactionId)).Info("Updating meter values")

	sessionWithTransactionId, err := i.GetSessionWithTransactionId(transactionId)
	if err != nil {
		return err
	}

	sessionWithTransactionId.AddSampledValue(values)

	return i.sessionRepository.UpdateSession(sessionWithTransactionId)
}

func (i *Impl) GetSession(evseId int, connectorId *int) (*session.Session, error) {
	i.logger.With(
		zap.Int("evseId", evseId),
		zap.Intp("connectorId", connectorId),
	).Info("Getting session")

	return i.sessionRepository.GetSession(evseId, connectorId)
}

func (i *Impl) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	i.logger.With(
		zap.String("transactionId", transactionId),
	).Info("Getting session with transaction id")

	return i.sessionRepository.GetSessionWithTransactionId(transactionId)
}

func (i *Impl) GetSessionWithTagId(tagId string) (*session.Session, error) {
	i.logger.With(
		zap.String("tagId", tagId),
	).Info("Getting session with tag id")
	return i.sessionRepository.GetSessionWithTagId(tagId)
}
