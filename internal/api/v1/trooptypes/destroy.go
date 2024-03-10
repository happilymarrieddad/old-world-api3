package trooptypes

import (
	"context"

	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteTroopType(ctx context.Context, req *pbtrooptypes.DeleteTroopTypeRequest) (res *pbtrooptypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
