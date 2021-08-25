FROM  golang:1.16 as base
RUN apt-get update \
    apt-get upgrade \
    apt-get install git
    apt-get install cmake \
COPY ../../docs/install-dependencies.sh ~/
RUN chmod +x /install-dependencies.sh
RUN ./install-dependencies.sh 0
RUN mkdir /ChargePi
COPY ../.. /ChargePi
WORKDIR /ChargePi
RUN cd ~

#Run the client
FROM base as run
# Build the client
RUN go build -o main .
CMD ["./ChargePi/main"]

# Test the client
FROM base as test
RUN cp -r configs/ test/ && cd test/ && chmod +x create-test-certs.sh
&& ./create-test-certs.sh
CMD ["go", "test -v"]