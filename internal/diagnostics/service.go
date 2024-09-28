package diagnostics

import (
	"archive/zip"
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ChargePi/ChargePi-go/pkg/observability"
	"github.com/avast/retry-go"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/tavsec/gin-healthcheck/checks"
)

type uploader interface {
	Upload(context context.Context, bytes []byte, location string) error
}

type Service interface {
	checks.Check
	GetLogs() ([]byte, error)
	GetLogsByDate(startDate, stopDate *time.Time) ([]byte, error)
	GetLogsAsZip(startDate, stopDate *time.Time) ([]byte, error)
	UploadLogs(url string, startDate, stopDate *time.Time, retryInterval time.Duration, retries int) error
}

type ServiceV1 struct {
	logger       log.FieldLogger
	mu           sync.Mutex
	uploadStatus string
	fac          *uploaderFactory
}

func NewService() (*ServiceV1, error) {
	return &ServiceV1{
		mu:     sync.Mutex{},
		fac:    &uploaderFactory{},
		logger: log.WithField("component", "diagnostics-manager"),
	}, nil
}

// GetLogs gets the logs of the application.
func (m *ServiceV1) GetLogs() ([]byte, error) {
	m.logger.Debug("Getting logs")

	files, err := os.ReadDir(observability.LogFileDir)
	if err != nil {
		return nil, err
	}

	// Search for files with specific name
	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		return file.Name(), strings.HasPrefix(file.Name(), observability.LogFileName)
	})

	logs, err := joinLogs(names)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLogsAsZip gets the logs of the application in a zip file
func (m *ServiceV1) GetLogsAsZip(startDate, stopDate *time.Time) ([]byte, error) {
	m.logger.Debug("Getting logs")

	files, err := os.ReadDir(observability.LogFileDir)
	if err != nil {
		return nil, err
	}

	// Search for files with specific name
	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		info, err := file.Info()
		if err != nil {
			return "", false
		}

		return file.Name(), strings.HasPrefix(file.Name(), observability.LogFileName) && info.ModTime().After(*startDate) && info.ModTime().Before(*stopDate)
	})

	logs, err := zipLogs(names)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLogsByDate gets the logs of the application from the given date range.
func (m *ServiceV1) GetLogsByDate(startDate, stopDate *time.Time) ([]byte, error) {
	m.logger.Debug("Getting logs by date")

	files, err := os.ReadDir(observability.LogFileDir)
	if err != nil {
		return nil, err
	}

	// Search for files with specific name
	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		info, err := file.Info()
		if err != nil {
			return "", false
		}

		return file.Name(), strings.HasPrefix(file.Name(), observability.LogFileName) && info.ModTime().After(*startDate) && info.ModTime().Before(*stopDate)
	})

	logs, err := joinLogs(names)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// UploadLogs uploads the logs to the given FTP server.
func (m *ServiceV1) UploadLogs(uri string, from, to *time.Time, retryInterval time.Duration, retries int) error {
	m.logger.Debug("Uploading logs to FTP server")

	// When uploading, we will zip the logs

	logStream, err := m.GetLogsByDate(from, to)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), retryInterval*time.Duration(retries))
	defer cancel()

	parse, err := url.Parse(uri)
	if err != nil {
		return err
	}

	logUploader, err := m.fac.UploaderFromType(uploaderType(parse.Scheme))
	if err != nil {
		return err
	}

	return retry.Do(
		func() error {
			return logUploader.Upload(ctx, logStream, uri)
		},
		retry.Attempts(uint(retries)),
		retry.Delay(retryInterval),
		retry.OnRetry(func(n uint, err error) {
			m.logger.WithError(err).WithField("attempt", n).Warn("Failed to upload logs, retrying")
		}),
	)
}

func (m *ServiceV1) Pass() bool {
	return true
}

func (m *ServiceV1) Name() string {
	return "diagnostics"
}

// zipLogs zips the files
func zipLogs(files []string) ([]byte, error) {
	temp, err := os.CreateTemp("", "logs_*.zip")
	if err != nil {
		return nil, err
	}

	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.WithError(err).Error("Failed to close temporary file")
		}
	}(temp)

	zipWriter := zip.NewWriter(temp)

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		// Add the file to zip
		writer, err := zipWriter.Create(filepath.Base(file))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(writer, in)
		if err != nil {
			return nil, err
		}

		err = in.Close()
		if err != nil {
			return nil, err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(temp.Name())
}

// joinLogs joins the logs from the given files into one byte array.
func joinLogs(files []string) ([]byte, error) {
	var combinedLogs []byte

	for _, file := range files {
		in, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// Append the file content to the existing byte array
		combinedLogs = append(combinedLogs, in...)
	}

	return combinedLogs, nil
}
