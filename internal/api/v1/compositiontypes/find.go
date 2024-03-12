package compositiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositionTypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetCompositionTypes(ctx context.Context, req *pbcompositiontypes.GetCompositionTypesRequest) (res *pbcompositiontypes.GetCompositionTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, count, err := gr.CompositionTypes().Find(ctx, req.GameId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	res = new(pbcompositiontypes.GetCompositionTypesReply)
	for _, at := range ats {
		res.CompositionTypes = append(res.CompositionTypes, &pbcompositiontypes.CompositionType{
			Id:        at.ID,
			GameId:    at.GameID,
			Name:      at.Name,
			CreatedAt: timestamppb.New(at.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
