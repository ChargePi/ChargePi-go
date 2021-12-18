package settings_manager

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	cache2 "github.com/xBlaz3kx/ChargePi-go/components/cache"
	"github.com/xBlaz3kx/ChargePi-go/data/settings"
	"os/exec"
	"testing"
	"time"
)

type SettingsManagerTestSuite struct {
	suite.Suite
	connector  settings.Connector
	session    settings.Session
	relay      settings.Relay
	powerMeter settings.PowerMeter
}

func (s *SettingsManagerTestSuite) SetupTest() {
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)
	s.session = settings.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	s.relay = settings.Relay{
		RelayPin:     1,
		InverseLogic: false,
	}

	s.powerMeter = settings.PowerMeter{
		Enabled:              false,
		Type:                 "",
		PowerMeterPin:        0,
		SpiBus:               0,
		PowerUnits:           "",
		Consumption:          0,
		ShuntOffset:          0,
		VoltageDividerOffset: 0,
	}

	s.connector = settings.Connector{
		EvseId:      1,
		ConnectorId: 1,
		Type:        "Schuko",
		Status:      "Available",
		Session:     s.session,
		Relay:       s.relay,
		PowerMeter:  s.powerMeter,
	}

	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", s.connector.EvseId, s.connector.ConnectorId)
		cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", s.connector.EvseId, s.connector.ConnectorId)
	)

	cache2.Cache.Set(cachePathKey, "./connector-1.json", cache.DefaultExpiration)
	cache2.Cache.Set(cacheConnectorKey, &s.connector, cache.DefaultExpiration)
}

func (s *SettingsManagerTestSuite) TestUpdateSessionInfo() {
	var (
		connectorFromFile settings.Connector
		newSession        = settings.Session{
			IsActive:      true,
			TransactionId: "Transaction1234",
			TagId:         "Tag1234",
			Started:       "",
			Consumption:   nil,
		}
	)

	UpdateConnectorSessionInfo(s.connector.EvseId, s.connector.ConnectorId, &newSession)

	err := fig.Load(&connectorFromFile, fig.File("connector-1.json"))
	s.Require().FileExists("./connector-1.json")
	s.Require().NoError(err)
	s.Require().EqualValues(newSession, connectorFromFile.Session)

	// Clean up
	exec.Command("rm connector-1.json")
}

func (s *SettingsManagerTestSuite) TestUpdateConnectorStatus() {
	var connectorFromFile settings.Connector

	UpdateConnectorStatus(s.connector.EvseId, s.connector.ConnectorId, core.ChargePointStatusCharging)

	err := fig.Load(&connectorFromFile, fig.File("connector-1.json"), fig.Dirs("."))
	s.Require().FileExists("./connector-1.json")
	s.Require().NoError(err)

	s.Require().EqualValues(core.ChargePointStatusCharging, connectorFromFile.Status)

	// Delete the unnecessary file
	exec.Command("rm connector-1.json")
}

func TestSettingsManager(t *testing.T) {
	suite.Run(t, new(SettingsManagerTestSuite))
}
