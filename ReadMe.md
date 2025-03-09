# ⚡ ChargePi-go

ChargePi is an open-source, Linux based ⚡Charging Point🔌 "framework". It makes your EV charging station smart by
abstracting and unifying the hardware components into a customizable,
modular Charge Point with OCPP 1.6 and 2.0.1 support, an intuitive UI for easy management and an API for custom
integrations. It is highly customizable, simple to make and supports both off-the-shelf and commercial hardware.

## Concepts

A charge point is a device, consisting of multiple EVSEs, that is capable of delivering energy to an electric vehicle.
Each EVSE has a connector, an EV Charge Controller and a Power Meter. Optionally, the charge point can have an RFID/NFC
reader as an authentication mechanism, a display as a user interface, a status indicator and a communication module for
connecting to the backend.

## 🔌 Charge point software specification

| OCPP version  | Core functionalities | Reservations | LocalAuthList | SmartCharging | FirmwareUpdate |       Security Extension        | ISO 15118 (Plug and Charge) |
|:-------------:|:--------------------:|:------------:|:-------------:|:-------------:|:--------------:|:-------------------------------:|:---------------------------:|
|  1.6 JSON/WS  |          ✔️          |      ✔️      |      ✔️       |    Planned    |  ✔️ (partial)  |             Planned             |           Planned           |
| 2.0.1 JSON/WS |       Planned        |   Planned    |    Planned    |    Planned    |    Planned     | Certificate management included |          Included           |

## 📦 Features

- **OCPP 1.6** - ChargePi supports OCPP 1.6, the most common protocol for EV charging
  stations. OCPP 2.0.1 support is planned.
- **Easily configurable** - ChargePi is easy to configure and customize. It comes with default settings that require
  minimal configuration.
- **Easy replication** - ChargePi is easy to replicate. It can be deployed on a Raspberry Pi or any other Linux-based
  device. Configurations can be easily exported from an existing device to other devices.
- **Hardware-agnostic** - ChargePi supports a variety of hardware components. You can mix and match the hardware
  components to create a charging station that suits your needs. You can also write custom drivers for currently
  unsupported hardware. We are always looking for contributors to add support for new hardware components.
- **Modular** - ChargePi is modular, allowing you to add or remove components as needed, based on your use case. You can
  easily add new hardware components or functionalities to your charging station.
- **API** - ChargePi comes with an API that allows you to integrate it with other systems or implement custom
  functionality.

## Roadmap

- **OCPP 1.6 Extension** - Security extension for OCPP 1.6, which adds additional security features to the protocol.
- **OCPP 2.0.1 Support** - Support for OCPP 2.0.1, the latest version of the OCPP protocol.
- **Simulator** - ChargePi offers a simulator mode, which allows you to simulate events such as charging, stopping
  charging, and RFID authentication. Mostly used for testing purposes.
- **Demo mode** - a functional charging station that can be used for demonstration purposes, runs with emulated
  hardware.
- **Smart charging** - Smart charging functionality, which allows you to optimize the charging process
  based on factors such as energy prices, grid load, and user preferences.
- **Customizable UI** - ChargePi comes with a customizable UI that allows you to easily manage your charging station.

## 🔌 Hardware support

All contributions for hardware support are welcome. If you want to add support for a new hardware component, please
check out the [hardware support guide](/docs/hardware/hardware-support.md). The more hardware is supported, the more
customizable the ChargePi is.

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

2. Wire your hardware according to the provided [schematics](/docs/hardware/hardware-support.md) (if the hardware is
   supported - otherwise, seek professional guidance or contribute to the project!).

   > [!CAUTION]
   > The hardware part of the project requires fairly good electrical knowledge. Make at your own risk or seek
   professional
   guidance. The project is not responsible for any damage caused by the hardware and will not take responsibility for
   any software misconfiguration or misbehavior.

3. Configure the settings files according to the hardware and desired functionality:
    - [client configuration](/docs/configuration/configuration.md)
    - [OCPP configuration](/docs/configuration/ocpp/ocpp-16.md)

4. Run the client. ChargePi-go can be deployed/run in multiple ways:

    - standalone (host machine)
    - Docker 🐳 and docker-compose

Details on how to run the client on are described in
the [client startup guide](docs/getting-started/installation/running-the-client.md).

### 🛠️ Configuration and settings

To configure the ChargePi, check out the [configuration guide](/docs/configuration/configuration.md).
Client comes with default settings which require minimal configuration.

## Contributing

Contributions are welcome! Whether you want to add support for new hardware components, improve the documentation,
found a bug or have an idea for a new feature, we appreciate your help!

Please check out the [contributing guide](/docs/contribution/contributing.md) for more information on how to contribute.

## License

ChargePi-go is licensed under the [MIT License](LICENSE.txt).