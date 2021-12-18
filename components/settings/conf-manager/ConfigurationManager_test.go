package conf_manager

import (
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	cache2 "github.com/xBlaz3kx/ChargePi-go/components/cache"
	"github.com/xBlaz3kx/ChargePi-go/data/ocpp"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type ConfigurationManagerTestSuite struct {
	suite.Suite
	config ocpp.Config
}

func (s *ConfigurationManagerTestSuite) SetupTest() {
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	s.config = ocpp.Config{
		Version: 1,
		Keys: []core.ConfigurationKey{
			{
				Key:      "Test1",
				Readonly: false,
				Value:    "60",
			}, {
				Key:      "Test2",
				Readonly: false,
				Value:    "ABCD",
			},
			{
				Key:      "MeterValuesSampledData",
				Readonly: false,
				Value: strings.Join(
					[]string{
						string(types.MeasurandEnergyActiveExportInterval),
						string(types.MeasurandCurrentExport),
						string(types.MeasurandVoltage),
					},
					",",
				),
			},
		}}

	cache2.Cache.Set("OCPPConfiguration", &s.config, cache.DefaultExpiration)
}

func (s *ConfigurationManagerTestSuite) TestGetConfiguration() {
	configuration, err := GetConfiguration()
	s.Require().NoError(err)
	s.Require().Equal(s.config, *configuration)

	cache2.Cache.Delete("OCPPConfiguration")
	configuration, err = GetConfiguration()

	s.Require().Error(err)
	s.Require().Nil(configuration)
}

func (s *ConfigurationManagerTestSuite) TestGetConfigurationValue() {
	// Ok case
	value, err := GetConfigurationValue("Test1")
	s.Require().NoError(err)
	s.Require().Equal("60", value)

	// No such key
	value, err = GetConfigurationValue("Test123")
	s.Require().Error(err)
	s.Require().Equal("", value)
}

func (s *ConfigurationManagerTestSuite) TestGetTypesToSample() {
	measurands := GetTypesToSample()
	s.Require().NotNil(measurands)
}

func (s *ConfigurationManagerTestSuite) TestUpdateConfigurationFile() {
	var (
		fileConfig ocpp.Config
	)
	cache2.Cache.Set("configurationFilePath", "./configuration.json", cache.DefaultExpiration)

	err := UpdateConfigurationFile()
	s.Require().NoError(err)
	err = fig.Load(&fileConfig, fig.File("configuration.json"), fig.Dirs("."))
	s.Require().NoError(err)
	s.Require().Equal(s.config, fileConfig)

	// Delete the config from cache
	cache2.Cache.Delete("OCPPConfiguration")
	err = UpdateConfigurationFile()
	s.Require().Error(err)

	// Delete the path from cache
	cache2.Cache.Delete("configurationFilePath")
	err = UpdateConfigurationFile()
	s.Require().Error(err)

	// Delete the unnecessary file
	exec.Command("rm configuration.json")
}

func TestConfigurationManager(t *testing.T) {
	suite.Run(t, new(ConfigurationManagerTestSuite))
}
