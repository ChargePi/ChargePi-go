# Configuration

## Configuring connectivity and basic information of the Charging Point

Settings can be found in the [_settings_](../../configs/settings.json) file and feature basic Charge Point information such as:

- vendor and model,
- unique registered charging point ID, server URI and logging server IP,
- default max charging time,
- OCPP protocol version,
- client current and target version for tracking updates,
- hardware settings for LCD, RFID/NFC reader and LEDs.

The settings file is supported in YAML or JSON format.

The table represents attributes, their values and descriptions that require more attention and might not be
self-explanatory. Some attributes can have multiple possible values, if any are empty, they will be treated as disabled or
might not work properly.

| Attribute| Description |Possible values | 
| :---:    | :---:    | :---:    | 
| id | ID of the charging point. Must be registered in the Central System | Default:"ChargePi" |
| protocolVersion | Version of the OCPP protocol. | "1.6", "2.0.1" |
| serverUri | URI of the Central System with the port and endpoint. | Default: "172.0.1.121:8080/steve/websocket/CentralSystemService" | 
| logServer | IP of the logging server. | Any valid IP | 
| info: maxChargingTime | Max charging time allowed on the Charging point in minutes. | Default:180 |
| rfidReader: readerModel | RFID/NFC reader model used. |  "PN532", ""| 
| ledIndicator: type | Type of the led indicator.  | "WS281x", ""|
| hardware: minPower| Minimum power draw needed to continue charging, if Power meter is configured. | Default:20|

Example settings:

```json
{
  "chargePoint": {
    "info": {
      "vendor": "UL FE",
      "model": "ChargePi",
      "id": "ChargePi",
      "protocolVersion": "1.6",
      "currentClientVersion": "1.0",
      "targetClientVersion": "1.0",
      "serverUri": "yourCSMSURL/",
      "logServer": "yourLoggingServer:12201",
      "maxChargingTime": 180,
      "tls": {
        "isEnabled": false,
        "CACertificatePath": "/usr/share/certs/rootCA.crt",
        "clientCertificatePath": "/usr/share/certs/charge-point.crt",
        "clientKeyPath": "/usr/share/certs/charge-point.key"
      }
    },
    "hardware": {
      "lcd": {
        "isSupported": true,
        "driver": "hd44780",
        "i2cAddress": "0x27"
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
      },
      "powerMeters": {
        "minPower": 20,
        "retries": 3
      }
    }
  }
}
```

## Configuring EVSEs and connectors

### General information

Connector (or equipment, hardware) settings can be found in _/connectors/connector-{id}.json_. To add a connector and
its settings, simply add a new json file that contains the structure below and modify it to your specs. The client will
scan the **connectors** folder at boot and configure the connectors from the files. Charging point can have multiple
EVSE and each can have one or more connectors, but only one connector of the EVSE can charge at a time.

### Attributes

Connector object contains a connector type and an ID of the connector, which must start with 1 and increment by one. The
status attribute changes according to the OCPP specification. The session object represents a Charging session and is
used to restore the connector's last known state when starting the client.

The relay and power meter objects are configurable to specific GPIO pins and SPI bus. The inverseLogic attribute in the
relay object indicates the logic of the relay. If inverseLogic is set to _true_, the relay will use negative logic. The
Power meter also contains some attributes for measurement calibration.

The table represents attributes, their values and descriptions that require more attention and might not be
self-explanatory. Some attributes can have multiple possible values, if any are empty, they will be treated as disabled
or might not work properly.

| Attribute| Description |Possible values | 
| :---:    | :---:    | :---:    | 
| evseId | ID of the EVSE | / |
| type | A type of the connector used in the build. | Refer to OCPP documentation. Default: "Schuko" |
| relay: inverseLogic | Uses negative logic for operating with the relay | false| 
| powerMeter: shuntOffset | Value of the shunt resistor used in the build to measure power. | Default: 0.01 | 
| powerMeter: voltageDividerOffset| Value of the voltage divider used in the build to measure power.| Default:1333 |

Example connector:

```json
{
  "evseId": 1,
  "connectorId": 1,
  "type": "Schuko",
  "status": "Available",
  "session": {
    "isActive": false,
    "transactionId": "",
    "tagId": "",
    "started": "",
    "consumption": []
  },
  "relay": {
    "relayPin": 26,
    "inverseLogic": false
  },
  "powerMeter": {
    "enabled": false,
    "type": "CS5460A",
    "powerMeterPin": 25,
    "spiBus": 0,
    "powerUnits": "kWh",
    "consumption": 0.0,
    "shuntOffset": 0.055,
    "voltageDividerOffset": 1333
  }
}
```