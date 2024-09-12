package ocpp

type DataTransferEVSEInfo struct {
	// The ID of the EVSE.
	EvseId int `json:"evseId" yaml:"evseId" mapstructure:"evseId"`

	// The maximum power that the EVSE can deliver.
	MaxPower float32 `json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`

	// The connectors that are available on the EVSE. Must contain at least one connector.
	Connectors []Connector `json:"connectors,omitempty" yaml:"connectors" mapstructure:"connectors"`
}

func NewEvseInfo(evseId int, maxPower float32, connectors ...Connector) DataTransferEVSEInfo {
	return DataTransferEVSEInfo{
		EvseId:     evseId,
		MaxPower:   maxPower,
		Connectors: connectors,
	}
}

type Connector struct {
	// The ID of the connector.
	ConnectorId int `json:"connectorId,omitempty" yaml:"connectorId" mapstructure:"connectorId"`

	// The type of the connector. This can be one of the following values:
	// CCS1, CCS2, CHADEMO, TESLA, TYPE_1, TYPE_2,
	Type string `json:"type,omitempty" yaml:"type" mapstructure:"type"`
}

func NewConnector(connectorId int, connectorType string) Connector {
	return Connector{
		ConnectorId: connectorId,
		Type:        connectorType,
	}
}
