package settings

const (
	currentFolder = "./configs"
	evseFolder    = "./configs/evses"
	dockerFolder  = "/etc/ChargePi/configs"
)

// Configuration variables
const (
	Model           = "chargepoint.info.ocpp.model"
	Vendor          = "chargepoint.info.ocpp.vendor"
	MaxChargingTime = "chargepoint.info.maxChargingTime"
	ProtocolVersion = "chargepoint.info.protocolVersion"
	LoggingFormat   = "chargepoint.logging.format"

	Debug      = "debug"
	ApiEnabled = "api.enabled"
	ApiAddress = "api.address"
	ApiPort    = "api.port"
)

// Flags
const (
	DebugFlag          = "debug"
	ApiFlag            = "api"
	ApiPortFlag        = "api-port"
	ApiAddressFlag     = "api-address"
	SettingsFlag       = "settings"
	EvseFlag           = "connector-folder"
	AuthFileFlag       = "auth"
	OcppConfigPathFlag = "ocpp-config"
)
