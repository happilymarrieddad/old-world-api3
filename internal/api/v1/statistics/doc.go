package statistics

import (
	"context"

	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbstatistics.UnimplementedV1StatisticsServer
}

func InitRoutes() pbstatistics.V1StatisticsServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetStatistic(ctx context.Context, req *pbstatistics.GetStatisticRequest) (res *pbstatistics.Statistic, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) UpdateStatistic(ctx context.Context, req *pbstatistics.UpdateStatisticRequest) (res *pbstatistics.Statistic, err error) {
	return nil, types.NewNotImplementedError()
}
