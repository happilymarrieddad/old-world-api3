package statistics

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) UpdateStatistic(ctx context.Context, req *pbstatistics.UpdateStatisticRequest) (res *pbstatistics.Statistic, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := gr.Statistics().Update(ctx, types.UpdateStatistic{
		ID:      req.GetId(),
		Name:    goaway.Censor(req.GetName()),
		Display: req.GetDisplay(),
	}); err != nil {
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
