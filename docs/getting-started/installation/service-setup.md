# Setting up ChargePi as a systemd service

**Golang should be installed and the binary should be added to PATH variable.**

1. Create a systemd unit file
   ```bash
   sudo nano /etc/systemd/system/ChargePi.service
   ```

2. Paste into ChargePi.service file:

   ```unit file
       [Unit]
       Description=ChargePi client 
       After=network.target
   
       [Service]
       Type=simple
       WorkingDirectory= /<path_to_dir>/ChargePi-go/
       ExecStart=go build main.go && ./main
       Restart=on-failure
       KillSignal=SIGTERM
   
       [Install]
       WantedBy=multi-user.target
   ```
3. Add the service to systemd:

   ```bash
   sudo chmod 640 /etc/systemd/system/ChargePi.service.service
   systemctl status ChargePi.service.service
   sudo systemctl daemon-reload
   sudo systemctl enable ChargePi.service.service
   sudo systemctl start ChargePi.service.service
   ```
