package diagnostics

import (
	"archive/zip"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/diagnostics/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	mockUploader *mocks.Mockuploader
	service      *ServiceV1
}

func (s *serviceTestSuite) SetupTest() {
	service, err := NewService()
	s.Require().NoError(err)

	s.service = service
}

func (s *serviceTestSuite) TearDownSuite() {

}

func (s *serviceTestSuite) TestUploadLogs() {
	tests := []struct {
		name          string
		uri           string
		dateFrom      *time.Time
		dateTo        *time.Time
		retryInterval time.Duration
		retryCount    int
		expectErr     bool
	}{
		{
			name: "Successful upload",
		},
		{
			name: "Invalid URL",
		},
		{
			name: "Invalid retry interval",
		},
		{
			name: "Invalid retry count",
		},
		{
			name: "Date to is before date from",
		},
		{
			name: "Date to is after now",
		},
		{
			name: "Date from is after now",
		},
		{
			name: "No logs to upload",
		},
		{
			name: "Failed to upload logs",
		},
		{
			name: "Failed to get logs",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.service.UploadLogs(tt.uri, tt.dateFrom, tt.dateTo, tt.retryInterval, tt.retryCount)
			if tt.expectErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *serviceTestSuite) TestGetLogsByDate() {
	tests := []struct {
		name         string
		dateFrom     *time.Time
		dateTo       *time.Time
		expectedLogs string
		wantErr      bool
	}{
		{
			name: "Successful get logs",
		},
		{
			name: "Date to is before date from",
		},
		{
			name: "Date to is after now",
		},
		{
			name: "Date from is after now",
		},
		{
			name: "No logs",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			logs, err := s.service.GetLogsByDate(tt.dateFrom, tt.dateTo)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(logs)
				s.EqualValues(tt.expectedLogs, string(logs))
			}
		})
	}
}

func (s *serviceTestSuite) TestGetLogsAsZip() {
	tests := []struct {
		name         string
		dateFrom     *time.Time
		dateTo       *time.Time
		expectedLogs string
		wantErr      bool
	}{
		{
			name: "Successful get logs",
		},
		{
			name: "Date to is before date from",
		},
		{
			name: "Date to is after now",
		},
		{
			name: "Date from is after now",
		},
		{
			name: "Unable to zip logs",
		},
		{
			name: "No logs",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			asZip, err := s.service.GetLogsAsZip(tt.dateFrom, tt.dateTo)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(asZip)
				// todo unzip and compare
			}
		})
	}
}

func (s *serviceTestSuite) TestGetLogs() {
	tests := []struct {
		name         string
		expectedLogs string
		wantErr      bool
	}{
		{
			name: "Successful get logs",
		},
		{
			name: "Date to is before date from",
		},
		{
			name: "Date to is after now",
		},
		{
			name: "Date from is after now",
		},
		{
			name: "No logs",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			logs, err := s.service.GetLogs()
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(logs)
				s.EqualValues(tt.expectedLogs, string(logs))
			}
		})
	}
}

func TestManager(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func Test_zipLogs(t *testing.T) {
	f1 := writeToFile(t, []byte("test1"))
	f2 := writeToFile(t, []byte("test2"))
	f3 := writeToFile(t, []byte("test3"))

	defer func() {
		os.Remove(f1)
		os.Remove(f2)
		os.Remove(f3)
	}()

	zipFile, err := zipLogs([]string{f1, f2, f3})
	assert.NoError(t, err)
	assert.NotEmpty(t, zipFile)

	archive, err := zip.OpenReader("")
	if err != nil {
		return
	}
	defer archive.Close()

	fmt.Println(archive.File)
	assert.Len(t, archive.File, 3)
	assert.Equal(t, archive.File[0].Name, f1)
	assert.Equal(t, archive.File[1].Name, f2)
	assert.Equal(t, archive.File[3].Name, f3)
}

func Test_joinLogs(t *testing.T) {
	f1 := writeToFile(t, []byte("test1"))
	f2 := writeToFile(t, []byte("test2"))
	f3 := writeToFile(t, []byte("test3"))

	defer func() {
		os.Remove(f1)
		os.Remove(f2)
		os.Remove(f3)
	}()

	logs, err := joinLogs([]string{f1, f2, f3})
	assert.NoError(t, err)
	assert.Equal(t, string(logs), "test1test2test3")

	// Logs dont exist
	logs, err = joinLogs([]string{"file1", "file2", "file3"})
	assert.Error(t, err)
}

func writeToFile(t *testing.T, content []byte) string {
	file, err := os.CreateTemp("", "log_*.txt")
	require.NoError(t, err)

	_, err = file.Write(content)
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	return file.Name()
}
