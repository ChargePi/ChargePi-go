package sessions

import (
	session "github.com/ChargePi/ChargePi-go/internal/sessions/models"
)

type SessionRepository interface {
	CreateSession(session *session.Session) error
	StopSession(transactionId string) error
	UpdateSession(session *session.Session) error
	GetSession(evseId int, connectorId *int) (*session.Session, error)
	GetSessions() ([]session.Session, error)
	GetActiveSessions() ([]session.Session, error)
	GetSessionWithTransactionId(transactionId string) (*session.Session, error)
	GetSessionWithTagId(tagId string) (*session.Session, error)
}
