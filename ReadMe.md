# ⚡ ChargePi-go

⚡ ChargePi is an open-source, simple and a configurable Raspberry Pi 4 based ⚡Charging Point🔌 project. It is made with
off-the-shelf hardware and supports multiple EVSEs and simple connectors 🔌. 

A connector🔌 consists of an indicator (🚥 RGB LED strip), a relay and a ⚡ power meter. The client is written in Go and
requires installation of a few C libraries. You can also add an RFID/NFC reader and a display to the Charging Point.

ChargePi-go client can be deployed/run in multiple ways:

- standalone
- Docker 🐳 and docker-compose

## 🔌 Charge point specifications

| OCPP version  | Core functionalities |  Offline charging   | Local authorization | Charging profiles |
|:-------------:|:--------------------:|:-------------------:|:-------------------:|:-----------------:|
|  1.6 JSON/WS  |          ✔️          |     ✔️(partial)     |         ✔️          |         ❌         |
| 2.0.1 JSON/WS | Will be implemented  | Will be implemented | Will be implemented |         ❌         |

### 🛠️ Configuration and settings

To configure the ChargePi client, check out the [configuration guide](/docs/client/configuration.md). Client comes with
default settings which require minimal configuration.

### 📑 Logging

Debugging and insight of the charge point made simple. ChargePi supports multiple logging outputs specifying the format
the [settings](/configs/settings.json) file.

Supported formats:

|            Logging type             | Supported |
|:-----------------------------------:|:---------:|
| [Graylog](https://www.graylog.org/) |     ✔     |
|               Syslog                |     ✔     |
|                File                 |     ✔     |

For details, check out the [logging](/docs/logging/logging.md) docs.

## ⚡ Quickstart

1. Wire your hardware according to the provided [schematics](/docs/hardware/hardware.md).

2. Install the [libraries](/docs/client/installing-libraries.md) (if needed).

3. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/client/configuration.md)
    - [OCPP configuration](/docs/ocpp/ocpp-16.md)

4. Run the client.

Details on how to run the client are described in the [client startup guide](docs/client/running-the-client.md).