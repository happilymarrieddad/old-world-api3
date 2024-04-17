package games

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetGames(ctx context.Context, req *pbgames.GetGamesRequest) (res *pbgames.GetGamesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	games, err := gr.Games().Find(ctx, &repos.FindGameOpts{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbgames.GetGamesReply)
	for _, g := range games {
		res.Games = append(res.Games, &pbgames.Game{
			Id:        g.ID,
			Name:      g.Name,
			CreatedAt: timestamppb.New(g.CreatedAt.UTC()),
		})
	}

	return res, nil
}
