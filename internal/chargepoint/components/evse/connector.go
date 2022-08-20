package evse

type (
	Connector struct {
		ConnectorId int    `fig:"ConnectorId" validate:"required" json:"ConnectorId,omitempty" yaml:"ConnectorId" mapstructure:"ConnectorId"`
		Type        string `fig:"Type" validate:"required" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		Status      string `fig:"Status" validation:"required" json:"status,omitempty" yaml:"status" mapstructure:"status"`
	}
)
