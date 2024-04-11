package armytypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetArmyType(ctx context.Context, req *pbarmytypes.GetArmyTypeRequest) (res *pbarmytypes.ArmyType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, _, err := gr.ArmyTypes().Find(ctx, &repos.FindArmyTypeOpts{
		IDs: []string{req.GetId()}, Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(ats) == 0 {
		return nil, types.NewNotFoundError("unable to get army type by id")
	}

	return &pbarmytypes.ArmyType{
		Id:        ats[0].ID,
		GameId:    ats[0].GameID,
		Name:      ats[0].Name,
		CreatedAt: timestamppb.New(ats[0].CreatedAt.UTC()),
	}, nil
}
