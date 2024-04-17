package games

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetGame(ctx context.Context, req *pbgames.GetGameRequest) (res *pbgames.GetGameReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	gm, err := gr.Games().Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pbgames.GetGameReply{
		Game: &pbgames.Game{
			Id:        gm.ID,
			Name:      gm.Name,
			CreatedAt: timestamppb.New(gm.CreatedAt.UTC()),
		},
	}, nil
}
