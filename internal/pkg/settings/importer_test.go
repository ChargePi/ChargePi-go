package settings

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type importerTestSuite struct {
	suite.Suite
}

func (s *importerTestSuite) SetupTest() {
}

func (s *importerTestSuite) TestImportEVSESettings() {
}

func (s *importerTestSuite) TestImportOcppConfiguration() {
}

func (s *importerTestSuite) TestImportLocalAuthList() {
}

func TestImporter(t *testing.T) {
	suite.Run(t, new(importerTestSuite))
}
