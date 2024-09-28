package diagnostics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ftpUploaderTestSuite struct {
	suite.Suite
	ftp *ftpUploader
}

func (s *ftpUploaderTestSuite) SetupSuite() {
	s.ftp = newFtpUploader(FtpUploaderSettings{})
}

func (s *ftpUploaderTestSuite) SetupTest() {

}

func (s *ftpUploaderTestSuite) TearDownSuite() {
}

func (s *ftpUploaderTestSuite) TestUpload() {
	tests := []struct {
		name     string
		stream   []byte
		location string
		err      error
	}{
		{
			name:     "Successful upload",
			stream:   []byte("test"),
			location: "ftp://user:password@localhost:21/test",
			err:      nil,
		},
		{
			name:     "Malformed URL",
			stream:   []byte("test"),
			location: "ftp://user:password@localhost::21/test",
			err:      nil,
		},
		{
			name:     "Timed out",
			stream:   []byte("test"),
			location: "ftp://user:password@localhost:21/test",
			err:      nil,
		},
		{
			name:     "Cannot dial FTP server",
			stream:   []byte("test"),
			location: "ftp://user:password@localhost:21/test",
			err:      nil,
		},
		{
			name:     "Cannot login to FTP server",
			stream:   []byte("test"),
			location: "ftp://user1:password1@localhost:21/test",
			err:      nil,
		},
		{
			name:     "Unable to store file",
			stream:   []byte("test"),
			location: "ftp://user:password@localhost::21/test",
			err:      nil,
		},
	}

	// todo mock FTP server
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := s.ftp.Upload(ctx, tt.stream, tt.location)
			if tt.err != nil {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestUploader(t *testing.T) {
	suite.Run(t, new(ftpUploaderTestSuite))
}
