package grpc

import (
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/api"
	"google.golang.org/grpc"
	"net"
)

func CreateAndRunGrpcServer(address string, sendChannel chan api.Message, receiveChannel chan api.Message) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.WithError(err).Fatalf("Unable to listen to provided address: %s", address)
	}

	api.RegisterChargePointServer(grpcServer, api.NewApiServer(log.StandardLogger(), sendChannel, receiveChannel))

	err = grpcServer.Serve(listener)
	if err != nil {
		log.WithError(err).Fatal("Cannot expose API")
	}
}
