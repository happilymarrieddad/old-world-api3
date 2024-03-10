package statistics

import (
	"context"

	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteStatistic(ctx context.Context, req *pbstatistics.DeleteStatisticRequest) (res *pbstatistics.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
