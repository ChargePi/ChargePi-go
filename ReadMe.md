# ⚡ ChargePi-go

⚡ChargePi is an open-source Raspberry Pi 4 based ⚡Charging Point🔌 project, which supports multiple EVSEs and simple
connectors🔌. You can also add an RFID/NFC reader and a display to the Charging Point. It is configurable, simple to
make and uses off-the-shelf hardware.

A connector🔌 consists of an indicator (🚥 RGB LED strip), a relay and a ⚡ power meter. The client is written in Go and
requires installation of a few C libraries.

ChargePi-go client can be deployed/run in multiple ways:

- standalone
- Docker 🐳 by building the image and running the container
- Docker-compose with SteVe Central System and Watchtower (**recommended for testing only**)
- Docker-compose (client only)

## 🔌 Charge point specifications

| OCPP version  | Core functionalities |  Offline charging   | Local authorization | Charging profiles |
|:-------------:|:--------------------:|:-------------------:|:-------------------:|:-----------------:|
|  1.6 JSON/WS  |          ✔️          |     ✔️(partial)     |         ✔️          |         ❌         |
| 2.0.1 JSON/WS | Will be implemented  | Will be implemented | Will be implemented |         ❌         |

### 🛠️ Configuration and settings

To configure the ChargePi client, check out the [configuration guide](/docs/client/configuration.md). Client comes with
default settings which require minimal configuration.

### 📑 Logging

ChargePi supports multiple logging outputs. Configure the `logging` property in the [settings](/configs/settings.json)
file. Check out the list below:

|            Logging type             | Supported |
|:-----------------------------------:|:---------:|
| [Graylog](https://www.graylog.org/) |     ✔     |
|               Syslog                |     ✔     |
|                File                 |     ✔     |

For details, check out the [logging](/docs/logging/logging.md) docs.

## ⚡ Quickstart

1. If you want to run SteVe on the same host (_Not recommended in production_):

   ```bash
   git clone https://github.com/RWTH-i5-IDSG/steve
   ```

   Replace SteVe's default Dockerfile with Dockerfile provided [here](build/SteVe/Steve.Dockerfile) to run on
   Raspberry Pi.

2. Wire your hardware according to the provided [schematics](/docs/hardware/hardware.md).

3. Install the [libraries](/docs/client/installing-libraries.md) (if needed).

4. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/client/configuration.md)
    - [OCPP configuration](/docs/ocpp/ocpp-16.md)

5. Run the client.

## 🏃 Running standalone

Running the client:

   ```bash
   go run .
   ```

or compiling and executing the client:

   ```bash
   go build -o chargepi .
   ./chargepi
   ```

More options available in the [client startup guide](docs/client/running-the-client.md).