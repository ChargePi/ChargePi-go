package diagnostics

import (
	"bytes"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/observability/logging"
)

type Manager interface {
	GetLogs() ([]byte, error)
	GetLogsByDate(startDate, stopDate *time.Time) ([]byte, error)
	UploadLogs(url string, startDate, stopDate *time.Time) error
}

type ManagerImpl struct {
	logger       log.FieldLogger
	uploadStatus string
}

func NewManager() *ManagerImpl {
	return &ManagerImpl{
		logger: log.WithField("component", "diagnostics-manager"),
	}
}

// GetLogs gets the logs of the application.
func (m *ManagerImpl) GetLogs() ([]byte, error) {
	m.logger.Debug("Getting logs")

	files, err := os.ReadDir(logging.LogFileDir)
	if err != nil {
		return nil, err
	}

	// Search for files with specific name
	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		return file.Name(), strings.HasPrefix(file.Name(), logging.LogFileName)
	})

	logs, err := m.joinLogs(names)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLogsByDate gets the logs of the application from the given date range.
func (m *ManagerImpl) GetLogsByDate(startDate, stopDate *time.Time) ([]byte, error) {
	m.logger.Debug("Getting logs by date")

	files, err := os.ReadDir(logging.LogFileDir)
	if err != nil {
		return nil, err
	}

	// Search for files with specific name
	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		info, err := file.Info()
		if err != nil {
			return "", false
		}

		return file.Name(), strings.HasPrefix(file.Name(), logging.LogFileName) && info.ModTime().After(*startDate) && info.ModTime().Before(*stopDate)
	})

	logs, err := m.joinLogs(names)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// UploadLogs uploads the logs to the given FTP server.
func (m *ManagerImpl) UploadLogs(location string, startDate, stopDate *time.Time) error {
	m.logger.Debug("Uploading logs to FTP server")

	logStream, err := m.GetLogsByDate(startDate, stopDate)
	if err != nil {
		return err
	}

	parse, err := url.Parse(location)
	if err != nil {
		return err
	}

	// Connect to the FTP server.
	connection, err := ftp.Dial(parse.Hostname(), ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return err
	}

	pass, hasPassword := parse.User.Password()
	if !hasPassword {
		pass = ""
	}

	// Login to the FTP server.
	err = connection.Login(parse.User.Username(), pass)
	if err != nil {
		return err
	}

	// Upload the local file to the remote FTP server.
	return connection.Stor(parse.Path, bytes.NewReader(logStream))
}

// joinLogs joins the logs from the given files into one byte array.
func (m *ManagerImpl) joinLogs(files []string) ([]byte, error) {
	var combinedLogs []byte

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		in.Close()

	}

	return combinedLogs, nil
}
