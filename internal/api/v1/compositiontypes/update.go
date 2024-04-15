package compositiontypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateCompositionType(ctx context.Context, req *pbcompositiontypes.UpdateCompositionTypeRequest) (res *pbcompositiontypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbcompositiontypes.EmptyReply{}, gr.CompositionTypes().Update(ctx, types.UpdateCompositionType{
		ID:   req.GetId(),
		Name: req.GetName(),
	})
}
