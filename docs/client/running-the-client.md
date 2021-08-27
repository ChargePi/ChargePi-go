# Running the client

## Standalone

This client uses **[SteVe](https://github.com/RWTH-i5-IDSG/steve)** for the Central System, but can connect to other
Central Systems as well.

Running the program in Golang:

   ```bash
   go run main.go
   ```

Compiling and executing the program:

   ```bash
   go build main.go
   ./main
   ```

### Client flags/options

The client comes with flags for configurations' filepath, so you can specify where your configuration files are located.
All the paths must be absolute. If no flags are provided, the client will default to the working directory.

Example:

 ```bash
   go run main.go --help
      -auth string
           Path to the authorization persistence file (default "/home/pi/ChargePi-go/configs/auth.json")
     -config-folder string
           Path to the configuration folder (default "/home/pi/ChargePi-go/configs")
     -connector-folder string
           Path to the connector folder (default "/home/pi/ChargePi-go/configs/connectors")
     -ocpp-config string
           Path to the OCPP configuration file (default "/home/pi/ChargePi-go/configs/configuration.json")
     -settings string
           Path to the settings file (default "/home/pi/ChargePi-go/configs/settings.json")
   ```

## Deploying on Docker

1. Build the client image on Docker:

   ```bash
   cd ChargePi/client
   docker build -t chargepi .
   ```

2. Run the container from built image:

   ```bash
   docker run --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --privileged chargepi
   ```

## Deploying using docker-compose

Alternatively, you can run the client, SteVe server and Watchtower on the same Pi using **docker-compose**.
The **[Watchtower](https://github.com/containrrr/watchtower)** service will automatically pull the newest image and run
it when it is available.

1. Change the IP address under __serverUri__ in the settings file to **172.0.1.121**.

2. Build services:

   ```bash
   docker-compose build
   ```

3. Run services in daemon mode:

   ```bash
   docker-compose up -d
   ```