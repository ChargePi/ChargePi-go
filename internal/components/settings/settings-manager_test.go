package settings

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	cache2 "github.com/xBlaz3kx/ChargePi-go/pkg/cache"
	"os/exec"
	"testing"
	"time"
)

type SettingsManagerTestSuite struct {
	suite.Suite
	connector  settings2.Connector
	session    settings2.Session
	relay      settings2.Relay
	powerMeter settings2.PowerMeter
}

func (s *SettingsManagerTestSuite) SetupTest() {
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)
	s.session = settings2.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	s.relay = settings2.Relay{
		RelayPin:     1,
		InverseLogic: false,
	}

	s.powerMeter = settings2.PowerMeter{
		Enabled:              false,
		Type:                 "",
		PowerMeterPin:        0,
		SpiBus:               0,
		Consumption:          0,
		ShuntOffset:          0,
		VoltageDividerOffset: 0,
	}

	s.connector = settings2.Connector{
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
		connectorFromFile settings2.Connector
		newSession        = settings2.Session{
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
	var connectorFromFile settings2.Connector

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
