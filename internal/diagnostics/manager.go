package diagnostics

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	log "github.com/sirupsen/logrus"
)

const logFileName = "chargepi.log"
const logFileDir = "/var/log/chargepi"

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

func (m *ManagerImpl) GetLogs() ([]byte, error) {
	m.logger.Debug("Getting logs")

	files, err := os.ReadDir(logFileDir)
	if err != nil {
		return nil, err
	}

	names := []string{}

	// Search for files with specific name
	for _, file := range files {
		if strings.HasPrefix(file.Name(), logFileName) {
			names = append(names, file.Name())
		}
	}

	return nil, nil
}

func (m *ManagerImpl) GetLogsByDate(startDate, stopDate *time.Time) ([]byte, error) {
	m.logger.Debug("Getting logs by date")

	files, err := os.ReadDir(logFileDir)
	if err != nil {
		return nil, err
	}

	names := []string{}

	// Search for files with specific name and between specific dates
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(info.Name(), logFileName) && info.ModTime().After(*startDate) && info.ModTime().Before(*stopDate) {
			names = append(names, file.Name())
		}

	}

	return nil, nil
}

func (m *ManagerImpl) UploadLogs(url string, startDate, stopDate *time.Time) error {
	m.logger.Debug("Uploading logs to FTP server")
	// todo filter logs by date

	err := m.joinLogs([]string{}, logFileName)
	if err != nil {
		return err
	}

	settings, err := m.getFtpUploadSettings(url)
	if err != nil {
		return err
	}

	// Connect to the FTP server.
	connection, err := ftp.Dial(settings["address"], ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return err
	}

	// Login to the FTP server.
	err = connection.Login(settings["username"], settings["password"])
	if err != nil {
		return err
	}

	// Open the local file.
	file, err := os.Open(fmt.Sprintf("%s/%s", logFileDir, logFileName))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Upload the local file to the remote FTP server.
	return connection.Stor(settings["path"], file)
}

func (m *ManagerImpl) joinLogs(files []string, destFile string) error {
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			// Remove the output file on error
			err := os.Remove(destFile)
			if err != nil {
				return err
			}

			return err
		}

		in.Close()

	}

	return nil
}

func (m *ManagerImpl) getFtpUploadSettings(url string) (map[string]string, error) {

	return nil, nil
}
