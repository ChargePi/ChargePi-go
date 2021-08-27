# Supported hardware and schematics

This client supports following hardware:

- Raspberry Pi 3B/4B
- Relay(s) 230V 10A
- PN532 RFID/NFC reader or any reader supported by [libnfc](http://nfc-tools.org/index.php/Libnfc)
- Power meter (CS5460A chip)
- LCD (optionally with PCF8574 I2C module)
- WS281x LED strip

Hardware must be configured in [_settings file_](../../configs/settings.json) and [_connectors_](../../configs/connectors) files.

## RFID/NFC reader

PN532 communicates through UART and the program uses the NFC go library, which is a wrapper for libnfc 1.8.0. You could
use any other libnfc compatible NFC/RFID reader, but the configuration steps may vary.

| RPI PIN |   PN532 PIN    | 
| :---:	| :---:	|
|  5V    |  VCC  |
|   GND    |  GND    | 
|   GPIO 14    |  TX    |
|   GPIO 15    |  RX    | 

## LCD I2C

LCD should be on I2C bus 1 with address 0x27. To find the I2C address, follow these steps:

1. Download i2c tools:

   ```bash
   sudo apt-get install -y i2c-tools
   ```

2. If needed, reboot.

3. Run the following command to get the I2C address:

   ```bash
   sudo i2cdetect -y 1 
   ```

| RPI PIN |   PCF8574 PIN    | 
| :---:	| :---:	|
|   2 or any 5V pin    |  VCC  |
|   14 or any ground pin    |  GND    | 
|   3 (GPIO 2)    |  SDA    |
|   5 (GPIO 3)    |  SCL    | 

## Relay (or relay module)

It is highly recommended splitting both GND and VCC between relays or using a relay module.

| RPI PIN |  RELAY PIN    | 
| ---	| :---:	|
|   4 or any 5V pin    |   VCC    | 
|   20 or any ground pin    |   GND    |  
|  37 (GPIO 26) or any free GPIO pin    |   S/Enable    |  

## Power meter

| RPI PIN|  CS5460A PIN    |  RPI PIN |   CS5460A PIN    |
| :---:	| :---:	| :---:	| :---:	|
|   4 or 2    |   VCC    |  38 (GPIO 20)    |   MOSI    |
|   25 or any ground pin    |   GND    |   35 (GPIO 19)    |   MISO    |
|   Any free pin    |   CE/SDA    |   /    |   /    |
|   40 (GPIO 21)    |   SCK    |   /    |  /    |

## WS281x LED strip

| RPI PIN|  WS281x PIN    |  RPI PIN |   WS281x PIN    |
| :---:	| :---:	| :---:	| :---:	|
|   External 12V    |   VCC    |  32 (GPIO 12)    |   Data |
|   External GND   |   GND    |   /    |  / |

## Wiring diagram

![](WiringSketch_eng.png)