package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateUnitType(ctx context.Context, req *pbunittypes.UpdateUnitTypeRequest) (res *pbunittypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbunittypes.EmptyReply{}, gr.UnitTypes().Update(ctx, types.UpdateUnitType{
		ID:             req.GetId(),
		Name:           req.GetName(),
		PointsPerModel: int(req.GetPointsPerModel()),
		MinModels:      int(req.GetMinModels()),
		MaxModels:      int(req.GetMaxModels()),
	})
}
