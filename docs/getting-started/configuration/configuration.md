# 🛠️ Configuration

There are three types of configurations:

1. [charge point](#the-charge-point-configuration)
2. [ocpp](../configuration/ocpp/ocpp-16.md#ocpp-16)
3. [evse](evse-configuration.md#the-evse-configuration)

## The `charge point` configuration

The charge point configuration contains basic information about the charge point and its operative details, such as:

- Charge Point ID,
- Connectivity details (central system URI and OCPP protocol version, etc.)
- OCPP information (vendor, model, firmware, etc.),
- logging settings,
- hardware settings

The following tables represent attributes, their values and descriptions that require more attention and might not be
self-explanatory. Some attributes can have multiple possible values, if any are empty, they will be treated as disabled
or might not work properly.

### chargePoint connectionSettings

|    Attribute    |                            Description                             |  Possible values   | 
|:---------------:|:------------------------------------------------------------------:|:------------------:|
|       id        | ID of the charging point. Must be registered in the Central System | Default:"ChargePi" |
| protocolVersion |                   Version of the OCPP protocol.                    |   "1.6", "2.0.1"   |
|    serverUri    |       URI of the Central System with the port and endpoint.        |         ""         |
|  basicAuthUser  |      HTTP username for authentication with the Central System      |  Any string value  |
|  basicAuthPass  |      HTTP username for authentication with the Central System      |  Any string value  |
|       tls       |                       TLS certificate paths                        |                    |

### chargePoint connectionSettings tls

|       Attribute       |        Description        | Possible values | 
|:---------------------:|:-------------------------:|:---------------:|
|   CACertificatePath   | Root/CA certificate path. | Any valid path  |
| clientCertificatePath |    Client certificate     | Any valid path  |
|     clientKeyPath     |    Client private key     | Any valid path  |

### chargePoint info

|    Attribute    |                            Description                            |      Possible values       | 
|:---------------:|:-----------------------------------------------------------------:|:--------------------------:|
| maxChargingTime | Maxiumum amount of time that a transaction can last (in minutes). |        Default: 180        |
|   ocpp.vendor   |                             Vendor ID                             |     Default: "xBlaz3k"     |
|   ocpp.model    |                               Model                               |    Default: "ChargePi"     |
|      type       |                     Type of the Charge point.                     |        "AC" or "DC"        |
|    maxPower     |     Maximum power the charge point can supply to EVs. In kWh      | Must be greater than zero. |

### chargePoint hardware

#### chargePoint hardware display

| Attribute  |                      Description                       | Possible values | 
|:----------:|:------------------------------------------------------:|:---------------:|
| isEnabled  |               Enable or disable display                |   true, false   |
|   driver   |                  Display driver type                   |       ""        |
| i2cAddress | Field specific for any display using I2C communication |       ""        |
|   i2cBus   | Field specific for any display using I2C communication |    "1", "0"     | 
|  language  |   Language selection for messages displayed to LCD.    |    "en",sl"     |

#### chargePoint hardware tagReader

|   Attribute   |                         Description                         |                         Possible values                          | 
|:-------------:|:-----------------------------------------------------------:|:----------------------------------------------------------------:|
|   isEnabled   |                 Enable or disable tagReader                 |                        Default:"ChargePi"                        |
|  readerModel  |                Version of the OCPP protocol.                |                          "1.6", "2.0.1"                          |
|   resetPin    |    URI of the Central System with the port and endpoint.    | Default: "172.0.1.121:8080/steve/websocket/CentralSystemService" |
| deviceAddress | Max charging time allowed on the Charging point in minutes. |                           Default:180                            |

#### chargePoint hardware ledIndicator

|    Attribute     |                         Description                         |                         Possible values                          | 
|:----------------:|:-----------------------------------------------------------:|:----------------------------------------------------------------:|
|    isEnabled     |               Enable or disable ledIndicator                |                        Default:"ChargePi"                        |
|       type       |                Version of the OCPP protocol.                |                          "1.6", "2.0.1"                          |
|     dataPin      |    URI of the Central System with the port and endpoint.    | Default: "172.0.1.121:8080/steve/websocket/CentralSystemService" |
| indicateCardRead | Max charging time allowed on the Charging point in minutes. |                           Default:180                            |
|      invert      |                 RFID/NFC reader model used.                 |                           "PN532", ""                            | 

Example settings:

```yaml
api:
  enabled: true
  address: 0.0.0.0:8080

chargePoint:
  connectionSettings:
    id: ChargePi
    protocolVersion: '1.6'
    serverUri: example.com
    basicAuthUser: ''
    basicAuthPass: ''
    tls:
      isEnabled: false
      CACertificatePath: /usr/share/certs/rootCA.crt
      clientCertificatePath: /usr/share/certs/charge-point.crt
      clientKeyPath: /usr/share/certs/charge-point.key

  info:
    type: AC
    maxPower: 11
    maxChargingTime: 5
    ocpp:
      vendor: UL FE
      model: ChargePi

  logging:
    type:
      - remote
      - file
    format: gelf
    host: logging.example.com
    port: 12201

  hardware:
    display:
      enabled: true
      driver: hd44780
      i2c:
        address: '0x27'
        bus: 1
      language: en
    reader:
      enabled: true
      readerModel: PN532
      device: /dev/ttyS0
      resetPin: 19
    indicator:
      enabled: true
      type: WS281x
      dataPin: 18
      indicateCardRead: true
      invert: false
```

## Importing

## Exporting