package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
)

func (h *grpcHandler) GetUserArmy(ctx context.Context, req *pbuserarmies.GetUserArmyRequest) (res *pbuserarmies.GetUserArmyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ua, err := gr.UserArmies().Get(ctx, usr.ID, req.Id)
	if err != nil {
		return nil, err
	}

	return &pbuserarmies.GetUserArmyReply{
		UserArmy: ua.GetPbUserArmy(),
	}, nil
}
