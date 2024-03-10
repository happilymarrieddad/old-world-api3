package unittypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
)

func (h *grpcHandler) GetUnitTypes(ctx context.Context, req *pbunittypes.GetUnitTypesRequest) (res *pbunittypes.GetUnitTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uts, count, err := gr.UnitTypes().Find(ctx, &repos.FindUnitTypesOpts{
		ArmyTypeID:             req.GetArmyTypeId(),
		Limit:                  int(req.GetLimit()),
		Offset:                 int(req.GetOffset()),
		IncludeUnitTypeOptions: req.GetIncludeUnitTypeOptions(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbunittypes.GetUnitTypesReply)
	res.Count = count
	for _, ut := range uts {
		res.UnitTypes = append(res.UnitTypes, ut.GetPbUnitType())
	}

	return res, nil
}
