package armytypes

import (
	"context"

	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteArmyType(ctx context.Context, req *pbarmytypes.DeleteArmyTypeRequest) (res *pbarmytypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
