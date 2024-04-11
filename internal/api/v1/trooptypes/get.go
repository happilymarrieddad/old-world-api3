package trooptypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetTroopType(ctx context.Context, req *pbtrooptypes.GetTroopTypeRequest) (res *pbtrooptypes.TroopType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cts, _, err := gr.TroopTypes().Find(ctx, &repos.FindTroopTypeOpts{
		IDs: []string{req.GetId()}, Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(cts) == 0 {
		return nil, types.NewNotFoundError("unable to get army type by id")
	}

	return &pbtrooptypes.TroopType{
		Id:        cts[0].ID,
		GameId:    cts[0].GameID,
		Name:      cts[0].Name,
		CreatedAt: timestamppb.New(cts[0].CreatedAt.UTC()),
	}, nil
}
