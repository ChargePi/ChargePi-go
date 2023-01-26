package ocpp

type (
	DataTransferEVSEInfo struct {
		EvseId     int         `json:"evseId" yaml:"evseId" mapstructure:"evseId"`
		MaxPower   float32     `json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		Connectors []Connector `json:"connectors,omitempty" yaml:"connectors" mapstructure:"connectors"`
	}

	Connector struct {
		ConnectorId int    `json:"connectorId,omitempty" yaml:"connectorId" mapstructure:"connectorId"`
		Type        string `json:"type,omitempty" yaml:"type" mapstructure:"type"`
	}
)

func NewEvseInfo(evseId int, maxPower float32, connectors ...Connector) DataTransferEVSEInfo {
	return DataTransferEVSEInfo{
		EvseId:     evseId,
		MaxPower:   maxPower,
		Connectors: connectors,
	}
}

func NewConnector(connectorId int, connectorType string) Connector {
	return Connector{
		ConnectorId: connectorId,
		Type:        connectorType,
	}
}
