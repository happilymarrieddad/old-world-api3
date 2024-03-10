package trooptypes

import (
	"context"

	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbtrooptypes.UnimplementedV1TroopTypesServer
}

func InitRoutes() pbtrooptypes.V1TroopTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetTroopType(ctx context.Context, req *pbtrooptypes.GetTroopTypeRequest) (res *pbtrooptypes.TroopType, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) UpdateTroopType(ctx context.Context, req *pbtrooptypes.UpdateTroopTypeRequest) (res *pbtrooptypes.TroopType, err error) {
	return nil, types.NewNotImplementedError()
}
