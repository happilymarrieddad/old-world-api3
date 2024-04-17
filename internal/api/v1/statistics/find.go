package statistics

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetStatistics(ctx context.Context, req *pbstatistics.GetStatisticsRequest) (res *pbstatistics.GetStatisticsReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	stats, count, err := gr.Statistics().Find(ctx, &repos.FindStatisticsOpts{
		GameIDs: []string{req.GetGameId()},
		Limit:   int(req.Limit),
		Offset:  int(req.Offset),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbstatistics.GetStatisticsReply)
	for _, stat := range stats {
		res.Statistics = append(res.Statistics, &pbstatistics.Statistic{
			Id:        stat.ID,
			GameId:    stat.GameID,
			Name:      stat.Name,
			Display:   stat.Display,
			Position:  int64(stat.Position),
			CreatedAt: timestamppb.New(stat.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
