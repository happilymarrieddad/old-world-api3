package optiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetOptionTypes(ctx context.Context, req *pboptiontypes.GetOptionTypesRequest) (res *pboptiontypes.GetOptionTypesReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ats, count, err := gr.UnitOptionTypes().Find(ctx, req.GameId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	res = new(pboptiontypes.GetOptionTypesReply)
	for _, at := range ats {
		res.OptionTypes = append(res.OptionTypes, &pboptiontypes.OptionType{
			Id:        at.ID,
			GameId:    at.GameID,
			Name:      at.Name,
			CreatedAt: timestamppb.New(at.CreatedAt.UTC()),
		})
	}
	res.Count = count

	return res, nil
}
