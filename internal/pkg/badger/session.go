package badger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/dgraph-io/badger/v3"

	session "github.com/ChargePi/ChargePi-go/internal/sessions/models"
)

const (
	sessionTransactionKey = "session-transaction"
)

func getSessionTransactionIdKey(transactionId string) []byte {
	return []byte(fmt.Sprintf("%s-%s", sessionTransactionKey, transactionId))
}

func (db *Database) CreateSession(ctx context.Context, session *session.Session) error {
	// todo handle cases where a session does not have a transaction ID yet
	db.logger.With(zap.Any("session", session)).Info("Creating session")
	return db.db.Update(func(txn *badger.Txn) error {
		// Check if already exists
		_, err := txn.Get(getSessionTransactionIdKey(session.TransactionId))
		if err != nil {
			return err
		}

		marshal, err := json.Marshal(session)
		if err != nil {
			return err
		}

		return txn.Set(getSessionTransactionIdKey(session.TransactionId), marshal)
	})
}

func (db *Database) StopSession(ctx context.Context, transactionId string) error {
	db.logger.With(zap.String("transactionId", transactionId)).Info("Stopping a session")
	// todo handle cases where a session does not have a transaction ID yet
	return db.db.Update(func(txn *badger.Txn) error {
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
		err = dSession.EndSession()
		if err != nil {
			return err
		}

		marshal, err := json.Marshal(dSession)
		if err != nil {
			return err
		}

		// Update the session with the new status and set TTL to 24 hours
		// todo keep session in history up to a week, but remove from active sessions after 24 hours
		return txn.SetEntry(badger.NewEntry(getSessionTransactionIdKey(transactionId), marshal).WithTTL(time.Hour * 24))
	})
}

func (db *Database) UpdateSession(ctx context.Context, sesh *session.Session) error {
	db.logger.Info("Updating a session")

	return db.db.Update(func(txn *badger.Txn) error {
		// Verify that the session exists
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

func (db *Database) GetSession(ctx context.Context, evseId int, connectorId *int) (*session.Session, error) {
	db.logger.With(zap.Int("evseId", evseId)).Info("Getting session")

	return nil, nil
}

func (db *Database) GetSessions(ctx context.Context) ([]session.Session, error) {
	db.logger.Info("Getting sessions")

	return nil, nil
}

func (db *Database) GetActiveSessions(ctx context.Context) ([]session.Session, error) {
	db.logger.Info("Getting sessions")

	return nil, nil
}

func (db *Database) GetSessionWithTransactionId(ctx context.Context, transactionId string) (*session.Session, error) {
	db.logger.With(zap.String("transactionId", transactionId)).Info("Getting session with transaction id")

	return nil, nil
}

func (db *Database) GetSessionWithTagId(ctx context.Context, tagId string) (*session.Session, error) {
	db.logger.With(zap.String("tagId", tagId)).Info("Getting session with tag id")

	return nil, nil
}
