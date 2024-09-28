package diagnostics

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type httpUploaderTestSuite struct {
	suite.Suite
}

func (s *httpUploaderTestSuite) SetupTest() {
}

func (s *httpUploaderTestSuite) TearDownSuite() {
}

func (s *httpUploaderTestSuite) TestUpload() {

}

func TestHttpUploader(t *testing.T) {
	suite.Run(t, new(httpUploaderTestSuite))
}
