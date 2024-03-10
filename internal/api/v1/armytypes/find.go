package armytypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetArmyTypes(ctx context.Context, req *pbarmytypes.GetArmyTypesRequest) (res *pbarmytypes.GetArmyTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, count, err := gr.ArmyTypes().Find(ctx, req.GameId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	res = new(pbarmytypes.GetArmyTypesReply)
	for _, at := range ats {
		res.ArmyTypes = append(res.ArmyTypes, &pbarmytypes.ArmyType{
			Id:        at.ID,
			GameId:    at.GameID,
			Name:      at.Name,
			CreatedAt: timestamppb.New(at.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
