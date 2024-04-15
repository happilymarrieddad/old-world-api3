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

	updateUnitType := types.UpdateUnitType{
		ID:                req.GetId(),
		Name:              req.GetName(),
		PointsPerModel:    int(req.GetPointsPerModel()),
		MinModels:         int(req.GetMinModels()),
		MaxModels:         int(req.GetMaxModels()),
		TroopTypeID:       req.GetTroopTypeId(),
		CompositionTypeID: req.GetCompositionTypeId(),
	}

	for _, stat := range req.Statistics {
		updateUnitType.UpdateStatistics = append(updateUnitType.UpdateStatistics, &types.UpdateUnitStatistic{
			ID:    stat.Id,
			Value: stat.Value,
		})
	}

	if err := gr.UnitTypes().Update(ctx, updateUnitType); err != nil {
		return nil, err
	}

	// TODO: update any existing units that are using this unit type points

	return &pbunittypes.EmptyReply{}, nil
}
