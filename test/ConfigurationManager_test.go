package test

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cache2 "github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"os/exec"
	"testing"
	"time"
)

func TestGetConfiguration(t *testing.T) {
	require := require.New(t)
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	var config = settings.OCPPConfig{
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
		}}

	cache2.Cache.Set("OCPPConfiguration", &config, cache.DefaultExpiration)
	configuration, err := settings.GetConfiguration()
	require.NoError(err)
	require.Equal(config, *configuration)

	cache2.Cache.Delete("OCPPConfiguration")
	configuration, err = settings.GetConfiguration()

	require.Error(err)
	require.Nil(configuration)
}

func TestGetConfigurationValue(t *testing.T) {
	require := require.New(t)
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	var config = settings.OCPPConfig{
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
			}, {
				Key:      "Test3",
				Readonly: false,
				Value:    "",
			},
		}}

	cache2.Cache.Set("OCPPConfiguration", &config, cache.DefaultExpiration)
	// ok case
	value, err := settings.GetConfigurationValue("Test1")
	require.NoError(err)
	require.Equal("60", value)

	//no such key
	value, err = settings.GetConfigurationValue("Test123")
	require.Error(err)
	require.Equal("", value)
}

func TestOCPPConfig_GetConfig(t *testing.T) {
	assert := assert.New(t)

	var keys = []core.ConfigurationKey{
		{
			Key:      "Test1",
			Readonly: false,
			Value:    "60",
		}, {
			Key:      "Test2",
			Readonly: false,
			Value:    "ABCD",
		}, {
			Key:      "Test3",
			Readonly: false,
			Value:    "",
		},
	}

	var config = settings.OCPPConfig{
		Version: 1,
		Keys:    keys,
	}

	assert.Equal(keys, config.GetConfig())
	// reset the config
	config = settings.OCPPConfig{
		Version: 1,
		Keys:    []core.ConfigurationKey{},
	}

	assert.Equal([]core.ConfigurationKey{}, config.GetConfig())
}

func TestOCPPConfig_UpdateKey(t *testing.T) {
	assert := assert.New(t)

	var keys = []core.ConfigurationKey{
		{
			Key:      "Test1",
			Readonly: false,
			Value:    "60",
		}, {
			Key:      "Test2",
			Readonly: true,
			Value:    "ABCD",
		}, {
			Key:      "Test3",
			Readonly: false,
			Value:    "",
		},
	}

	var config = settings.OCPPConfig{
		Version: 1,
		Keys:    keys,
	}

	// ok case
	err := config.UpdateKey("Test1", "1234")
	assert.NoError(err)
	value, err := config.GetConfigurationValue("Test1")
	assert.NoError(err)
	assert.Equal("1234", value)

	// invalid key
	err = config.UpdateKey("Test4", "1234")
	assert.Error(err)

	// is readOnly
	err = config.UpdateKey("Test2", "BCDEF")
	assert.Error(err)
	value, err = config.GetConfigurationValue("Test2")
	assert.NoError(err)
	assert.Equal("ABCD", value)
}

func TestUpdateConfigurationFile(t *testing.T) {
	require := assert.New(t)
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	var fileConfig settings.OCPPConfig
	var config = settings.OCPPConfig{
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
			}, {
				Key:      "Test3",
				Readonly: false,
				Value:    "",
			},
		}}

	cache2.Cache.Set("OCPPConfiguration", &config, cache.DefaultExpiration)
	cache2.Cache.Set("configurationFilePath", "./configuration.json", cache.DefaultExpiration)

	err := settings.UpdateConfigurationFile()
	require.NoError(err)
	err = fig.Load(&fileConfig, fig.File("configuration.json"), fig.Dirs("."))
	require.NoError(err)
	require.Equal(config, fileConfig)

	cache2.Cache.Delete("OCPPConfiguration")
	err = settings.UpdateConfigurationFile()
	require.Error(err)

	cache2.Cache.Delete("configurationFilePath")
	err = settings.UpdateConfigurationFile()
	require.Error(err)

	// delete the unnecessary file
	exec.Command("rm configuration.json")
}

