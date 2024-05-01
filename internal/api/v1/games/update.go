package games

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateGame(ctx context.Context, req *pbgames.UpdateGameRequest) (res *pbgames.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbgames.EmptyReply{}, gr.Games().Update(ctx, types.UpdateGame{
		ID:   req.GetId(),
		Name: goaway.Censor(req.GetName()),
	})
}
