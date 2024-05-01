package unittypes

import (
	"context"
	"log"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/utils"
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
		Name:              goaway.Censor(req.GetName()),
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

	go updateExistingUnits(ctx, updateUnitType.ID)

	return &pbunittypes.EmptyReply{}, nil
}

func updateExistingUnits(ctx context.Context, unitTypeID string) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		log.Printf("unable to GetGlobalRepoFromContext in UpdateUnitType.updateExistingUnits with err: %s\n", err.Error())
		return
	}

	existingUserArmyUnits, err := gr.UserArmies().FindUserUnitIDsByUnitTypeID(ctx, unitTypeID)
	if err != nil {
		log.Printf("unable to GetGlobalRepoFromContext in UpdateUnitType.FindUserUnitIDsByUnitTypeID with err: %s\n", err.Error())
		return
	}

	for _, uau := range existingUserArmyUnits {
		OptionValues := []*types.UpdateUserArmyUnitOptionValue{}

		for _, uauo := range uau.OptionValues {
			OptionValues = append(OptionValues, &types.UpdateUserArmyUnitOptionValue{
				ID:            uauo.ID,
				IsSelected:    uauo.IsSelected,
				IndexSelected: uauo.IndexSelected,
				IDsSelected:   uauo.IDsSelected,
				QtySelected:   uauo.QtySelected,
			})
		}

		if err = gr.UserArmies().UpdateUnit(ctx, types.UpdateUserArmyUnit{
			ID:           uau.ID,
			Qty:          utils.Ref(uau.Quantity),
			Points:       utils.Ref(uau.CalculatePoints()),
			OptionValues: OptionValues,
		}); err != nil {
			log.Printf("unable to update UserArmyUnit with err: %s\n", err.Error())
		}
	}
}
