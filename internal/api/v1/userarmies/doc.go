package userarmies

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbuserarmies.UnimplementedV1UserArmiesServer
}

func InitRoutes() pbuserarmies.V1UserArmiesServer {
	return &grpcHandler{}
}

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

func (h *grpcHandler) AddUnit(ctx context.Context, req *pbuserarmies.AddUnitRequest) (res *pbuserarmies.EmptyReply, err error) {
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

	if _, err = uaRepo.AddUnits(ctx, req.UserArmyId, &types.CreateUserArmyUnit{
		UserArmyID: ua.ID,
		UnitTypeID: req.UnitTypeid,
	}); err != nil {
		return nil, err
	}

	return new(pbuserarmies.EmptyReply), nil
}

func (h *grpcHandler) RemoveUnit(ctx context.Context, req *pbuserarmies.RemoveUnitRequest) (res *pbuserarmies.EmptyReply, err error) {
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

	if err = uaRepo.RemoveUnits(ctx, ua.ID, req.Id); err != nil {
		return nil, err
	}

	return new(pbuserarmies.EmptyReply), nil
}
