## Connecting to Mobile carrier via Huawei E3372 LTE dongle:

1. Update and upgrade:
   > sudo apt-get update & sudo apt-get upgrade
2. Install PPP:
   > sudo apt-get install ppp
3. Download Sakis3g client:
   > wget "http://raspberry-at-home.com/files/sakis3g.tar.gz"
4. Create a dir for the client:
   > sudo mkdir /usr/bin/modem3g
5. Give executable permissions:
   > sudo chmod 777 /usr/bin/modem3g
6. Copy to the created dir:
   > sudo cp sakis3g.tar.gz /usr/bin/modem3g
7. Move to the dir:
   > cd /usr/bin/modem3g
8. Extract Sakis3G client:
   > sudo tar -zxvf sakis3g.tar.gz
9. Add executable permissions
   > sudo chmod +x sakis3g
10. Run the client interactively:
    > sudo /usr/bin/modem3g/sakis3g connect --interactive
11. Get default settings for your modem:
    > lsusb
12. Paste default settings in a file:
    > sudo nano /etc/sakis3g.conf

    Config file example:

    ```
    USBDRIVER="option"
    USBINTERFACE="0"
    APN="internet"
    MODEM="12d1:155e"
    ```

## Running as service script:

1. Make two service files:

   > sudo nano /etc/systemd/system/modem-connection.service

   > sudo nano /etc/systemd/system/ChargePi.service

2. Paste into modem-connection.service file:

    ```
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

    ```
    [Unit]
    Description=ChargePi client 
    After=network.target modem-connection.service

    [Service]
    Type=simple
    WorkingDirectory=/<path_to_dir>/ChargePi-go 
    ExecStart=/usr/local/bin/go build main && ./main
    Restart=on-failure
    KillSignal=SIGTERM

    [Install]
    WantedBy=multi-user.target
    ```
**Golang should be installed and $GOPATH should be set to root user.**

Repeat next steps for both files:

4. Give permissions:
   > sudo chmod 640 /etc/systemd/system/modem-connection.service
5. Check the status of the service:
   > systemctl status modem-connection.service
6. Reload the daemon:
   > sudo systemctl daemon-reload
7. Enable service autostart:
   > sudo systemctl enable modem-connection.service
8. Start the service:
   > sudo systemctl start modem-connection.service

### References:

* [Sakis3g client](http://raspberry-at-home.com/installing-3g-modem/#more-138)
* [Systemd services](https://www.howtogeek.com/687970/how-to-run-a-linux-program-at-startup-with-systemd/)
* [Detailed Modem tutorial](https://lawrencematthew.wordpress.com/2013/08/07/connect-raspberry-pi-to-a-3g-network-automatically-during-its-boot/)