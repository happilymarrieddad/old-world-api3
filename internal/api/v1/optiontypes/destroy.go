package optiontypes

import (
	"context"

	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteOptionType(ctx context.Context, req *pboptiontypes.DeleteOptionTypeRequest) (res *pboptiontypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
