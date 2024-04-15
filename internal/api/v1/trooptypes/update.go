package trooptypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateTroopType(ctx context.Context, req *pbtrooptypes.UpdateTroopTypeRequest) (res *pbtrooptypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbtrooptypes.EmptyReply{}, gr.TroopTypes().Update(ctx, types.UpdateTroopType{
		ID:   req.GetId(),
		Name: req.GetName(),
	})
}
