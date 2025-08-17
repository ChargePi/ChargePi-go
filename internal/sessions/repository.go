package sessions

import (
	"context"

	session "github.com/ChargePi/ChargePi-go/internal/sessions/models"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *session.Session) error
	StopSession(ctx context.Context, transactionId string) error
	UpdateSession(ctx context.Context, session *session.Session) error
	GetSession(ctx context.Context, evseId int, connectorId *int) (*session.Session, error)
	GetSessions(ctx context.Context) ([]session.Session, error)
	GetActiveSessions(ctx context.Context) ([]session.Session, error)
	GetSessionWithTransactionId(ctx context.Context, transactionId string) (*session.Session, error)
	GetSessionWithTagId(ctx context.Context, tagId string) (*session.Session, error)
}
