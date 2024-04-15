package optiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateOptionType(ctx context.Context, req *pboptiontypes.UpdateOptionTypeRequest) (res *pboptiontypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pboptiontypes.EmptyReply{}, gr.UnitOptionTypes().Update(ctx, types.UpdateUnitOptionType{
		ID:   req.GetId(),
		Name: req.GetName(),
	})
}
