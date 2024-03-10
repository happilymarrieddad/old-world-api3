package optiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateOptionType(ctx context.Context, req *pboptiontypes.CreateOptionTypeRequest) (res *pboptiontypes.OptionType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newAt, err := gr.UnitOptionTypes().FindOrCreate(ctx, types.CreateUnitOptionType{
		Name:   req.GetName(),
		GameID: req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pboptiontypes.OptionType)
	res.Id = newAt.ID
	res.Name = newAt.Name
	res.GameId = newAt.GameID
	res.CreatedAt = timestamppb.New(newAt.CreatedAt.UTC())

	return res, nil
}
