package settings

const DatabasePath = "/tmp/chargepi"

const (
	CurrentFolder = "./configs"
	EvseFolder    = "./configs/evses"
	DockerFolder  = "/etc/ChargePi/configs"
)

// Configuration variables
const (
	Model           = "chargepoint.info.ocpp.model"
	Vendor          = "chargepoint.info.ocpp.vendor"
	MaxChargingTime = "chargepoint.info.maxChargingTime"
	ProtocolVersion = "chargepoint.info.protocolVersion"

	Debug      = "debug"
	ApiEnabled = "api.enabled"
	ApiAddress = "api.address"
	ApiPort    = "api.port"
)

// Flags
const (
	DebugFlag          = "debug"
	ApiAddressFlag     = "api-address"
	SettingsFlag       = "settings"
	EvseFlag           = "evse"
	AuthFileFlag       = "auth"
	OcppConfigPathFlag = "ocppCfg"
	OcppVersion        = "v"
)
