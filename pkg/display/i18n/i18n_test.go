package i18n

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type I18NTestSuite struct {
	suite.Suite
}

func (suite *I18NTestSuite) SetupTest() {
}

func (suite *I18NTestSuite) Test() {
}

func TestI18N(t *testing.T) {
	suite.Run(t, new(I18NTestSuite))
}
