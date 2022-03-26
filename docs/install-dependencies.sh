#!/bin/bash
# Update
sudo apt-get update -y
sudo apt-get upgrade -y

# Download libnfc
cd ~ && mkdir libnfc
cd libnfc/
wget https://github.com/nfc-tools/libnfc/releases/download/libnfc-1.8.0/libnfc-1.8.0.tar.bz2
tar -xvjf libnfc-1.8.0.tar.bz2
cd libnfc-1.8.0 && sudo mkdir /etc/nfc /etc/nfc/devices.d

# Configure and install libnfc
sudo apt-get install autoconf libtool libpcsclite-dev libusb-dev -y
autoreconf -vis

if [ "$1" = "pn532_i2c" ]; then
  # Add configuration for PN532_I2C
  touch /etc/nfc/devices.d/pn532_i2c.conf
  echo name = "PN532 board via I2C" >>/etc/nfc/devices.d/pn532_i2c.conf
  echo connstring = pn532_i2c:/dev/i2c-1 >>/etc/nfc/devices.d/pn532_i2c.conf
  echo allow_intrusive_scan = true >>/etc/nfc/devices.d/pn532_i2c.conf
  ./configure --with-drivers=pn532_i2c --sysconfdir=/etc --prefix=/usr

elif [ "$1" = "pn532_spi" ]; then
  # Add configuration for PN532_SPI
  touch /etc/nfc/devices.d/pn532_spi.conf
  echo name = "PN532 board via SPI" >>/etc/nfc/devices.d/pn532_spi.conf
  echo connstring = pn532_i2c:/dev/spidev0.0:500000 >>/etc/nfc/devices.d/pn532_spi.conf
  echo allow_intrusive_scan = true >>/etc/nfc/devices.d/pn532_spi.conf
  ./configure --with-drivers=pn532_spi --sysconfdir=/etc --prefix=/usr

elif [ "$1" = "pn532_uart" ]; then
  # Add configuration for PN532_UART
  touch /etc/nfc/devices.d/pn532_uart.conf
  echo name = "PN532 board via UART" >>/etc/nfc/devices.d/pn532_uart.conf
  echo connstring = pn532_uart:/dev/ttyAMA0 >>/etc/nfc/devices.d/pn532_uart.conf
  echo allow_intrusive_scan = true >>/etc/nfc/devices.d/pn532_uart.conf
  ./configure --with-drivers=pn532_uart --enable-serial-autoprobe --sysconfdir=/etc --prefix=/usr

fi

sudo make clean
sudo make install all

# Install WS281x drivers
sudo apt-get install cmake -y
cd ~
git clone https://github.com/jgarff/rpi_ws281x
cd rpi_ws281x && mkdir build
cd build
cmake -D BUILD_SHARED=OFF -D BUILD_TEST=ON ..
cmake --build .
sudo make install
sudo cp *.a /usr/local/lib
sudo cp *.h /usr/local/include

# Optionally, install Go
if [ "$2" -eq 1 ]; then
  # installing Go with the help of this script: https://github.com/canha/golang-tools-install-script
  wget -q -O - https://git.io/vQhTU | bash
fi
