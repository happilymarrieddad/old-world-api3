package statistics

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetStatistic(ctx context.Context, req *pbstatistics.GetStatisticRequest) (res *pbstatistics.Statistic, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	stat, err := gr.Statistics().Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pbstatistics.Statistic{
		Id:        stat.ID,
		Name:      stat.Name,
		Display:   stat.Display,
		Position:  int64(stat.Position),
		CreatedAt: timestamppb.New(stat.CreatedAt),
	}, nil
}
