package grpc

import (
	"github.com/golang/protobuf/ptypes/empty"

	grpc "github.com/ChargePi/ChargePi-go/gen/proto/logs/v1"
)

type LogHandler struct {
	grpc.UnimplementedLogServiceServer
}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (s *LogHandler) GetLogs(e *empty.Empty, server grpc.LogService_GetLogsServer) error {
	// todo either a file-watcher or pipe directly from logrus (hook)?
	return nil
}

func (s *LogHandler) mustEmbedUnimplementedLogServer() {
}
