package compositiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateCompositionType(ctx context.Context, req *pbcompositiontypes.CreateCompositionTypeRequest) (res *pbcompositiontypes.CompositionType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newAt, err := gr.CompositionTypes().FindOrCreate(ctx, types.CreateCompositionType{
		Name:   req.GetName(),
		GameID: req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbcompositiontypes.CompositionType)
	res.Id = newAt.ID
	res.Name = newAt.Name
	res.GameId = newAt.GameID
	res.CreatedAt = timestamppb.New(newAt.CreatedAt.UTC())

	return res, nil
}
