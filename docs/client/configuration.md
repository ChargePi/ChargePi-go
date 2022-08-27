# 🛠️ Configuration

## 🚩 Flags and environment variables

Available flags:

|        Flag         | Short |           Description           | Default value |
|:-------------------:|:-----:|:-------------------------------:|:-------------:|
|     `-settings`     |   /   |   Path to the settings file.    |               |
| `-connector-folder` |   /   |  Path to the connector folder.  |               |
|   `-ocpp-config`    |   /   | Path to the OCPP configuration. |               |
|       `-auth`       |   /   | Path to the authorization file. |               |
|      `-debug`       | `--d` |           Debug mode            |     false     |
|       `-api`        | `--a` |         Expose the API          |     false     |
|   `-api-address`    |   /   |           API address           |  "localhost"  |
|     `-api-port`     |   /   |            API port             |     4269      |

Environment variables are created automatically thanks to [Viper](https://github.com/spf13/viper) and are prefixed
with `CHARGEPI`. Only the main settings are bound to the env variables. Debug mode and API settings flags are also bound
to the environment variables.

Example environment variable: `CHARGEPI_CHARGEPOINT_CONNECTIONSETTINGS_ID`.

## 🛠 Configuration files

There are three **required** configuration files:

1. [`settings`](../../configs/settings.json)
2. [`configuration`](../../configs/configuration.json)
3. [`evses`](../../configs/evses/connector-1.json)

The settings files are supported in `YAML`, `JSON` or `TOML` format.

### The `settings` file

The `settings` file contains basic information about the charge point and provides connectivity details:

- Charge Point ID,
- central system URI and OCPP protocol version,
- OCPP information (vendor, model, firmware, etc.),
- logging settings,
- TLS settings,
- default max charging time,
- settings for LCD, RFID/NFC reader and LEDs.

The table represents attributes, their values and descriptions that require more attention and might not be
self-explanatory. Some attributes can have multiple possible values, if any are empty, they will be treated as disabled
or might not work properly.

#### chargePoint connectionSettings

|    Attribute    |                            Description                             |  Possible values   | 
|:---------------:|:------------------------------------------------------------------:|:------------------:|
|       id        | ID of the charging point. Must be registered in the Central System | Default:"ChargePi" |
| protocolVersion |                   Version of the OCPP protocol.                    |   "1.6", "2.0.1"   |
|    serverUri    |       URI of the Central System with the port and endpoint.        |         ""         |
|  basicAuthUser  |      HTTP username for authentication with the Central System      |  Any string value  |
|  basicAuthPass  |      HTTP username for authentication with the Central System      |  Any string value  |
|       tls       |                       TLS certificate paths                        |                    |

#### chargePoint connectionSettings tls

|       Attribute       |        Description        | Possible values | 
|:---------------------:|:-------------------------:|:---------------:|
|   CACertificatePath   | Root/CA certificate path. | Any valid path  |
| clientCertificatePath |    Client certificate     | Any valid path  |
|     clientKeyPath     |    Client private key     | Any valid path  |

#### chargePoint info

|    Attribute    |                            Description                             |   Possible values   | 
|:---------------:|:------------------------------------------------------------------:|:-------------------:|
| maxChargingTime | Maxiumum amount of time that a transaction can last (in minutes).  |    Default: 180     |
|   ocpp.vendor   |                             Vendor ID                              | Default: "xBlaz3k"  |
|   ocpp.model    |                               Model                                | Default: "ChargePi" |

#### chargePoint hardware

##### chargePoint hardware display

| Attribute  |                      Description                       | Possible values | 
|:----------:|:------------------------------------------------------:|:---------------:|
| isEnabled  |               Enable or disable display                |   true, false   |
|   driver   |                  Display driver type                   |       ""        |
| i2cAddress | Field specific for any display using I2C communication |       ""        |
|   i2cBus   | Field specific for any display using I2C communication |    "1", "0"     | 
|  language  |   Language selection for messages displayed to LCD.    |    "en",sl"     |

##### chargePoint hardware tagReader

|   Attribute   |                         Description                         |                         Possible values                          | 
|:-------------:|:-----------------------------------------------------------:|:----------------------------------------------------------------:|
|   isEnabled   |                 Enable or disable tagReader                 |                        Default:"ChargePi"                        |
|  readerModel  |                Version of the OCPP protocol.                |                          "1.6", "2.0.1"                          |
|   resetPin    |    URI of the Central System with the port and endpoint.    | Default: "172.0.1.121:8080/steve/websocket/CentralSystemService" |
| deviceAddress | Max charging time allowed on the Charging point in minutes. |                           Default:180                            |

##### chargePoint hardware ledIndicator

|    Attribute     |                         Description                         |                         Possible values                          | 
|:----------------:|:-----------------------------------------------------------:|:----------------------------------------------------------------:|
|    isEnabled     |               Enable or disable ledIndicator                |                        Default:"ChargePi"                        |
|       type       |                Version of the OCPP protocol.                |                          "1.6", "2.0.1"                          |
|     dataPin      |    URI of the Central System with the port and endpoint.    | Default: "172.0.1.121:8080/steve/websocket/CentralSystemService" |
| indicateCardRead | Max charging time allowed on the Charging point in minutes. |                           Default:180                            |
|      invert      |                 RFID/NFC reader model used.                 |                           "PN532", ""                            | 

Example settings:

```json
{
  "chargePoint": {
    "info": {
      "id": "ChargePi",
      "protocolVersion": "1.6",
      "serverUri": "example.com",
      "basicAuthUser": "",
      "basicAuthPass": "",
      "maxChargingTime": 5,
      "ocpp": {
        "vendor": "UL FE",
        "model": "ChargePi"
      }
    },
    "logging": {
      "type": [
        "remote",
        "file"
      ],
      "format": "gelf",
      "host": "logging.example.com",
      "port": 12201
    },
    "tls": {
      "isEnabled": false,
      "CACertificatePath": "/usr/share/certs/rootCA.crt",
      "clientCertificatePath": "/usr/share/certs/charge-point.crt",
      "clientKeyPath": "/usr/share/certs/charge-point.key"
    },
    "hardware": {
      "lcd": {
        "isSupported": true,
        "driver": "hd44780",
        "i2cAddress": "0x27",
        "i2cBus": 1,
        "language": "en"
      },
      "tagReader": {
        "isSupported": true,
        "readerModel": "PN532",
        "device": "/dev/ttyS0",
        "resetPin": 19
      },
      "ledIndicator": {
        "enabled": true,
        "type": "WS281x",
        "dataPin": 18,
        "indicateCardRead": true,
        "invert": false
      }
    }
  }
}
```

### 🔌 The `evse` file(s) - EVSE configuration

EVSE settings files can be found in the `evse` folder. To add and configure the EVSE, simply add a new file with the
file structure as in the example. The client will scan the folder at boot, validate the configuration files and add the
EVSE with the provided configuration.

#### Attributes

| Attribute  |                Description                |             Possible values             | 
|:----------:|:-----------------------------------------:|:---------------------------------------:|
|   evseId   |              ID of the EVSE               |                   >1                    |
| connectors | List of available connectors for the EVSE ||
|    evcc    |            Charging Controller            ||
|   status   |            Status of the EVSE             | "Available", "Charging", "Faulted",...  |
|  session   |          Persistent session data          ||
| powerMeter |    Power Meter configuration for EVSE     || 

#### evcc

|   Attribute   |              Description              |  Possible values   | 
|:-------------:|:-------------------------------------:|:------------------:|
|     type      |         Type of the EVCC used         | "Relay", "Phoenix" |
|   relayPin    |   Attribute specific for the Relay.   |    Any GPIO pin    |
| inverseLogic  |   Attribute specific for the Relay.   |     true,false     |
| deviceAddress | Attribute specific for smarter EVCCs. |  Any string value  |

#### powerMeter

|      Attribute       |            Description            | Possible values | 
|:--------------------:|:---------------------------------:|:---------------:|
|      isEnabled       | Enable or disable the power meter |   true, false   |
|         type         |      Type of the power meter      |    "CS5460A"    |
|    powerMeterPin     |  Attribute specific for CS5460A   |  Any GPIO pin   |
|        spiBus        |  Attribute specific for CS5460A   |       0,1       |
|     consumption      |  Attribute specific for CS5460A   |                 |
|     shuntOffset      |  Attribute specific for CS5460A   |  Default: 1337  |
| voltageDividerOffset |  Attribute specific for CS5460A   |  Default: 0.01  |

#### connectors

|   Attribute   |       Description       |            Possible values            | 
|:-------------:|:-----------------------:|:-------------------------------------:|
|  connectorId  |      connector ID       |          Default:"ChargePi"           |
|     type      |  type of the connector  | "Schuko", "Type1","Type2", "CCS", ... |
|    status     | Status of the connector |     "Available", "Charging", ...      |

Example EVSE configuration:

```json
{
  "evseId": 1,
  "connectors": [
    {
      "connectorId": 1,
      "type": "Schuko",
      "status": "Available"
    }
  ],
  "evcc": {
    "type": "Relay",
    "relayPin": 26,
    "inverseLogic": false
  },
  "status": "Available",
  "session": {
    "isActive": false,
    "transactionId": "",
    "tagId": "",
    "started": "",
    "consumption": []
  },
  "powerMeter": {
    "enabled": false,
    "type": "CS5460A",
    "powerMeterPin": 25,
    "spiBus": 0,
    "consumption": 0.0,
    "shuntOffset": 0.055,
    "voltageDividerOffset": 1333
  }
}
```