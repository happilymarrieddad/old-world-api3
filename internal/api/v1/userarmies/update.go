package userarmies

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateUserArmy(ctx context.Context, req *pbuserarmies.UpdateUserArmyRequest) (res *pbuserarmies.GetUserArmyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	uaRepo := gr.UserArmies()
	ua, err := uaRepo.Get(ctx, usr.ID, req.UserArmyId)
	if err != nil {
		return nil, err
	} else if ua.UserID != usr.ID {
		return nil, types.NewUnauthorizedError("unauthorized")
	}

	updateReq := types.UpdateUserArmy{
		ID: req.GetUserArmyId(),
	}

	if len(req.Name) > 0 {
		req.Name = goaway.Censor(req.Name)
		updateReq.Name = &req.Name
	}

	if req.Points > 0 {
		updateReq.Points = &req.Points
	}

	if err := uaRepo.Update(ctx, updateReq); err != nil {
		return nil, err
	}

	ua, err = uaRepo.Get(ctx, usr.ID, req.UserArmyId)
	if err != nil {
		return nil, err
	}

	return &pbuserarmies.GetUserArmyReply{
		UserArmy: ua.GetPbUserArmy(),
	}, nil
}
