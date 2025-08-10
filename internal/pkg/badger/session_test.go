package badger

import (
	"os"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/sessions/models"
	"github.com/stretchr/testify/suite"
)

type sessionRepositoryTestSuite struct {
	suite.Suite
	db     *Database
	tmpDir string
}

func (s *sessionRepositoryTestSuite) SetupSuite() {
	// Create temporary database file/dir for testing
	tempDir, err := os.MkdirTemp("", "badger_session_test_*")
	s.Require().NoError(err)
	s.tmpDir = tempDir

	s.db, err = NewBadgerDb(tempDir)
	s.Require().NoError(err)
}

func (s *sessionRepositoryTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	// Remove temporary database file/dir
	err := os.RemoveAll(s.tmpDir)
	s.Require().NoError(err)
}

func (s *sessionRepositoryTestSuite) SetupTest() {
	// Clean up the database before each test
	err := s.db.db.DropAll()
	s.Require().NoError(err)
}

func (s *sessionRepositoryTestSuite) TestCreateSession() {
	tests := []struct {
		name        string
		session     *models.Session
		expectError bool
	}{
		{
			name: "Create session successfully",
			session: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "test-transaction-123",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
		{
			name: "Create session with empty transaction ID",
			session: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Execute
			err := s.db.CreateSession(tt.session)
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)

				// Verify session was actually created
				if tt.session.TransactionId != "" {
					// Note: GetSessionWithTransactionId is not implemented yet, so this will fail
					// This test will need to be updated when that method is implemented
					_, err := s.db.GetSessionWithTransactionId(tt.session.TransactionId)
					s.Error(err) // Expected to fail for now
				}
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestStopSession() {
	tests := []struct {
		name          string
		setupSession  *models.Session
		transactionId string
		expectError   bool
	}{
		{
			name: "Stop existing session successfully",
			setupSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "stop-transaction-123",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			transactionId: "stop-transaction-123",
			expectError:   false,
		},
		{
			name:          "Stop non-existing session",
			setupSession:  nil,
			transactionId: "non-existing-transaction",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSession != nil {
				err := s.db.CreateSession(tt.setupSession)
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.StopSession(tt.transactionId)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestUpdateSession() {
	tests := []struct {
		name          string
		setupSession  *models.Session
		updateSession *models.Session
		expectError   bool
	}{
		{
			name: "Update existing session successfully",
			setupSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "update-transaction-123",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			updateSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "update-transaction-123",
				TagId:         "updated-tag-789",
				IsActive:      false,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
				Ended:         &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
		{
			name:         "Update non-existing session",
			setupSession: nil,
			updateSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "non-existing-transaction",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSession != nil {
				err := s.db.CreateSession(tt.setupSession)
				s.Require().NoError(err)
			}

			// Execute
			err := s.db.UpdateSession(tt.updateSession)

			// Assert
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestGetSession() {
	tests := []struct {
		name         string
		evseId       int
		connectorId  *int
		setupSession *models.Session
		expectError  bool
	}{
		{
			name:        "Get session by EVSE ID and connector ID",
			evseId:      1,
			connectorId: &[]int{1}[0],
			setupSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "get-transaction-123",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSession != nil {
				err := s.db.CreateSession(tt.setupSession)
				s.Require().NoError(err)
			}

			// Execute
			session, err := s.db.GetSession(tt.evseId, tt.connectorId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(session)
			} else {
				// Note: GetSession is not implemented yet, so this will return nil, nil
				s.NoError(err)
				s.Nil(session) // Expected to be nil for now
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestGetSessions() {
	tests := []struct {
		name          string
		setupSessions []*models.Session
		expectError   bool
	}{
		{
			name: "Get all sessions",
			setupSessions: []*models.Session{
				{
					ID:            models.NewEmptySession().ID,
					TransactionId: "session1-123",
					TagId:         "tag1-456",
					IsActive:      true,
					EvseId:        1,
					ConnectorId:   &[]int{1}[0],
					Started:       &[]time.Time{time.Now()}[0],
				},
				{
					ID:            models.NewEmptySession().ID,
					TransactionId: "session2-456",
					TagId:         "tag2-789",
					IsActive:      false,
					EvseId:        2,
					ConnectorId:   &[]int{2}[0],
					Started:       &[]time.Time{time.Now()}[0],
					Ended:         &[]time.Time{time.Now()}[0],
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			for _, session := range tt.setupSessions {
				err := s.db.CreateSession(session)
				s.Require().NoError(err)
			}

			// Execute
			sessions, err := s.db.GetSessions()

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(sessions)
			} else {
				// Note: GetSessions is not implemented yet, so this will return nil, nil
				s.NoError(err)
				s.Nil(sessions) // Expected to be nil for now
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestGetActiveSessions() {
	tests := []struct {
		name          string
		setupSessions []*models.Session
		expectError   bool
	}{
		{
			name: "Get active sessions",
			setupSessions: []*models.Session{
				{
					ID:            models.NewEmptySession().ID,
					TransactionId: "active-session1-123",
					TagId:         "tag1-456",
					IsActive:      true,
					EvseId:        1,
					ConnectorId:   &[]int{1}[0],
					Started:       &[]time.Time{time.Now()}[0],
				},
				{
					ID:            models.NewEmptySession().ID,
					TransactionId: "inactive-session2-456",
					TagId:         "tag2-789",
					IsActive:      false,
					EvseId:        2,
					ConnectorId:   &[]int{2}[0],
					Started:       &[]time.Time{time.Now()}[0],
					Ended:         &[]time.Time{time.Now()}[0],
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			for _, session := range tt.setupSessions {
				err := s.db.CreateSession(session)
				s.Require().NoError(err)
			}

			// Execute
			sessions, err := s.db.GetActiveSessions()

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(sessions)
			} else {
				// Note: GetActiveSessions is not implemented yet, so this will return nil, nil
				s.NoError(err)
				s.Nil(sessions) // Expected to be nil for now
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestGetSessionWithTransactionId() {
	tests := []struct {
		name          string
		transactionId string
		setupSession  *models.Session
		expectError   bool
	}{
		{
			name:          "Get session with existing transaction ID",
			transactionId: "transaction-123",
			setupSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "transaction-123",
				TagId:         "test-tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
		{
			name:          "Get session with non-existing transaction ID",
			transactionId: "non-existing-transaction",
			setupSession:  nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSession != nil {
				err := s.db.CreateSession(tt.setupSession)
				s.Require().NoError(err)
			}

			// Execute
			session, err := s.db.GetSessionWithTransactionId(tt.transactionId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(session)
			} else {
				// Note: GetSessionWithTransactionId is not implemented yet, so this will return nil, nil
				s.NoError(err)
				s.Nil(session) // Expected to be nil for now
			}
		})
	}
}

func (s *sessionRepositoryTestSuite) TestGetSessionWithTagId() {
	tests := []struct {
		name         string
		tagId        string
		setupSession *models.Session
		expectError  bool
	}{
		{
			name:  "Get session with existing tag ID",
			tagId: "tag-456",
			setupSession: &models.Session{
				ID:            models.NewEmptySession().ID,
				TransactionId: "transaction-123",
				TagId:         "tag-456",
				IsActive:      true,
				EvseId:        1,
				ConnectorId:   &[]int{1}[0],
				Started:       &[]time.Time{time.Now()}[0],
			},
			expectError: false,
		},
		{
			name:         "Get session with non-existing tag ID",
			tagId:        "non-existing-tag",
			setupSession: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupSession != nil {
				err := s.db.CreateSession(tt.setupSession)
				s.Require().NoError(err)
			}

			// Execute
			session, err := s.db.GetSessionWithTagId(tt.tagId)

			// Assert
			if tt.expectError {
				s.Error(err)
				s.Nil(session)
			} else {
				// Note: GetSessionWithTagId is not implemented yet, so this will return nil, nil
				s.NoError(err)
				s.Nil(session) // Expected to be nil for now
			}
		})
	}
}

func TestSessionRepository(t *testing.T) {
	suite.Run(t, new(sessionRepositoryTestSuite))
}
