package armytypes

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) CreateArmyType(ctx context.Context, req *pbarmytypes.CreateArmyTypeRequest) (res *pbarmytypes.ArmyType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newAt, err := gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{
		Name:   goaway.Censor(req.GetName()),
		GameID: req.GetGameId(),
	})
	if err != nil {
		return nil, err
	}

	res = new(pbarmytypes.ArmyType)
	res.Id = newAt.ID
	res.Name = newAt.Name
	res.GameId = newAt.GameID
	res.CreatedAt = timestamppb.New(newAt.CreatedAt.UTC())

	return res, nil
}
