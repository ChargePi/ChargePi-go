# 🏃 Running the ChargePi

This guide assumes you already have a working OCPP management/central system. If you do not have one already,
check out [SteVe](https://github.com/RWTH-i5-IDSG/steve) and set it up according to their guide.

## Standalone

Before you run/connect the client, make sure the backend is available and the charge point is
registered. Also, **libnfc** should be installed (a convenience script is available).

Running the client:

   ```bash
   go run -tags=rpi .
   ```

or compiling and executing the client:

   ```bash
   GOOS=linux
   GOARCH=arm64
   go build -o chargepi .
   ./chargepi
   ```

## 🐳 Deploying on Docker

1. Build the client image on Docker:

   ```bash
   cd ChargePi/client
   docker build -t chargepi .
   ```

   or pull it from Docker Hub:
   ```bash
      docker pull xblaz3kx/ChargePi-go:latest 
   ```

2. Run the container from built image:

   ```bash
   docker run --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --privileged chargepi
   ```

## Deploying using docker-compose

1. Build the ChargePi client:

   ```bash
   docker-compose -f ./deployments build .
   ```

2. Run services in daemon mode:

   ```bash
   docker-compose -f ./deployments up -d
   ```