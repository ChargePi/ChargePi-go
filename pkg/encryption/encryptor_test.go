package encryption

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type encryptorTestSuite struct {
	suite.Suite
	encryptor Encryptor
}

func (s *encryptorTestSuite) SetupTest() {
	enc, err := NewEncryptor("N1PCdw3M2B1TfJhoaY2mL736p2vCUc47")
	s.Require().NoError(err)

	s.encryptor = enc

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func (s *encryptorTestSuite) TearDownSuite() {
}

func (s *encryptorTestSuite) TestNewFromEnv() {
	tests := []struct {
		name string
		err  bool
	}{
		{
			name: "Env key is set",
			err:  false,
		},
		{
			name: "Env key is not set",
			err:  true,
		},
		{
			name: "Env key is too short",
			err:  true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.name == "Env key is set" {
				s.Require().NoError(os.Setenv("STORAGE_ENCRYPTION_KEY", "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"))
			} else if tt.name == "Env key is not set" {
				s.Require().NoError(os.Setenv("STORAGE_ENCRYPTION_KEY", ""))
			} else if tt.name == "Env key is too short" {
				s.Require().NoError(os.Setenv("STORAGE_ENCRYPTION_KEY", "imtooshort"))
			}

			_, err := NewEncryptorFromEnv()
			if tt.err {
				s.Error(err)
				return
			}

			s.NoError(err)
		})
	}
}

func (s *encryptorTestSuite) TestEncrypt() {
	sampleText := "test123"

	encryptedText, err := s.encryptor.Encrypt(sampleText)
	s.NoError(err)
	s.NotNil(encryptedText)

	dec, err := s.encryptor.Decrypt(*encryptedText)
	s.NoError(err)
	s.NotNil(dec)
	s.Equal(sampleText, *dec)
}

func (s *encryptorTestSuite) TestNoop() {
	s.encryptor = NewNoopEncryptor()

	sampleText := "test123"

	encryptedText, err := s.encryptor.Encrypt(sampleText)
	s.NoError(err)
	s.EqualValues(sampleText, *encryptedText)

	dec, err := s.encryptor.Decrypt(*encryptedText)
	s.NoError(err)
	s.NotNil(dec)
	s.EqualValues(sampleText, *dec)
}

func Test(t *testing.T) {
	suite.Run(t, new(encryptorTestSuite))
}
