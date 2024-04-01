package grpc

import (
	"github.com/ChargePi/ChargePi-go/pkg/grpc"
	"github.com/golang/protobuf/ptypes/empty"
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