func TestUpdateConnectorSessionInfo(t *testing.T) {
	require := assert.New(t)
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	connector := settings.Connector{
		EvseId:      1,
		ConnectorId: 1,
		Type:        "Schuko",
		Status:      "Available",
		Session: struct {
			IsActive      bool   `fig:"IsActive"`
			TransactionId string `fig:"TransactionId" default:""`
			TagId         string `fig:"TagId" default:""`
			Started       string `fig:"Started" default:""`
			Consumption   []types.MeterValue
		}{
			IsActive:      false,
			TransactionId: "",
			TagId:         "",
			Started:       "",
			Consumption:   nil,
		},
		Relay: struct {
			RelayPin     int  `fig:"RelayPin" validate:"required"`
			InverseLogic bool `fig:"InverseLogic"`
		}{
			RelayPin:     1,
			InverseLogic: false,
		},
		PowerMeter: struct {
			Enabled              bool    `fig:"Enabled"`
			Type                 string  `fig:"Type"`
			PowerMeterPin        int     `fig:"PowerMeterPin"`
			SpiBus               int     `fig:"SpiBus" default:"0"`
			PowerUnits           string  `fig:"PowerUnits"`
			Consumption          float64 `fig:"Consumption"`
			ShuntOffset          float64 `fig:"ShuntOffset"`
			VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
		}{
			Enabled:              false,
			Type:                 "",
			PowerMeterPin:        0,
			SpiBus:               0,
			PowerUnits:           "",
			Consumption:          0,
			ShuntOffset:          0,
			VoltageDividerOffset: 0,
		},
	}

	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", connector.EvseId, connector.ConnectorId)
		cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId)
		connectorFromFile settings.Connector

		newSession = settings.Session{
			IsActive:      true,
			TransactionId: "Transaction1234",
			TagId:         "Tag1234",
			Started:       "",
			Consumption:   nil,
		}
	)

	cache2.Cache.Set(cachePathKey, "./connector-1.json", cache.DefaultExpiration)
	cache2.Cache.Set(cacheConnectorKey, &connector, cache.DefaultExpiration)

	settings.UpdateConnectorSessionInfo(connector.EvseId, connector.ConnectorId, &newSession)

	err := fig.Load(&connectorFromFile, fig.File("connector-1.json"), fig.Dirs("."))
	require.FileExists("./connector-1.json")
	require.NoError(err)
	require.EqualValues(newSession, connectorFromFile.Session)

	// delete the unnecessary file
	exec.Command("rm connector-1.json")
}

func TestUpdateConnectorStatus(t *testing.T) {
	require := assert.New(t)
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)

	connector := settings.Connector{
		EvseId:      1,
		ConnectorId: 1,
		Type:        "Schuko",
		Status:      "Available",
		Session: struct {
			IsActive      bool   `fig:"IsActive"`
			TransactionId string `fig:"TransactionId" default:""`
			TagId         string `fig:"TagId" default:""`
			Started       string `fig:"Started" default:""`
			Consumption   []types.MeterValue
		}{
			IsActive:      false,
			TransactionId: "",
			TagId:         "",
			Started:       "",
			Consumption:   nil,
		},
		Relay: struct {
			RelayPin     int  `fig:"RelayPin" validate:"required"`
			InverseLogic bool `fig:"InverseLogic"`
		}{
			RelayPin:     1,
			InverseLogic: false,
		},
		PowerMeter: struct {
			Enabled              bool    `fig:"Enabled"`
			Type                 string  `fig:"Type"`
			PowerMeterPin        int     `fig:"PowerMeterPin"`
			SpiBus               int     `fig:"SpiBus" default:"0"`
			PowerUnits           string  `fig:"PowerUnits"`
			Consumption          float64 `fig:"Consumption"`
			ShuntOffset          float64 `fig:"ShuntOffset"`
			VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
		}{
			Enabled:              false,
			Type:                 "",
			PowerMeterPin:        0,
			SpiBus:               0,
			PowerUnits:           "",
			Consumption:          0,
			ShuntOffset:          0,
			VoltageDividerOffset: 0,
		},
	}

	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", connector.EvseId, connector.ConnectorId)
		cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId)
		connectorFromFile settings.Connector
	)

	cache2.Cache.Set(cachePathKey, "./connector-1.json", cache.DefaultExpiration)
	cache2.Cache.Set(cacheConnectorKey, &connector, cache.DefaultExpiration)

	settings.UpdateConnectorStatus(connector.EvseId, connector.ConnectorId, core.ChargePointStatusCharging)

	err := fig.Load(&connectorFromFile, fig.File("connector-1.json"), fig.Dirs("."))
	require.FileExists("./connector-1.json")
	require.NoError(err)

	require.EqualValues(core.ChargePointStatusCharging, connectorFromFile.Status)

	// delete the unnecessary file
	exec.Command("rm connector-1.json")
}

func TestWriteToFile(t *testing.T) {
	require := assert.New(t)

	Test123 := struct {
		Enabled bool   `json:"enabled"`
		Type    string `json:"type"`
	}{
		Enabled: false,
		Type:    "",
	}

	err := settings.WriteToFile("test123.json", &Test123)
	require.NoError(err)
	require.FileExists("test123.json")

	err = settings.WriteToFile("test123.yaml", &Test123)
	require.NoError(err)
	require.FileExists("test123.yaml")

	err = settings.WriteToFile("test123.o", &Test123)
	require.Error(err)

}
