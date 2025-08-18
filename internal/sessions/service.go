package sessions

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/tavsec/gin-healthcheck/checks"
	"go.uber.org/zap"

	session "github.com/ChargePi/ChargePi-go/internal/sessions/models"
)

type (
	Service interface {
		checks.Check
		// StartSession starts a new session for the given EVSE and tagId. If a session already exists, it does nothing.
		StartSession(ctx context.Context, evseId int, connectorId *int, tagId, transactionId string) error
		AddTransactionIdToSession(ctx context.Context, evseId int, connectorId *int, transactionId string) error
		StopSession(ctx context.Context, evseId int, connectorId *int, tagId, transactionId *string) error
		UpdateMeterValues(ctx context.Context, transactionId string, values ...types.SampledValue) error
		GetSession(ctx context.Context, evseId int, connectorId *int) (*session.Session, error)
		GetSessions(ctx context.Context, evseId *int) ([]session.Session, error)
		GetSessionWithTransactionId(ctx context.Context, transactionId string) (*session.Session, error)
		GetSessionWithTagId(ctx context.Context, tagId string) (*session.Session, error)
	}

	Impl struct {
		sessionRepository SessionRepository
		logger            *zap.Logger
	}
)

func NewSessionService(logger *zap.Logger, sessionRepository SessionRepository) (*Impl, error) {
	return &Impl{
		sessionRepository: sessionRepository,
		logger:            logger.Named("session_service"),
	}, nil
}

func (i *Impl) StartSession(ctx context.Context, evseId int, connectorId *int, tagId, transactionId string) error {
	i.logger.With(
		zap.Int("evseId", evseId),
		zap.Intp("connectorId", connectorId),
		zap.String("tagId", tagId),
		zap.String("transactionId", transactionId),
	).Info("Starting session")

	// Check if a session already exists
	_, err := i.GetSession(ctx, evseId, connectorId)
	if err == nil {
		return nil
	}

	// Create a new session
	newSession := session.NewEmptySession()
	err = newSession.StartSession(tagId)
	if err != nil {
		return err
	}

	return i.sessionRepository.CreateSession(ctx, newSession)
}

func (i *Impl) StopSession(ctx context.Context, evseId int, connectorId *int, tagId, transactionId *string) error {
	i.logger.With(zap.Stringp("transaction_id", transactionId)).Info("Stopping a session")

	return i.sessionRepository.StopSession(ctx, *transactionId)
}

func (i *Impl) UpdateMeterValues(ctx context.Context, transactionId string, values ...types.SampledValue) error {
	i.logger.With(zap.String("transaction", transactionId)).Info("Updating meter values")

	sessionWithTransactionId, err := i.GetSessionWithTransactionId(ctx, transactionId)
	if err != nil {
		return err
	}

	err = sessionWithTransactionId.AddSampledValue(values)
	if err != nil {
		return err
	}

	return i.sessionRepository.UpdateSession(ctx, sessionWithTransactionId)
}

func (i *Impl) GetSession(ctx context.Context, evseId int, connectorId *int) (*session.Session, error) {
	i.logger.With(
		zap.Int("evseId", evseId),
		zap.Intp("connectorId", connectorId),
	).Info("Getting session")

	return i.sessionRepository.GetSession(ctx, evseId, connectorId)
}

func (i *Impl) GetSessionWithTransactionId(ctx context.Context, transactionId string) (*session.Session, error) {
	i.logger.With(
		zap.String("transactionId", transactionId),
	).Info("Getting session with transaction id")

	return i.sessionRepository.GetSessionWithTransactionId(ctx, transactionId)
}

func (i *Impl) GetSessionWithTagId(ctx context.Context, tagId string) (*session.Session, error) {
	i.logger.With(
		zap.String("tagId", tagId),
	).Info("Getting session with tag id")
	return i.sessionRepository.GetSessionWithTagId(ctx, tagId)
}

func (i *Impl) AddTransactionIdToSession(ctx context.Context, evseId int, connectorId *int, transactionId string) error {
	i.logger.With(
		zap.String("transactionId", transactionId),
		zap.Int("evseId", evseId),
	).Info("Adding transaction id to session")
	return nil
}

func (i *Impl) GetSessions(ctx context.Context, evseId *int) ([]session.Session, error) {
	i.logger.With(
		zap.Intp("evseId", evseId),
	).Info("Getting sessions")
	return nil, nil
}

func (i *Impl) Pass() bool {
	return true
}

func (i *Impl) Name() string {
	return "session-service"
}
