package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
)

func (h *grpcHandler) GetBasicUnitTypesByArmyTypeID(
	ctx context.Context, req *pbunittypes.ArmyTypeIdRequest,
) (res *pbunittypes.ArmyTypeIdReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uts, err := gr.UnitTypes().GetNamesByArmyTypeID(ctx, req.ArmyTypeId)
	if err != nil {
		return nil, err
	}

	res = new(pbunittypes.ArmyTypeIdReply)
	for _, ut := range uts {
		res.UnitTypes = append(res.UnitTypes, ut.GetPbUnitType())
	}

	return res, nil
}
