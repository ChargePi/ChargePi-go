package evse

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
)

func (evse *Impl) Lock() {
	evse.evcc.Lock()
}

func (evse *Impl) Unlock() {
	evse.evcc.Unlock()
}

func (evse *Impl) GetConnectors() []settings.Connector {
	return evse.connectors
}

func (evse *Impl) AddConnector(connector settings.Connector) error {
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
	return evse.evcc
}

func (evse *Impl) SetEvcc(e evcc.EVCC) {
	evse.evcc = e
}
