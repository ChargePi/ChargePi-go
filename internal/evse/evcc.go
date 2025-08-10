package evse

import (
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/pkg/evcc"
	"go.uber.org/zap"
)

func (evse *Impl) Lock() {
	evse.logger.Debug("Locking EVCC")
	evse.evcc.Lock()
}

func (evse *Impl) Unlock() {
	evse.logger.Debug("Unlocking EVCC")
	evse.evcc.Unlock()
}

func (evse *Impl) GetConnectors() []settings.Connector {
	evse.logger.Debug("Getting connectors for EVSE")
	return evse.connectors
}

func (evse *Impl) AddConnector(connector settings.Connector) error {
	evse.logger.With(zap.Any("connector", connector)).Debug("Adding connector to EVSE")
	for _, c := range evse.connectors {
		// Do not add if they're the same connector
		if c.ConnectorId == connector.ConnectorId {
			return ErrConnectorExists
		}
	}

	evse.connectors = append(evse.connectors, connector)
	return nil
}

func (evse *Impl) GetEvcc() evcc.EVCC {
	evse.logger.Debug("Getting EVCC")
	return evse.evcc
}

func (evse *Impl) SetEvcc(e evcc.EVCC) {
	evse.logger.Debug("Setting EVCC")

	// Cleanup the previous EVCC
	err := evse.evcc.Cleanup()
	if err != nil {
		evse.logger.With(zap.Error(err)).Error("Error cleaning up EVCC")
	}

	evse.evcc = e
}
