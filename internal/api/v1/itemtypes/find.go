package itemtypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetItemTypes(ctx context.Context, req *pbitemtypes.GetItemTypesRequest) (res *pbitemtypes.GetItemTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, count, err := gr.ItemTypes().Find(ctx, &repos.FindItemTypeOpts{
		GameID: req.GetGameId(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbitemtypes.GetItemTypesReply)
	for _, at := range ats {
		res.ItemTypes = append(res.ItemTypes, &pbitemtypes.ItemType{
			Id:        at.ID,
			GameId:    at.GameID,
			Name:      at.Name,
			CreatedAt: timestamppb.New(at.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
