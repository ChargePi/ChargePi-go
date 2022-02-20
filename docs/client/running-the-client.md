# 🏃 Running the client

## Standalone

This client uses **[SteVe](https://github.com/RWTH-i5-IDSG/steve)** for the Central System, but can connect to other
Central Systems as well. Before you run/connect the client, make sure the backend is available and the charge point is
registered.

1. Running the program in Golang:

   ```bash
   go run .
   ```

2. Compiling and executing the program:

   ```bash
   go build -o chargepi .
   ./chargepi
   ```

### 🚩 Client flags/options

|        Flag         | Short |           Description           | Default value |
|:-------------------:|:-----:|:-------------------------------:|:-------------:|
|     `-settings`     |   /   |   Path to the settings file.    |               |
| `-connector-folder` |   /   |  Path to the connector folder.  |               |
|   `-ocpp-config`    |   /   | Path to the OCPP configuration. |               |
|       `-auth`       |   /   | Path to the authorization file. |               |
|      `-debug`       | `--d` |           Debug mode            |     false     |
|       `-api`        | `--a` |         Expose the API          |     false     |
|   `-api-address`    |   /   |           API address           |  "localhost"  |
|     `-api-port`     |   /   |            API port             |     4269      |

## 🐳 Deploying on Docker

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