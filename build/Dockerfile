FROM golang:1.18 as base

COPY . /ChargePi/src
WORKDIR /ChargePi/src

ARG READER_TYPE
ENV READER_TYPE=$READER_TYPE

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_TYPE 0

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO=1

# Compile the client
RUN go mod download
RUN go mod verify
RUN go build -o ChargePi .

FROM node:19 as build-ui

COPY ./ui /ui
WORKDIR /ui

RUN npm install
RUN npm run build

FROM alpine as chargepi

WORKDIR /etc/ChargePi

ARG READER_CONNECTION_TYPE
ENV READER_CONNECTION_TYPE=$READER_TYPE

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_TYPE 0

# Copy the xecutable
COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/ChargePi /usr/bin/ChargePi

# Copy the UI static files
COPY --from=build-ui /ui/build/ /usr/bin/ChargePi/ui/build

CMD ["./chargepi"]

# Test the client
FROM base as test

RUN cp -r configs/ test/  \
    && cd test/  \
    && chmod +x create-test-certs.sh \
    && ./create-test-certs.sh \

# Gen mocks

CMD ["go", "test","-v","./..."]