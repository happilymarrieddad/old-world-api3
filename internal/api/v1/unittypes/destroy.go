package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
)

func (h *grpcHandler) DeleteUnitType(ctx context.Context, req *pbunittypes.DeleteUnitTypeRequest) (res *pbunittypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbunittypes.EmptyReply{}, gr.UnitTypes().Destroy(ctx, req.GetId())
}
