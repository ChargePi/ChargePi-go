# 🛠️ Configuring mobile connectivity

## Setting up & running Sakis3G

1. Update and upgrade:

   ```bash
   sudo apt-get update & sudo apt-get upgrade
   ```

2. Install the Sakis3g client:

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

4. Get default settings for your modem:

   ```bash
   lsusb
   ```

## Running as a service

1. Paste default settings in a file:

   ```bash
   sudo nano /etc/sakis3g.conf
   ```

   Config file example:

    ```text
    USBDRIVER="option"
    USBINTERFACE="0"
    APN="internet"
    MODEM="12d1:155e"
    ```

2. Make a systemd service file:

   ```bash
   sudo nano /etc/systemd/system/modem-connection.service
   ```

3. Paste into modem-connection.service file:

   ```unit file
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

4. Give permissions and add services to **systemd**:

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