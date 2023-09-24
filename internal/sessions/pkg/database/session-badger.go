package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	session "github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models"
)

type SessionBadgerDb struct {
	db *badger.DB
}

func NewSessionBadgerDb(db *badger.DB) *SessionBadgerDb {
	return &SessionBadgerDb{db: db}
}

func getSessionTransactionIdKey(transactionId string) []byte {
	return []byte(fmt.Sprintf("session-transaction-%s", transactionId))
}

func (s *SessionBadgerDb) CreateSession(session *session.Session) error {
	return s.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(session)
		if err != nil {
			return err
		}

		return txn.Set(getSessionTransactionIdKey(session.TransactionId), marshal)
	})
}

func (s *SessionBadgerDb) StopSession(transactionId string) error {
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
	return nil, nil
}

func (s *SessionBadgerDb) GetSessions() ([]session.Session, error) {
	return nil, nil
}

func (s *SessionBadgerDb) GetActiveSessions() ([]session.Session, error) {
	return nil, nil
}

func (s *SessionBadgerDb) GetSessionWithTransactionId(transactionId string) (*session.Session, error) {
	return nil, nil
}

func (s *SessionBadgerDb) GetSessionWithTagId(tagId string) (*session.Session, error) {
	return nil, nil
}
