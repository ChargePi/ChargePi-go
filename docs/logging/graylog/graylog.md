# 📜 Graylog logging server

## Prerequisites

You should have a publicly accessible Linux or Windows server with Docker and docker-compose installed.

See [Docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/)
for installation guides.

## 🏃 Deploying Graylog server with docker-compose

1. Copy the following code into a docker-compose.yaml or use the provided [docker-compose](docker-compose.yaml). There
   are two docker-compose files - the first one is the base and the second one uses Traefik as a reverse proxy for
   routing.

   ```yaml
   version: '3'
   services:
     mongo:
       image: mongo:4.4.1
       networks:
         - graylog
       volumes:
         - mongo_data:/data/db
     elasticsearch:
       image: docker.elastic.co/elasticsearch/elasticsearch-oss:7.10.0
       environment:
         - http.host=0.0.0.0
         - transport.host=localhost
         - network.host=0.0.0.0
         - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
       networks:
         - graylog
       volumes:
         - es_data:/usr/share/elasticsearch/data
     graylog:
       image: graylog/graylog:4.0
       environment:
         - GRAYLOG_PASSWORD_SECRET=<password_sercret>
         # Password: admin
         - GRAYLOG_ROOT_PASSWORD_SHA2=571dc059e9aac28f272adbb668914318f95f2a2e98cae34d807f54710f6da8dd
         - GRAYLOG_HTTP_EXTERNAL_URI=http://<URI>:8081/
       volumes:
         - graylog_data:/usr/share/graylog/data
       networks:
         - graylog
       restart: always
       depends_on:
         - mongo
         - elasticsearch
       ports:
         - "8081:9000"
         - "1514:1514"
         - "1514:1514/udp"
         - "12201:12201"
         - "12201:12201/udp"
   networks:
     graylog:
       driver: bridge
   volumes:
     mongo_data:
       driver: local
     es_data:
       driver: local
     graylog_data:
       driver: local
   ```

2. Change the password hint and SHA256 as well as the external `GRAYLOG_HTTP_EXTERNAL_URI` environment variable.

3. Run using:

   ```bash
   sudo docker-compose up -d
   ```

   or

   ```bash
   sudo docker-compose -f docker-compose.yaml -f docker-compose.traefik.yaml up -d
   ```