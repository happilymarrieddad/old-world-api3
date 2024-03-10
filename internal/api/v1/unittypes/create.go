package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateUnitType(ctx context.Context, req *pbunittypes.CreateUnitTypeRequest) (res *pbunittypes.UnitType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ut, err := gr.UnitTypes().FindOrCreate(ctx, types.CreateUnitType{
		Name:              req.GetName(),
		GameID:            req.GetGameId(),
		ArmyTypeID:        req.GetArmyTypeId(),
		TroopTypeID:       req.GetTroopTypeId(),
		CompositionTypeID: req.GetCompositionTypeId(),
		UnitTypeID:        req.GetUnitTypeId(),
		PointsPerModel:    int(req.GetPointsPerModel()),
		MinModels:         int(req.GetMinModels()),
		MaxModels:         int(req.GetMaxModels()),
	})
	if err != nil {
		return nil, err
	}

	return &pbunittypes.UnitType{
		Id:                  ut.ID,
		Name:                ut.Name,
		GameId:              ut.GameID,
		ArmyTypeId:          ut.ArmyTypeID,
		TroopTypeId:         ut.TroopTypeID,
		TroopTypeName:       ut.TroopTypeName,
		CompositionTypeId:   ut.CompositionTypeID,
		CompositionTypeName: ut.CompositionTypeName,
		PointsPerModel:      int64(ut.PointsPerModel),
		MinModels:           int64(ut.MinModels),
		MaxModels:           int64(ut.MaxModels),
		// These will get added later
		Statistics: []*pbunittypes.UnitStatistic{},
		Children:   []*pbunittypes.UnitTypeChild{},
		Options:    []*pbunittypes.UnitTypeOption{},
		CreatedAt:  timestamppb.New(ut.CreatedAt.UTC()),
	}, nil
}
