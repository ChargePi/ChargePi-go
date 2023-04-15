package settings

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var exporter Exporter

func GetExporter() Exporter {
	if exporter == nil {
		log.Debug("Creating an exporter")
		db := database.Get()
		exporter = &ExporterImpl{
			db:              db,
			tagManager:      auth.NewTagManager(db),
			settingsManager: GetManager(),
		}
	}

	return exporter
}

type Exporter interface {
	ExportEVSESettings() []settings.EVSE
	ExportOcppConfiguration() configuration.Config
	ExportLocalAuthList() (*settings.AuthList, error)
	ExportChargePointSettings() (*settings.Settings, error)
}

type ExporterImpl struct {
	db              *badger.DB
	tagManager      auth.TagManager
	settingsManager Manager
}

func (i *ExporterImpl) ExportEVSESettings() []settings.EVSE {
	log.Debug("Exporting EVSE settings from the database")
	return database.GetEvseSettings(i.db)
}

func (i *ExporterImpl) ExportOcppConfiguration() configuration.Config {
	log.Debug("Exporting OCPP configuration from the database")
	getConfiguration, _ := i.settingsManager.GetOcppConfiguration(configuration.OCPP16)
	return configuration.Config{
		Version: 1,
		Keys:    getConfiguration,
	}
}

func (i *ExporterImpl) ExportLocalAuthList() (*settings.AuthList, error) {
	log.Debug("Exporting Local auth list from the database")
	return &settings.AuthList{
		Version: i.tagManager.GetAuthListVersion(),
		Tags:    i.tagManager.GetTags(),
	}, nil
}

func (i *ExporterImpl) ExportChargePointSettings() (*settings.Settings, error) {
	log.Debug("Exporting charge point settings from the database")
	return i.settingsManager.GetSettings()
}
