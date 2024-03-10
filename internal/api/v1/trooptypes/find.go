package trooptypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetTroopTypes(ctx context.Context, req *pbtrooptypes.GetTroopTypesRequest) (res *pbtrooptypes.GetTroopTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, count, err := gr.TroopTypes().Find(ctx, req.GameId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	res = new(pbtrooptypes.GetTroopTypesReply)
	for _, at := range ats {
		res.TroopTypes = append(res.TroopTypes, &pbtrooptypes.TroopType{
			Id:        at.ID,
			GameId:    at.GameID,
			Name:      at.Name,
			CreatedAt: timestamppb.New(at.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
