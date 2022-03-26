package api

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
)

type (
	GrpcServer struct {
		UnimplementedChargePointServer
		sendChannel    chan Message
		receiveChannel chan Message
		logger         *log.Logger
	}

	Message struct {
		MessageId string
		Type      string
		Data      interface{}
	}
)

func NewApiServer(logger *log.Logger, sendChannel chan Message, receiveChannel chan Message) *GrpcServer {
	return &GrpcServer{
		logger:         logger,
		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
	}
}

func (s *GrpcServer) GetConnectorStatus(statusServer ChargePoint_GetConnectorStatusServer) error {
	for {
		_, err := statusServer.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (s *GrpcServer) StartTransaction(ctx context.Context, request *StartTransactionRequest) (*StartTransactionResponse, error) {
	response := &StartTransactionResponse{}

	return response, nil
}

func (s *GrpcServer) StopTransaction(ctx context.Context, request *StopTransactionRequest) (*StopTransactionResponse, error) {
	response := &StopTransactionResponse{}
	return response, nil
}

func (s *GrpcServer) HandleCharging(ctx context.Context, request *HandleChargingRequest) (*HandleChargingResponse, error) {
	response := &HandleChargingResponse{}
	return response, nil
}

func (s *GrpcServer) mustEmbedUnimplementedChargePointServer() {
}
