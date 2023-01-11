package service

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
)

type LogService struct {
	grpc.UnimplementedLogServer
}

func NewLogService() *LogService {
	return &LogService{}
}

func (s *LogService) GetLogs(e *empty.Empty, server grpc.Log_GetLogsServer) error {
	// todo either a file-watcher or pipe directly from logrus (hook)?
	return nil
}

func (s *LogService) mustEmbedUnimplementedLogServer() {
}
