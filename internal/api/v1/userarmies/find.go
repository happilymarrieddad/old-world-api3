package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
)

func (h *grpcHandler) GetUserArmies(ctx context.Context, req *pbuserarmies.GetUserArmiesRequest) (res *pbuserarmies.GetUserArmiesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uas, count, err := gr.UserArmies().Find(ctx, usr.ID, &repos.FindUserArmyOpts{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbuserarmies.GetUserArmiesReply)
	res.Count = count
	for _, ua := range uas {
		res.UserArmies = append(res.UserArmies, ua.GetPbUserArmy())
	}

	return res, nil
}
