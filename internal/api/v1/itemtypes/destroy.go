package itemtypes

import (
	"context"

	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteItemType(ctx context.Context, req *pbitemtypes.DeleteItemTypeRequest) (res *pbitemtypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
