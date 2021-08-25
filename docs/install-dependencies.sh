#!/bin/bash
# update and install updates
sudo apt-get update -y
sudo apt-get upgrade -y

# download libnfc
cd ~ && mkdir libnfc
cd libnfc/
wget https://github.com/nfc-tools/libnfc/releases/download/libnfc-1.8.0/libnfc-1.8.0.tar.bz2
tar -xvjf libnfc-1.8.0.tar.bz2
cd libnfc-1.8.0 && sudo mkdir /etc/nfc /etc/nfc/devices.d
# add configuration file
touch /etc/nfc/devices.d/pn532_i2c.conf
echo name = "PN532 board via I2C" >>/etc/nfc/devices.d/pn532_i2c.conf
echo connstring = pn532_i2c:/dev/i2c-1 >>/etc/nfc/devices.d/pn532_i2c.conf
echo allow_intrusive_scan = true >>/etc/nfc/devices.d/pn532_i2c.conf

#install dependencies
sudo apt-get install autoconf libtool libpcsclite-dev libusb-dev -y
autoreconf -vis
./configure --with-drivers=pn532_i2c --sysconfdir=/etc --prefix=/usr
sudo make clean
sudo make install all

# install WS281x drivers
sudo apt-get install cmake -y
cd ~
git clone https://github.com/jgarff/rpi_ws281x
cd rpi_ws281x && mkdir build
cd build
cmake -D BUILD_SHARED=OFF -D BUILD_TEST=ON ..
cmake --build .
sudo make install
cp *.a /usr/local/lib
cp *.h /usr/local/include

if [ $1 -eq 1 ]; then
  # install golang
  export GOLANG="$(curl https://golang.org/dl/ | grep armv6l | grep -v beta | head -1 | awk -F\> {'print $3'} | awk -F\< {'print $1'})"
  wget https://golang.org/dl/$GOLANG
  sudo tar -C /usr/local -xzf $GOLANG
  rm $GOLANG
  unset GOLANG
fi