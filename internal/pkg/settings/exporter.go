package settings

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var exporter Exporter

func GetExporter() Exporter {
	if exporter == nil {
		log.Debug("Creating EVSE manager")
		exporter = &ExporterImpl{db: database.Get()}
	}

	return exporter
}

type Exporter interface {
	ExportEVSESettings() []settings.EVSE
	ExportOcppConfiguration() configuration.Config
	ExportLocalAuthList() (*settings.AuthList, error)
}

type ExporterImpl struct {
	db                  *badger.DB
	ocppVariableManager ocppConfigManager.Manager
	tagManager          auth.TagManager
}

func (i *ExporterImpl) ExportEVSESettings() []settings.EVSE {
	log.Info("Exporting EVSE settings")
	return database.GetEvseSettings(i.db)
}

func (i *ExporterImpl) ExportOcppConfiguration() configuration.Config {
	log.Info("Exporting OCPP configuration")
	getConfiguration, _ := i.ocppVariableManager.GetConfiguration()
	return configuration.Config{
		Version: 1,
		Keys:    getConfiguration,
	}
}

func (i *ExporterImpl) ExportLocalAuthList() (*settings.AuthList, error) {
	log.Info("Exporting Local auth list")
	return &settings.AuthList{
		Version: i.tagManager.GetAuthListVersion(),
		Tags:    i.tagManager.GetTags(),
	}, nil
}
