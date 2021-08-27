# Configuring mobile connectivity

## Setting up & running Sakis3G

1. Update and upgrade:

   ```bash
   sudo apt-get update & sudo apt-get upgrade
   ```

2. Installing dependencies and- Sakis3g client:

   ```bash
   sudo apt-get install ppp
   wget "http://raspberry-at-home.com/files/sakis3g.tar.gz"
   sudo mkdir /usr/bin/modem3g
   sudo chmod +x /usr/bin/modem3g
   sudo cp sakis3g.tar.gz /usr/bin/modem3g
   cd /usr/bin/modem3g
   ```

3. Run the client interactively:

   ```bash
   sudo /usr/bin/modem3g/sakis3g connect --interactive
   ```

5. Get default settings for your modem:

   ```bash
   lsusb
   ```

6. Paste default settings in a file:

   ```bash
   sudo nano /etc/sakis3g.conf
   ```

   Config file example:

    ```
    USBDRIVER="option"
    USBINTERFACE="0"
    APN="internet"
    MODEM="12d1:155e"
    ```

## Running as service script:

**Golang should be installed and the binary should be added to PATH variable.**

1. Make two service files:

   ```bash
   sudo nano /etc/systemd/system/modem-connection.service
   sudo nano /etc/systemd/system/ChargePi.service
   ```

2. Paste into modem-connection.service file:

   ```bash
       [Unit]
       Description=Modem connection service
   
       [Service]
       Type=simple 
       ExecStart=/usr/bin/modem3g/sakis3g --sudo connect 
       Restart=on-failure 
       RestartSec=5  
       KillMode=process
   
       [Install]
       WantedBy=multi-user.target
   ```

3. Paste into ChargePi.service file:

   ```bash
       [Unit]
       Description=ChargePi client 
       After=network.target modem-connection.service
   
       [Service]
       Type=simple
       WorkingDirectory= /<path_to_dir>/ChargePi-go/
       ExecStart=go build main.go && ./main
       Restart=on-failure
       KillSignal=SIGTERM
   
       [Install]
       WantedBy=multi-user.target
   ```

4. Give permissions and add services to **systemd** (repeat for both service files):

   ```bash
   sudo chmod 640 /etc/systemd/system/modem-connection.service
   systemctl status modem-connection.service
   sudo systemctl daemon-reload
   sudo systemctl enable modem-connection.service
   sudo systemctl start modem-connection.service
   ```

### References:

* [Sakis3g client](http://raspberry-at-home.com/installing-3g-modem/#more-138)
* [Systemd services](https://www.howtogeek.com/687970/how-to-run-a-linux-program-at-startup-with-systemd/)
* [Detailed Modem tutorial](https://lawrencematthew.wordpress.com/2013/08/07/connect-raspberry-pi-to-a-3g-network-automatically-during-its-boot/)