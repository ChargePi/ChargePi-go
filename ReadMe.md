# ChargePi-go

⚡ChargePi⚡ is an open-source Raspberry Pi 4 based ⚡Charging Point🔌 project, which supports multiple EVSEs and simple
connectors🔌. A charging connector🔌 consists of a WS281x RGB 🚥 LED strip, a relay and a power meter. The client is
written in Golang and requires installation of a few C libraries.

ChargePi client can be deployed/run in multiple ways:

- standalone
- Docker 🐳 by building the image and running the container
- Docker-compose to be deployed with SteVe Central System and Watchtower (**recommended for dev/testing only**)
- Docker-compose by running the client

### Charging station specifications

| Protocol implementation | Core functionalities | Offline charging | Local authorization | Charging profiles |
| :---:    | :---:    | :---:    |:---:    | :---:    |
| OCPP 1.6 JSON/WS | ✔️ | ✔️ | ✔️ | ❌ |
| OCPP 2.0.1 JSON/WS | Will be implemented | Will be implemented | Will be implemented | ❌ |

### Configuration and settings

To configure the ChargePi client, check out the [guide](/docs/client/configuration.md). Client comes with predefined
settings which require minimal configuration.

### Logging

ChargePi uses [Graylog](https://www.graylog.org/) logging server for remote logging, so a server should be up and
running if you want the logs to be stored (check setup [instructions](/docs/services/graylog.md)). Logs are sent through
UDP protocol in GELF format. The library used for sending logs is [go-gelf](https://github.com/Graylog2/go-gelf).

Configure the `logServer` property in the [settings](/configs/settings.json) file with your server IP/domain name with
the port.

## Initial setup

1. If you want to run SteVe on the same host:

   ```bash
   git clone https://github.com/RWTH-i5-IDSG/steve
   ```

   *_When cloning Steve from GitHub, steve directory should be automatically generated._
   Replace SteVe's default Dockerfile with Dockerfile provided [here](build/package/Dockerfile-steve) to run on
   Raspberry Pi.

2. Wire your hardware according to the provided [schematics](/docs/hardware/hardware.md).

3. Install the [libraries](/docs/client/installing-libraries.md) (if needed).

4. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/client/configuration.md)
    - [OCPP configuration](/docs/ocpp/ocpp-16.md)

5. Run the client.

## Running standalone

Running the client:

   ```bash
   go run .
   ```

or compiling and executing the client:

   ```bash
   go build -o chargepi .
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

1. Change the IP address under `serverUri` in the settings file to **172.0.1.121**.

2. Build services:

    ```bash
    docker-compose build
    ```

3. Run services in daemon mode:

   ```bash
   docker-compose up -d
   ```