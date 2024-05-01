package itemtypes

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateItemType(ctx context.Context, req *pbitemtypes.CreateItemTypeRequest) (res *pbitemtypes.ItemType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newAt, err := gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{
		Name:   goaway.Censor(req.GetName()),
		GameID: req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbitemtypes.ItemType)
	res.Id = newAt.ID
	res.Name = newAt.Name
	res.GameId = newAt.GameID
	res.CreatedAt = timestamppb.New(newAt.CreatedAt.UTC())

	return res, nil
}
