package games

import (
	"context"

	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteGame(context.Context, *pbgames.DeleteGameRequest) (*pbgames.EmptyReply, error) {
	return nil, types.NewNotImplementedError()
}
