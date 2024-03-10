package trooptypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateTroopType(ctx context.Context, req *pbtrooptypes.CreateTroopTypeRequest) (res *pbtrooptypes.TroopType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newAt, err := gr.TroopTypes().FindOrCreate(ctx, types.CreateTroopType{
		Name:   req.GetName(),
		GameID: req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbtrooptypes.TroopType)
	res.Id = newAt.ID
	res.Name = newAt.Name
	res.GameId = newAt.GameID
	res.CreatedAt = timestamppb.New(newAt.CreatedAt.UTC())

	return res, nil
}
