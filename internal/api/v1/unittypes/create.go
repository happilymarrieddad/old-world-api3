package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateUnitType(ctx context.Context, req *pbunittypes.CreateUnitTypeRequest) (res *pbunittypes.UnitType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	parentUnitTypeID := req.GetUnitTypeId()

	newUnitType := types.CreateUnitType{
		Name:              req.GetName(),
		GameID:            req.GetGameId(),
		ArmyTypeID:        req.GetArmyTypeId(),
		TroopTypeID:       req.GetTroopTypeId(),
		CompositionTypeID: req.GetCompositionTypeId(),
		UnitTypeID:        parentUnitTypeID,
		PointsPerModel:    int(req.GetPointsPerModel()),
		MinModels:         int(req.GetMinModels()),
		MaxModels:         int(req.GetMaxModels()),
	}

	// TODO: add statistics to the CreateUnitTypeRequest object so it can be added instead
	if len(parentUnitTypeID) > 0 {
		put, err := gr.UnitTypes().Get(ctx, parentUnitTypeID)
		if err != nil {
			return nil, err
		}

		for _, stat := range put.Statistics {
			newUnitType.Statistics = append(newUnitType.Statistics, &types.CreateUnitStatistics{
				Display: stat.Statistic.Display,
				Value:   stat.Value,
			})
		}
	} else {
		stats, _, err := gr.Statistics().Find(ctx, &repos.FindStatisticsOpts{
			GameIDs: []string{req.GetGameId()},
			Limit:   10000,
		})
		if err != nil {
			return nil, err
		}

		for _, stat := range stats {
			newUnitType.Statistics = append(newUnitType.Statistics, &types.CreateUnitStatistics{
				Display: stat.Display,
				Value:   "-",
			})
		}
	}

	ut, err := gr.UnitTypes().FindOrCreate(ctx, newUnitType)
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
