package evse

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
)

func (evse *Impl) Lock() {
	evse.logger.Debugf("Locking EVCC")
	evse.evcc.Lock()
}

func (evse *Impl) Unlock() {
	evse.logger.Debugf("Unlocking EVCC")
	evse.evcc.Unlock()
}

func (evse *Impl) GetConnectors() []settings.Connector {
	evse.logger.Debugf("Getting connectors for EVSE")
	return evse.connectors
}

func (evse *Impl) AddConnector(connector settings.Connector) error {
	evse.logger.WithField("connectorId", connector.ConnectorId).Debug("Adding connector to EVSE")
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
	evse.logger.Debugf("Getting EVCC")
	return evse.evcc
}

func (evse *Impl) SetEvcc(e evcc.EVCC) {
	evse.logger.Debugf("Setting EVCC")

	// Cleanup the previous EVCC
	err := evse.evcc.Cleanup()
	if err != nil {
		evse.logger.Errorf("Error cleaning up EVCC: %s", err)
	}

	evse.evcc = e
}
