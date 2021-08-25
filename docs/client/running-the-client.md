# Running the client

## Standalone

This client uses **[SteVe](https://github.com/RWTH-i5-IDSG/steve)** for the Central System, but can connect to other
Central Systems as well.

Running the program in Golang:

	go run main.go

Compiling and executing the program:

```
go build chargepi
./chargepi
```

## Deploying on Docker

1. Build the client image on Docker:

```
cd ChargePi/client
docker build -t chargepi .
```

2. Run the container from built image:

> docker run --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --privileged chargepi

## Deploying using docker-compose

Alternatively, you can run the client, SteVe server and Watchtower on the same Pi using **docker-compose**.
The **[Watchtower](https://github.com/containrrr/watchtower)** service will automatically pull the newest image and run
it when it is available.

1. Change the IP address under __serverUri__ in the settings file to **172.0.1.121**.

2. Build services:

   > docker-compose build

3. Run the client in daemon mode:

   > docker-compose up -d

## Deploying with other services with docker-compose 