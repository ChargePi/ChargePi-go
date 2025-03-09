FROM --platform=$BUILDPLATFORM golang:1.23.0 AS base

WORKDIR /ChargePi/src
COPY . /ChargePi/src

ARG READER_TYPE
ENV READER_TYPE=$READER_TYPE

ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on

# Install dependencies
RUN chmod +x /ChargePi/src/scripts/install-dependencies.sh
RUN /ChargePi/src/scripts/install-dependencies.sh $READER_TYPE 0

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO=1

# Compile the client
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN go mod verify
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ChargePi .

FROM node:19 as build-ui

WORKDIR /ui

COPY ./ui/package*.json .
COPY ./ui .

RUN --mount=type=cache,target=/root/.npm npm install
RUN npm run build

# Test the client
FROM --platform=$BUILDPLATFORM vektra/mockery:v2.9.4 AS test

COPY --from=base /ChargePi/src /ChargePi/src
WORKDIR /ChargePi/src

RUN cp -r configs/ test/  \
    && cd test/  \
    && chmod +x create-test-certs.sh \
    && ./create-test-certs.sh

# Gen mocks
RUN mockery

# Run tests
CMD ["go", "test","-v","./..."]

FROM --platform=$BUILDPLATFORM alpine AS chargepi

WORKDIR /etc/ChargePi

ARG READER_TYPE
ENV READER_TYPE=$READER_TYPE

COPY --from=base /ChargePi/src/scripts /etc/ChargePi/scripts

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_TYPE 0

# Copy the executable
COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/ChargePi /usr/bin/ChargePi

# Copy the UI static files
COPY --from=build-ui /ui/build/ /usr/bin/ChargePi/ui/build

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s CMD curl --fail http://localhost:8081/healthz || exit 1

ENTRYPOINT ["ChargePi"]
CMD ["-c", "/etc/ChargePi/configs/config.yaml", "run"]