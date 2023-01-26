package settings

import (
	"encoding/json"
	"sync"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var (
	manager Manager
	once    = sync.Once{}
)

func init() {
	once.Do(func() {
		GetManager()
		GetExporter()
		GetImporter()
	})
}

func GetManager() Manager {
	if manager == nil {
		log.Debug("Creating EVSE manager")
		manager = NewManager(database.Get())
	}

	return manager
}

type (
	Manager interface {
		SetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string)
		GetChargePointSettings() settings.ChargePoint
		SetChargePointSettings(settings settings.ChargePoint) error
	}

	Impl struct {
		db                  *badger.DB
		ocppVariableManager ocppConfigManager.Manager
	}
)

func (i *Impl) GetChargePointSettings() settings.ChargePoint {
	// todo
	return settings.ChargePoint{}
}

func (i *Impl) SetChargePointSettings(settings settings.ChargePoint) error {
	// todo
	return nil
}

func NewManager(db *badger.DB) *Impl {
	return &Impl{
		db:                  db,
		ocppVariableManager: ocppConfigManager.GetManager(),
	}
}

func (i *Impl) SetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string) {
	logInfo := log.WithField("version", version)
	logInfo.Info("Setting up OCPP configuration")

	i.ocppVariableManager.SetVersion(version)
	i.ocppVariableManager.SetSupportedProfiles(supportedProfiles...)

	ocppConfig := configuration.Config{}
	err := i.db.View(func(txn *badger.Txn) error {

		confg, err := txn.Get(database.GetOcppConfigurationKey())
		if err != nil {
			return err
		}

		bytes, err := confg.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bytes, &ocppConfig)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logInfo.WithError(err).Errorf("Error reading configuration from database")
		return
	}

	err = i.ocppVariableManager.SetConfiguration(ocppConfig)
	if err != nil {
		logInfo.WithError(err).Errorf("Error setting the configuration to the manager")
		return
	}
}
