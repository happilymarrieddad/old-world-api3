package statistics

import (
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
)

type grpcHandler struct {
	pbstatistics.UnimplementedV1StatisticsServer
}

func InitRoutes() pbstatistics.V1StatisticsServer {
	return &grpcHandler{}
}
