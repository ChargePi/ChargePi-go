package service

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
)

type LogService struct {
	grpc.UnimplementedLogServer
}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) GetLogs(e *empty.Empty, server grpc.Log_GetLogsServer) error {
	return nil
}

func (s *LogService) mustEmbedUnimplementedLogServer() {
}
