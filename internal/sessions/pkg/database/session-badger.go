package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	session "github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models"
)

type SessionBadgerDb struct {
	db     *badger.DB
	logger log.FieldLogger
}

func NewSessionBadgerDb(db *badger.DB) *SessionBadgerDb {
	return &SessionBadgerDb{
		db:     db,
		logger: log.WithField("component", "session-database"),
	}
}

func getSessionTransactionIdKey(transactionId string) []byte {
	return []byte(fmt.Sprintf("session-transaction-%s", transactionId))
}

func (s *SessionBadgerDb) CreateSession(session *session.Session) error {
	s.logger.WithField("transactionId", session.TransactionId).Info("Creating session")
	return s.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(session)
		if err != nil {
			return err
		}

		return txn.Set(getSessionTransactionIdKey(session.TransactionId), marshal)
	})
}

func (s *SessionBadgerDb) StopSession(transactionId string) error {
	s.logger.WithField("transactionId", transactionId).Info("Stopping a session")

	return s.db.Update(func(txn *badger.Txn) error {
		s, err := txn.Get(getSessionTransactionIdKey(transactionId))
		if err != nil {
			return err
		}
		var dSession session.Session

		err = s.Value(func(val []byte) error {
			return json.Unmarshal(val, &dSession)
		})
		if err != nil {
			return err
		}

		// End the session
		dSession.EndSession()

		marshal, err := json.Marshal(dSession)
		if err != nil {
			return err
		}

		// Update the session with the new status and set TTL to 24 hours
		return txn.SetEntry(badger.NewEntry(getSessionTransactionIdKey(transactionId), marshal).WithTTL(time.Hour * 24))
	})
}

func (s *SessionBadgerDb) UpdateSession(sesh *session.Session) error {
	s.logger.Info("Updating a session")

	return s.db.Update(func(txn *badger.Txn) error {
		// Verify that it exists
		_, err := txn.Get(getSessionTransactionIdKey(sesh.TransactionId))
		if err != nil {
			return err
		}

		marshal, err := json.Marshal(sesh)
		if err != nil {
			return err
		}

		return txn.Set(getSessionTransactionIdKey(sesh.TransactionId), marshal)
	})
}

func (s *SessionBadgerDb) GetSession(evseId int, connectorId *int) (*session.Session, error) {
	s.logger.WithField("evseId", evseId).Info("Getting session")

	return nil, nil
}

func (s *SessionBadgerDb) GetSessions() ([]session.Session, error) {
	s.logger.Info("Getting sessions")

	return nil, nil
}

func (s *SessionBadgerDb) GetActiveSessions() ([]session.Session, error) {
	s.logger.Info("Getting sessions")

	return nil, nil
}

func (s *SessionBadgerDb) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	s.logger.WithField("transactionId", transactionId).Info("Getting session with transaction id")

	return nil, nil
}

func (s *SessionBadgerDb) GetSessionWithTagId(tagId string) (*session.Session, error) {
	s.logger.WithField("tagId", tagId).Info("Getting session with tag id")

	return nil, nil
}
