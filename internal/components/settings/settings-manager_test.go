package settings

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	settingsData "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	cache2 "github.com/xBlaz3kx/ChargePi-go/pkg/cache"
	"os/exec"
	"testing"
	"time"
)

const (
	fileName = "connector-1.json"
)

type SettingsManagerTestSuite struct {
	suite.Suite
	connector  settingsData.Connector
	session    settingsData.Session
	relay      settingsData.Relay
	powerMeter settingsData.PowerMeter
}

func (s *SettingsManagerTestSuite) SetupTest() {
	cache2.Cache = cache.New(time.Minute*10, time.Minute*10)
	s.session = settingsData.Session{
		IsActive:      false,
		TransactionId: "",
		TagId:         "",
		Started:       "",
		Consumption:   nil,
	}

	s.relay = settingsData.Relay{
		RelayPin:     1,
		InverseLogic: false,
	}

	s.powerMeter = settingsData.PowerMeter{
		Enabled:              false,
		Type:                 "",
		PowerMeterPin:        0,
		SpiBus:               0,
		Consumption:          0,
		ShuntOffset:          0,
		VoltageDividerOffset: 0,
	}

	s.connector = settingsData.Connector{
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

	cache2.Cache.Set(cachePathKey, "./"+fileName, cache.DefaultExpiration)
	cache2.Cache.Set(cacheConnectorKey, &s.connector, cache.DefaultExpiration)
}

func (s *SettingsManagerTestSuite) TestUpdateSessionInfo() {
	var (
		connectorFromFile settingsData.Connector
		newSession        = settingsData.Session{
			IsActive:      true,
			TransactionId: "Transaction1234",
			TagId:         "Tag1234",
			Started:       "",
			Consumption:   nil,
		}
	)

	UpdateConnectorSessionInfo(s.connector.EvseId, s.connector.ConnectorId, &newSession)

	err := fig.Load(&connectorFromFile, fig.File(fileName))
	s.Require().FileExists("./" + fileName)
	s.Require().NoError(err)
	s.Require().EqualValues(newSession, connectorFromFile.Session)

	// Clean up
	cmd := exec.Command("rm", fileName)
	err = cmd.Run()
	s.Require().NoError(err)
}

func (s *SettingsManagerTestSuite) TestUpdateConnectorStatus() {
	cmd := exec.Command("touch", fileName)
	err := cmd.Run()
	s.Require().NoError(err)

	var connectorFromFile settingsData.Connector

	UpdateConnectorStatus(s.connector.EvseId, s.connector.ConnectorId, core.ChargePointStatusCharging)

	time.Sleep(time.Second)

	err = fig.Load(&connectorFromFile, fig.File(fileName))
	s.Require().FileExists("./" + fileName)
	s.Require().NoError(err)

	s.Require().EqualValues(core.ChargePointStatusCharging, connectorFromFile.Status)

	// Clean up
	cmd = exec.Command("rm", fileName)
	err = cmd.Run()
	s.Require().NoError(err)
}

func TestSettingsManager(t *testing.T) {
	suite.Run(t, new(SettingsManagerTestSuite))
}
