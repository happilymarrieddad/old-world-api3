package games

import (
	"context"
	"fmt"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
)

func (h *grpcHandler) CreateGame(ctx context.Context, req *pbgames.CreateGameRequest) (res *pbgames.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if _, err = gr.Games().FindOrCreate(ctx, goaway.Censor(req.Name)); err != nil {
		fmt.Println("Game failed to create err: " + err.Error())
		return nil, err
	}

	return new(pbgames.EmptyReply), nil
}
