# ChargePi-go

ChargePi is an open-source Raspberry Pi 4 based Charging Point project, which supports multiple EVSEs and simple
connectors. A charging connector consists of a WS2811 LED strip, a relay and a power meter. It is written in Golang and
requires installation of a few C libraries.

ChargePi client can be deployed/run in multiple ways:

- standalone
- Docker by building the image and running the container
- Docker-compose to be deployed with SteVe Central System and Watchtower (**recommended for testing only**)
- Docker-compose by running the client

### Charging station specifications

| Protocol implementation | Core functionalities | Offline charging | Local authorization | Charging profiles |
| :---:    | :---:    | :---:    |:---:    | :---:    |
| OCPP 1.6 JSON/WS | Yes | Yes | Yes | No |
| OCPP 2.0.1 JSON/WS | Will be implemented | Will be implemented | Will be implemented | No |

### Configuration and settings

To configure the ChargePi client, check out the **configuration.md** guide under __docs/configuration__. Client comes with predefined
settings which require minimal configuration.

### Logging

ChargePi uses [Graylog](https://www.graylog.org/) logging server for remote logging, so a server should be up and
running if you want the logs to be stored (check setup instructions for Graylog under _services/logging_). Logs are in
GELF format and are sent through UDP protocol. The library used for sending logs
is [go-gelf](https://github.com/Graylog2/go-gelf).

Configure the **"logServer"** property in the _settings.json_ file with your server IP/domain name with the port.

## Initial setup

1. If you want to run SteVe on the same host:

> git clone https://github.com/RWTH-i5-IDSG/steve

*_When cloning Steve from GitHub, steve directory should be automatically generated._
Replace SteVe's default Dockerfile with Dockerfile provided in ChargePi/steve/Dockerfile to run on Raspberry Pi.

2. Wire your hardware according to the provided schematics. *[Useful reference for Raspberry Pi](https://pinout.xyz/)*

3. Configure the settings files according to the hardware and desired functionality. Check the **docs/configuration**
   folder for configuration instructions.

### Building libnfc for PN532

1. Get and extract the libnfc:

```
 cd ~
 mkdir libnfc && cd libnfc/
 wget https://github.com/nfc-tools/libnfc/releases/download/libnfc-1.8.0/libnfc-1.8.0.tar.bz2
 tar -xvjf libnfc-1.8.0.tar.bz2
```

**Next two steps may vary for your reader**

2. Create PN532 UART configuration:

```
 cd libnfc-1.8.0
 sudo mkdir /etc/nfc
 sudo mkdir /etc/nfc/devices.d
 sudo cp contrib/libnfc/pn532_uart_on_rpi.conf.sample /etc/nfc/devices.d/pn532_uart_on_rpi.conf 
 sudo nano /etc/nfc/devices.d/pn532_uart_on_rpi.conf
```

3. Add the line at the end of the file:

   > allow_intrusive_scan = true

4. Install dependencies and configure:

```
 sudo apt-get install autoconf libtool libpcsclite-dev libusb-dev
 autoreconf -vis
 ./configure --with-drivers=pn532_uart --sysconfdir=/etc --prefix=/usr
```

5. Build the library:

```
sudo make clean
sudo make install all
```

### Installing rpi-ws281x library

Follow the instructions on [maintainer's repository](https://github.com/jgarff/rpi_ws281x).

**TLDR; Use the instructions here**

```
git clone https://github.com/jgarff/rpi_ws281x
cd rpi_ws281x
mkdir build
cd build
cmake -D BUILD_SHARED=OFF -D BUILD_TEST=ON ..
cmake --build .
sudo make install
```

To be able to use this C library in Go, it must be installed. You can do this by copying `*.h` to `/usr/local/include`
and `'.a` files to `/usr/local/lib`. If not working, check the Go wrapper library
instructions [here](https://github.com/rpi-ws281x/rpi-ws281x-go).

## Running standalone

ChargePi client uses **[SteVe](https://github.com/RWTH-i5-IDSG/steve)** for the Central System, but can connect to other
Central Systems as well.
**Optional: Clone the repository and run SteVe.**

Running the client in Golang:

   > go run main.go

or compiling and executing the client:

```
go build chargepi
./chargepi
```

## Deploying on Docker

1. Build the client image on Docker:

   > cd ChargePi/client
   > docker build -t chargepi .

2. Run the container from built image:

   > docker run --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --privileged chargepi

### Deploying using docker-compose

Alternatively, you can run the client, SteVe server and Watchtower on the same Pi using **docker-compose**.
The **[Watchtower](https://github.com/containrrr/watchtower)** service will automatically pull the newest image and run
it when it is available.

1. Change the IP address under __serverUri__ in the settings file to **172.0.1.121**.

2. Build services:

   > docker-compose build

3. Run services in daemon mode:

   > docker-compose up -d

## Helpful references

- [rpi_ws281x C library](https://github.com/jgarff/rpi_ws281x)
- [libnfc](https://github.com/nfc-tools/libnfc)
- [Go wrapper for rpi_ws281x](https://github.com/rpi-ws281x/rpi-ws281x-go)
- [Go wrapper for libnfc](https://github.com/clausecker/nfc)
- [CS5460A Implementation inspired by this project](https://github.com/cbm80amiga/ST7789_power_meter_cs5460a_display/)
- [CS5460A Datasheet](https://statics.cirrus.com/pubs/proDatasheet/CS5460A_F5.pdf)
- [SteVe, the open-source Central System](https://github.com/RWTH-i5-IDSG/steve)
- [OCPP 1.6 specification](https://www.oasis-open.org/committees/download.php/58944/ocpp-1.6.pdf)
- [OCPP 2.0.1 specification](https://github.com/mobilityhouse/ocpp/tree/master/docs/v201)
- [Docker](https://docs.docker.com/)
- [Installing Docker on Pi](https://www.docker.com/blog/happy-pi-day-docker-raspberry-pi/)
- [Watchtower](https://github.com/containrrr/watchtower)
- [Raspberry Pi pinout](https://pinout.xyz/)
- [Mender Docs](https://docs.mender.io/get-started/preparation/prepare-a-raspberry-pi-device)
- [GoCron](https://github.com/go-co-op/gocron)
- [go-gelf](https://github.com/Graylog2/go-gelf)