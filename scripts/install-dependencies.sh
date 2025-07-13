#!/bin/sh

# Prepare the system for building libnfc and WS281x drivers
apt-get update -y && apt-get install -y \
    pkg-config \
    build-essential \
    cmake \
    make \
    autoconf \
    libtool \
    libpcsclite-dev \
    libusb-dev \
    bzip2 \
    git \
    wget \
    gcc-aarch64-linux-gnu \
    && rm -rf /var/lib/apt/lists/*

# Download and install libnfc
latest_tag=$(curl -s https://api.github.com/repos/nfc-tools/libnfc/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
wget "https://github.com/nfc-tools/libnfc/releases/download/${latest_tag}/libnfc-${latest_tag#libnfc-}.tar.bz2"
tar -xvjf "libnfc-${latest_tag}.tar.bz2"
cd "libnfc-${latest_tag}"
make clean
make install all

# Configure libnfc
mkdir -p /etc/nfc /etc/nfc/devices.d

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

# Install WS281x drivers
cd ..
git clone https://github.com/jgarff/rpi_ws281x
cd rpi_ws281x && mkdir -p build && cd build

cmake -D BUILD_SHARED=OFF -D BUILD_TEST=ON ..
cmake --build .
make install

# Copy libraries and headers to appropriate directories
cp *.a /usr/local/lib
cp *.h /usr/local/include