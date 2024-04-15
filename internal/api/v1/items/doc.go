package items

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	pbitems "github.com/happilymarrieddad/old-world/api3/pb/proto/items"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbitems.UnimplementedV1ItemsServer
}

func InitRoutes() pbitems.V1ItemsServer {
	return &grpcHandler{}
}

func (h *grpcHandler) CreateItem(ctx context.Context, req *pbitems.CreateItemRequest) (res *pbitems.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if _, err = gr.Items().Create(ctx, types.NewItemFromPbCreateItem(req.GetItem())); err != nil {
		return nil, err
	}

	return &pbitems.EmptyReply{}, nil
}

func (h *grpcHandler) GetArmyItems(ctx context.Context, req *pbitems.GetArmyItemsRequest) (res *pbitems.GetItemsReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, _, err := gr.ArmyTypes().Find(ctx, &repos.FindArmyTypeOpts{
		IDs:   []string{req.GetArmyTypeId()},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(ats) == 0 {
		return nil, types.NewBadRequestError("invalid army type id passed in")
	}

	opts := &repos.FindItemsOpts{
		GameID:     ats[0].GameID,
		ArmyTypeID: utils.Ref(req.GetArmyTypeId()),
		Limit:      int(req.GetLimit()),
		Offset:     int(req.GetOffset()),
	}

	itms, err := gr.Items().Find(ctx, opts)
	if err != nil {
		return nil, err
	}

	res = new(pbitems.GetItemsReply)
	for _, itm := range itms {
		res.Items = append(res.Items, types.PbItemFromItem(itm))
	}

	return res, nil
}

func (h *grpcHandler) GetGameItems(ctx context.Context, req *pbitems.GetGameItemsRequest) (res *pbitems.GetItemsReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	opts := &repos.FindItemsOpts{
		GameID: req.GetGameId(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}

	itms, err := gr.Items().Find(ctx, opts)
	if err != nil {
		return nil, err
	}

	res = new(pbitems.GetItemsReply)
	for _, itm := range itms {
		res.Items = append(res.Items, types.PbItemFromItem(itm))
	}

	return res, nil
}

func (h *grpcHandler) GetItem(ctx context.Context, req *pbitems.GetItemRequest) (res *pbitems.GetItemReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	opts := &repos.FindItemsOpts{
		GameID: req.GetGameId(),
		IDs:    []string{req.GetId()}, Limit: 1,
	}

	itms, err := gr.Items().Find(ctx, opts)
	if err != nil {
		return nil, err
	} else if len(itms) == 0 {
		return nil, types.NewNotFoundError("item not found")
	}

	return &pbitems.GetItemReply{
		Item: types.PbItemFromItem(itms[0]),
	}, nil
}

func (h *grpcHandler) UpdateItem(ctx context.Context, req *pbitems.UpdateItemRequest) (res *pbitems.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = gr.Items().Update(ctx, types.NewItemFromPbUpdateItem(req.GetItem())); err != nil {
		return nil, err
	}

	return &pbitems.EmptyReply{}, nil
}
