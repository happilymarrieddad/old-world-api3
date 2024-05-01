package armytypes

import (
	"context"

	goaway "github.com/TwiN/go-away"
	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) UpdateArmyType(ctx context.Context, req *pbarmytypes.UpdateArmyTypeRequest) (res *pbarmytypes.EmptyReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pbarmytypes.EmptyReply{}, gr.ArmyTypes().Update(ctx, types.UpdateArmyType{
		ID:   req.GetId(),
		Name: goaway.Censor(req.GetName()),
	})
}
