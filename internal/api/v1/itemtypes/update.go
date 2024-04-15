package itemtypes

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateItemType(ctx context.Context, req *pbitemtypes.UpdateItemTypeRequest) (res *pbitemtypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbitemtypes.EmptyReply{}, gr.ItemTypes().Update(ctx, types.UpdateItemType{
		ID:   req.GetId(),
		Name: req.GetName(),
	})
}
