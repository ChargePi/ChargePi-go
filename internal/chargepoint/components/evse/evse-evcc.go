package evse

import "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/evcc"

func (evse *Impl) Lock() {
	evse.evcc.Lock()
}

func (evse *Impl) Unlock() {
	evse.evcc.Unlock()
}

func (evse *Impl) GetConnectors() []Connector {
	return evse.connectors
}

func (evse *Impl) AddConnector(connector Connector) error {
	for _, c := range evse.connectors {
		// Do not add if they're the same connector
		if c.ConnectorId == connector.ConnectorId {
			return nil
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
