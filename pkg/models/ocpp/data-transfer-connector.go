package ocpp

type (
	DataTransferEVSEInfo struct {
		EvseId     int         `fig:"evseId" json:"evseId" yaml:"evseId" mapstructure:"evseId"`
		MaxPower   float32     `fig:"maxPower" json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		Connectors []Connector `fig:"connectors" json:"connectors,omitempty" yaml:"connectors" mapstructure:"connectors"`
	}

	Connector struct {
		ConnectorId int    `fig:"ConnectorId" json:"ConnectorId,omitempty" yaml:"ConnectorId" mapstructure:"ConnectorId"`
		Type        string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
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
