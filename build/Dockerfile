FROM --platform=$TARGETPLATFORM golang:1.23.0-bookworm AS base

WORKDIR /ChargePi/src
COPY . /ChargePi/src

ARG READER_TYPE
ARG TARGETOS
ARG TARGETARCH
ENV READER_TYPE=$READER_TYPE

ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV CGO_ENABLED=1
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on
ENV PKG_CONFIG_PATH=/usr/lib/pkgconfig:/usr/local/lib/pkgconfig
ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib"
ENV GOARM=7

# Install dependencies
RUN chmod +x /ChargePi/src/scripts/install-dependencies.sh
RUN /ChargePi/src/scripts/install-dependencies.sh $READER_TYPE

# Compile the client
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN go mod verify
RUN --mount=type=cache,target="/root/.cache/go-build" \
     go build -o ChargePi -v -tags="linux" .

FROM --platform=$BUILDPLATFORM node:19 AS build-ui
WORKDIR /ui

COPY ./ui/package*.json .
COPY ./ui .

RUN --mount=type=cache,target=/root/.npm npm install
RUN npm run build

# Test the client
FROM --platform=$TARGETPLATFORM golang:1.23.0-bookworm AS test

WORKDIR /ChargePi/src

ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV CGO_ENABLED=1
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on
ENV PKG_CONFIG_PATH=/usr/lib/pkgconfig:/usr/local/lib/pkgconfig
ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib"
ENV GOARM=7

# Copy installed C libraries and headers
COPY --from=base /usr/local/bin/ /usr/local/bin/
COPY --from=base /usr/local/lib/ /usr/local/lib/
COPY --from=base /usr/local/include/ /usr/local/include/
COPY --from=base /usr/include/ /usr/include/
COPY --from=base /usr/lib/ /usr/lib/

COPY --from=base /ChargePi/src /ChargePi/src

RUN go install github.com/vektra/mockery/v3@v3.5.1

# Gen mocks
RUN make gen-mocks

# Run tests
ENTRYPOINT ["go", "test","-v","./...","-coverpkg=./...","-short" ,"-coverprofile=unit_coverage.out"]

FROM --platform=$TARGETPLATFORM alpine AS chargepi

WORKDIR /etc/ChargePi

# Copy installed C libraries and headers
COPY --from=base /usr/local/bin/ /usr/local/bin/
COPY --from=base /usr/local/lib/ /usr/local/lib/
COPY --from=base /usr/local/include/ /usr/local/include/
COPY --from=base /usr/include/ /usr/include/
COPY --from=base /usr/lib/ /usr/lib/

# Copy the executable
#COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/ChargePi /usr/bin/ChargePi

# Copy the UI static files
COPY --from=build-ui /ui/build/* /usr/bin/ui/build/*

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s CMD curl --fail http://localhost:8081/healthz || exit 1

ENTRYPOINT ["ChargePi"]
CMD ["-c", "/etc/ChargePi/configs/config.yaml", "run"]