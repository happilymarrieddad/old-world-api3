package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) CreateUserArmy(ctx context.Context, req *pbuserarmies.CreateUserArmyRequest) (res *pbuserarmies.CreateUserArmyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ua, err := gr.UserArmies().Create(ctx, types.CreateUserArmy{
		Name:       req.GetName(),
		UserID:     usr.ID,
		GameID:     req.GetGameId(),
		ArmyTypeID: req.GetArmyTypeId(),
		Points:     int(req.GetPoints()),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbuserarmies.CreateUserArmyReply)
	res.UserArmy = ua.GetPbUserArmy()

	return res, nil
}
