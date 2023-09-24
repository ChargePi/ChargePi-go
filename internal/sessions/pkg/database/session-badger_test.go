package database

// Create a test for the session repository

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	session "github.com/xBlaz3kx/ChargePi-go/internal/sessions/pkg/models"
)

func TestSessionRepository(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a session repository
	repo := NewSessionBadgerDb(db)

	// Create a session
	newSession := session.NewEmptySession()
	err = newSession.StartSession("1", "tag-id")
	if err != nil {
		t.Fatal(err)
	}

	// Create the session
	err = repo.CreateSession(newSession)
	if err != nil {
		t.Fatal(err)
	}

	// Get the session
	s, err := repo.GetSessionWithTransactionId("1")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the session is the same
	assert.Equal(t, newSession, s, "session should be the same")

	// End the session
	err = repo.StopSession("1")
	if err != nil {
		t.Fatal(err)
	}

	// Get the session
	s, err = repo.GetSessionWithTransactionId("transaction-id")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the session is ended
	assert.True(t, s.IsActive, "session should be ended")
}
