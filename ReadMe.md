# ChargePi-go

ChargePi is an open-source Raspberry Pi 4 based Charging Point project, which supports multiple EVSEs and simple
connectors. A charging connector consists of a WS2811 LED strip, a relay and a power meter. It is written in Golang and
requires installation of a few C libraries.

ChargePi client can be deployed/run in multiple ways:

- standalone
- Docker by building the image and running the container
- Docker-compose to be deployed with SteVe Central System and Watchtower (**recommended for dev/testing only**)
- Docker-compose by running the client

### Charging station specifications

| Protocol implementation | Core functionalities | Offline charging | Local authorization | Charging profiles |
| :---:    | :---:    | :---:    |:---:    | :---:    |
| OCPP 1.6 JSON/WS | Yes | Yes | Yes | No |
| OCPP 2.0.1 JSON/WS | Will be implemented | Will be implemented | Will be implemented | No |

### Configuration and settings

To configure the ChargePi client, check out the [guide](/docs/client/configuration.md). Client comes with predefined
settings which require minimal configuration.

### Logging

ChargePi uses [Graylog](https://www.graylog.org/) logging server for remote logging, so a server should be up and
running if you want the logs to be stored (check setup [instructions](/docs/services/graylog.md)). Logs are sent through
UDP protocol in GELF format. The library used for sending logs is [go-gelf](https://github.com/Graylog2/go-gelf).

Configure the **"logServer"** property in the [settings](/configs/settings.json) file with your server IP/domain name
with the port.

## Initial setup

1. If you want to run SteVe on the same host:

   ```bash
   git clone https://github.com/RWTH-i5-IDSG/steve
   ```

   *_When cloning Steve from GitHub, steve directory should be automatically generated._
   Replace SteVe's default Dockerfile with Dockerfile provided in _build/package/Dockerfile-steve_ to run on Raspberry
   Pi.

2. Wire your hardware according to the provided [schematics](/docs/hardware/hardware.md).

3. Install the [libraries](/docs/client/installing-libraries.md) (if needed).

4. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/client/configuration.md)
    - [OCPP configuration](/docs/ocpp/ocpp-16.md)

5. Run the client.

## Running standalone

Running the client in Golang:

   ```bash
   go run main.go
   ```

or compiling and executing the client:

   ```bash
   go build chargepi
   ./chargepi
   ```

## Deploying on Docker

1. Build the client image on Docker:

   ```bash
   cd ChargePi-go/
   docker build -t chargepi .
   ```

2. Run the container from built image:

   ```bash
   docker run --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --privileged chargepi
   ```

### Deploying using docker-compose

Alternatively, you can run the client, SteVe server and Watchtower on the same Pi using **docker-compose**.
The **[Watchtower](https://github.com/containrrr/watchtower)** service will automatically pull the newest image and run
it when it is available.

1. Change the IP address under __serverUri__ in the settings file to **172.0.1.121**.

2. Build services:

    ```bash
    docker-compose build
    ```

3. Run services in daemon mode:

   ```bash
   docker-compose up -d
   ```

## Libraries and guides

- [rpi_ws281x C library](https://github.com/jgarff/rpi_ws281x)
- [libnfc](https://github.com/nfc-tools/libnfc)
- [Go wrapper for rpi_ws281x](https://github.com/rpi-ws281x/rpi-ws281x-go)
- [Go wrapper for libnfc](https://github.com/clausecker/nfc)
- [CS5460A Implementation inspired by this project](https://github.com/cbm80amiga/ST7789_power_meter_cs5460a_display/)
- [CS5460A Datasheet](https://statics.cirrus.com/pubs/proDatasheet/CS5460A_F5.pdf)
- [SteVe, the open-source Central System](https://github.com/RWTH-i5-IDSG/steve)
- [Docker](https://docs.docker.com/)
- [Installing Docker on Pi](https://www.docker.com/blog/happy-pi-day-docker-raspberry-pi/)
- [Watchtower](https://github.com/containrrr/watchtower)
- [Raspberry Pi pinout](https://pinout.xyz/)
- [GoCron](https://github.com/go-co-op/gocron)
- [go-gelf](https://github.com/Graylog2/go-gelf)
- [ocpp-go](https://github.com/lorenzodonini/ocpp-go)