# ⚡ ChargePi-go

ChargePi is an open-source, Linux based ⚡Charging Point🔌 project with OCPP support with an API and UI.
It is highly customizable, simple to make and supports both off-the-shelf and commercial hardware.

Multiple connectors/EVSEs🔌 can be configured for the Charge point, which consist of an indicator (🚥 RGB LED strip),
an EV Charge Controller and a ⚡Power Meter. Optionally, the charge point supports an RFID/NFC reader and a display as a
user interface.

## 🔌 Charge point specifications

| OCPP version  | Core functionalities | Reservations | LocalAuthList | SmartCharging | FirmwareUpdate |
|:-------------:|:--------------------:|:------------:|:-------------:|:-------------:|:--------------:|
|  1.6 JSON/WS  |          ✔️          |      ✔️      |      ✔️       |    Planned    |       ❌        |
| 2.0.1 JSON/WS |       Planned        |   Planned    |    Planned    |    Planned    |    Planned     |

## 🔌 Hardware support

### RFID/NFC readers

| Reader | Is supported | 
|:------:|:------------:|
| PN532  |      ✔       |

### Display

| Display | Is supported | 
|:-------:|:------------:|
| HD44780 |      ✔       |

### Power Meters

| Power meter | Is supported | 
|:-----------:|:------------:|
|   CS5460A   |      ✔       |
|     ETI     |   Planned    |

### EVCC

EV charging controller (EVCC) controls the communication with the EV and allows or denies the charging. It can also set
the charging current limit.

|         EVCC          | Is supported | 
|:---------------------:|:------------:|
|         Relay         |      ✔       |
| Phoenix Contact EVSEs |   Planned    |

### LED indicators

| Indicator | Is supported | 
|:---------:|:------------:|
|  WS2812b  |      ✔       |
|  WS2811   |      ✔       |

## ⚡ Quickstart

1. Wire your hardware according to the provided [schematics](/docs/getting-started/hardware/hardware-support.md).

2. Install the [libraries](/docs/getting-started/installation/libraries.md) (if needed).

3. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/getting-started/configuration/configuration.md)
    - [OCPP configuration](/docs/getting-started/configuration/ocpp/ocpp-16.md)

4. Run the client. ChargePi-go can be deployed/run in multiple ways:

    - standalone
    - Docker 🐳 and docker-compose

Details on how to run the client are described in the [client startup guide](docs/getting-started/installation/running-the-client.md).

### 🛠️ Configuration and settings

To configure the ChargePi client, check out the [configuration guide](/docs/getting-started/configuration/configuration.md).
Client comes with default settings which require minimal configuration.

## Note

The hardware part of the project requires fairly good electrical knowledge. Make at your own risk or seek professional
guidance.