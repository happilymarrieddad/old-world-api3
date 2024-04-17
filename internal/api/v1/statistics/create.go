package statistics

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateStatistic(ctx context.Context, req *pbstatistics.CreateStatisticRequest) (res *pbstatistics.Statistic, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newStat, err := gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{
		Name:    req.GetName(),
		Display: req.GetDisplay(),
		GameID:  req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbstatistics.Statistic)
	res.Id = newStat.ID
	res.Name = newStat.Name
	res.GameId = newStat.GameID
	res.Position = int64(newStat.Position)
	res.CreatedAt = timestamppb.New(newStat.CreatedAt.UTC())

	return res, nil
}
