package compositiontypes

import (
	"context"

	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositionTypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) DeleteCompositionType(ctx context.Context, req *pbcompositiontypes.DeleteCompositionTypeRequest) (res *pbcompositiontypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
