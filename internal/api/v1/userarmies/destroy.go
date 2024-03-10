package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DestroyUserArmy(ctx context.Context, req *pbuserarmies.DestroyUserArmyRequest) (res *pbuserarmies.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uarepo := gr.UserArmies()
	ua, err := uarepo.Get(ctx, usr.ID, req.GetId())
	if err != nil {
		return nil, err
	} else if ua.UserID != usr.ID {
		return nil, types.NewUnauthorizedError("unauthorized")
	}

	// TODO: add to repo
	return &pbuserarmies.EmptyReply{}, types.NewNotImplementedError()
}
