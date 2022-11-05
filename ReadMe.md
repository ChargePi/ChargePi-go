# ⚡ ChargePi-go

ChargePi is an open-source Linux based OCPP-enabled ⚡Charging Point🔌 project. It is configurable, simple to make and
uses off-the-shelf hardware.

It supports multiple connectors🔌, which consist of an indicator (🚥 RGB LED strip), a EV Charge Controller and a ⚡ power
meter. Optionally, the charge point supports an RFID/NFC reader and a display as an user interface.

ChargePi-go client can be deployed/run in multiple ways:

- standalone
- Docker 🐳 and docker-compose

## 🔌 Charge point specifications

| OCPP version  | Core functionalities |    Reservations     |    LocalAuthList    | SmartCharging | FirmwareUpdate |
|:-------------:|:--------------------:|:-------------------:|:-------------------:|:-------------:|:--------------:|
|  1.6 JSON/WS  |          ✔️          |     ✔️(partial)     |         ✔️          |       ❌       |       ❌        |
| 2.0.1 JSON/WS | Will be implemented  | Will be implemented | Will be implemented |       ❌       |                |

## ⚡ Quickstart

1. Wire your hardware according to the provided [schematics](/docs/hardware/hardware.md).

2. Install the [libraries](/docs/client/installing-libraries.md) (if needed).

3. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/client/configuration.md)
    - [OCPP configuration](/docs/ocpp/ocpp-16.md)

4. Run the client.

Details on how to run the client are described in the [client startup guide](docs/client/running-the-client.md).

### 🛠️ Configuration and settings

To configure the ChargePi client, check out the [configuration guide](/docs/client/configuration.md). Client comes with
default settings which require minimal configuration.
