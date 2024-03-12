package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbuserarmies.UnimplementedV1UserArmyUnitsServer
}

func InitRoutes() pbuserarmies.V1UserArmyUnitsServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetUnit(ctx context.Context, req *pbuserarmies.GetUnitRequest) (res *pbuserarmies.GetUnitReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uau, err := gr.UserArmies().GetUnit(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pbuserarmies.GetUnitReply{
		Unit: uau.GetPbUserArmyType(),
	}, nil
}

func (h *grpcHandler) UpdateUnit(ctx context.Context, req *pbuserarmies.UpdateUnitRequest) (res *pbuserarmies.UpdateUnitReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	opts := types.UpdateUserArmyUnit{
		ID:     req.UpdateArmyUnit.GetId(),
		Points: utils.Ref(int(req.UpdateArmyUnit.GetPoints())),
	}

	if !req.UpdateArmyUnit.GetQtyNull() {
		opts.Qty = utils.Ref(int(req.UpdateArmyUnit.GetQtyValue()))
	}

	options := req.UpdateArmyUnit.GetOptionValues()
	for _, opt := range options {
		opts.OptionValues = append(opts.OptionValues, &types.UpdateUserArmyUnitOptionValue{
			ID:            opt.GetId(),
			IsSelected:    opt.GetIsSelected(),
			IndexSelected: opt.GetIndexSelected(),
			IDsSelected:   opt.GetIdsSelected(),
			QtySelected:   int(opt.GetQtySelected()),
		})
	}

	if err := gr.UserArmies().UpdateUnit(ctx, opts); err != nil {
		return nil, err
	}

	return new(pbuserarmies.UpdateUnitReply), nil
}
