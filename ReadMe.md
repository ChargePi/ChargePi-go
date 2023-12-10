# ⚡ ChargePi-go

## What is ChargePi-go?

ChargePi is an open-source, Linux based ⚡Charging Point🔌 project with OCPP support with an API and UI.
It is highly customizable, simple to make and supports both off-the-shelf and commercial hardware.

## Concepts

A charge point is a device, consisting of multiple EVSEs, that is capable of delivering energy to an electric vehicle.
Each EVSE has a connector, an EV Charge Controller and a Power Meter. Optionally, the charge point can have an RFID/NFC
reader as an authentication mechanism, a display as a user interface, a status indicator and a communication module for
connecting to the backend.

ChargePi abstracts and unifies these components into a customizable, modular Charge Point with OCPP 1.6 and 2.0.1
support,an intuitive UI for easy management and an API for custom integrations.

## 🔌 Charge point software specifications

| OCPP version  | Core functionalities | Reservations | LocalAuthList | SmartCharging | FirmwareUpdate |
|:-------------:|:--------------------:|:------------:|:-------------:|:-------------:|:--------------:|
|  1.6 JSON/WS  |          ✔️          |      ✔️      |      ✔️       |    Planned    |  ✔️ (partial)  |
| 2.0.1 JSON/WS |       Planned        |   Planned    |    Planned    |    Planned    |    Planned     |

## 🔌 Hardware support

All contributions for hardware support are welcome. If you want to add support for a new hardware component, please
check out the [hardware support guide](/docs/getting-started/hardware/hardware-support.md). The more hardware is
supported, the more customizable the ChargePi is.

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

1. Choose the desired functionality and hardware.

2. Wire your hardware according to the provided [schematics](/docs/getting-started/hardware/hardware-support.md) (if
   applicable).

3. If deploying on a host machine, install the [dependencies](/docs/getting-started/installation/libraries.md).

4. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/getting-started/configuration/configuration.md)
    - [OCPP configuration](/docs/getting-started/configuration/ocpp/ocpp-16.md)

5. Run the client. ChargePi-go can be deployed/run in multiple ways:

    - standalone (host machine)
    - Docker 🐳 and docker-compose

Details on how to run the client on are described in
the [client startup guide](docs/getting-started/installation/running-the-client.md).

### 🛠️ Configuration and settings

To configure the ChargePi, check out
the [configuration guide](/docs/getting-started/configuration/configuration.md).
Client comes with default settings which require minimal configuration.

## Note

The hardware part of the project requires fairly good electrical knowledge. Make at your own risk or seek professional
guidance. The project is not responsible for any damage caused by the hardware and will not take responsibility for any
software misconfiguration or misbehavior.