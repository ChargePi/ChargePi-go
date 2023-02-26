package settings

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type exporterTestSuite struct {
	suite.Suite
}

func (s *exporterTestSuite) SetupTest() {
}

func (s *exporterTestSuite) TestExportEVSESettings() {
}

func (s *exporterTestSuite) TestExportOcppConfiguration() {
}

func (s *exporterTestSuite) TestExportLocalAuthList() {
}

func TestExporter(t *testing.T) {
	suite.Run(t, new(exporterTestSuite))
}
