FROM  golang:1.17 as base

COPY . /ChargePi/src
WORKDIR /ChargePi/src

ARG READER_CONNECTION_TYPE
ENV READER_CONNECTION_TYPE=$READER_CONNECTION_TYPE

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_CONNECTION_TYPE 0

# Compile the client
RUN go mod download
RUN go mod verify
RUN go build -tags="rpi" -o ChargePi .

FROM base as dev
ENTRYPOINT ["go","run", "-tags=rpi","."]

FROM alpine as chargepi

WORKDIR /etc/ChargePi

ARG READER_CONNECTION_TYPE
ENV READER_CONNECTION_TYPE=$READER_CONNECTION_TYPE

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_CONNECTION_TYPE 0

COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/ChargePi /usr/bin/ChargePi

ENTRYPOINT ["./chargepi"]

# Test the client
FROM base as test
RUN cp -r configs/ test/  \
    && cd test/  \
    && chmod +x create-test-certs.sh \
    && ./create-test-certs.sh \
CMD ["go", "test","-v", "-tags=dev" ,"./..."]