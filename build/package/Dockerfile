FROM  golang:1.16 as base
RUN apt-get update \
    apt-get upgrade \
    apt-get install git
    apt-get install cmake \

# install dependencies
COPY ../../docs/install-dependencies.sh ~/
RUN chmod +x /install-dependencies.sh
RUN ./install-dependencies.sh 0

COPY ../.. /ChargePi/src
WORKDIR /ChargePi/src

# Compile the client
RUN cd ~
RUN go build -o chargepi .

FROM base as dev
ENTRYPOINT ["go","run","."]

FROM alpine as chargepi

COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/chargepi /etc/ChargePi/chargepi
RUN cd /etc/ChargePi

ENTRYPOINT ["./chargepi"]

# Test the client
FROM base as test
RUN cp -r configs/ test/  \
    && cd test/  \
    && chmod +x create-test-certs.sh
    && ./create-test-certs.sh \
CMD ["go", "test","-v"]