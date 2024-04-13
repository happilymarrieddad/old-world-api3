package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
)

func (h *grpcHandler) GetUnitType(ctx context.Context, req *pbunittypes.GetUnitTypeRequest) (res *pbunittypes.UnitType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ut, err := gr.UnitTypes().Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return ut.GetPbUnitType(), nil
}
