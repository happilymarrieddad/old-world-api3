package unittypes

import (
	"context"

	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbunittypes.UnimplementedV1UnitTypesServer
}

func InitRoutes() pbunittypes.V1UnitTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetUnitType(ctx context.Context, req *pbunittypes.GetUnitTypeRequest) (res *pbunittypes.UnitType, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) DeleteUnitType(ctx context.Context, req *pbunittypes.DeleteUnitTypeRequest) (res *pbunittypes.EmptyReply, err error) {
	return nil, types.NewNotImplementedError()
}
