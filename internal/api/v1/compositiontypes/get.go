package compositiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *grpcHandler) GetCompositionType(ctx context.Context, req *pbcompositiontypes.GetCompositionTypeRequest) (res *pbcompositiontypes.CompositionType, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cts, _, err := gr.CompositionTypes().Find(ctx, &repos.FindCompositionTypesOpts{
		IDs: []string{req.GetId()}, Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(cts) == 0 {
		return nil, types.NewNotFoundError("unable to get army type by id")
	}

	return &pbcompositiontypes.CompositionType{
		Id:        cts[0].ID,
		GameId:    cts[0].GameID,
		Name:      cts[0].Name,
		CreatedAt: timestamppb.New(cts[0].CreatedAt.UTC()),
	}, nil
}
