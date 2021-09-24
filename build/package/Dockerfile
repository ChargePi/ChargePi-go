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
RUN cd ~
ENTRYPOINT ["go","run","."]

FROM base as chargepi
RUN cd ~
ENTRYPOINT ["./chargepi"]

# Test the client
FROM base as test
RUN cp -r configs/ test/ && cd test/ && chmod +x create-test-certs.sh
&& ./create-test-certs.sh
CMD ["go", "test","-v"]