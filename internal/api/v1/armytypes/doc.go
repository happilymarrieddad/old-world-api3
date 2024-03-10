package armytypes

import (
	"context"

	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbarmytypes.UnimplementedV1ArmyTypesServer
}

func InitRoutes() pbarmytypes.V1ArmyTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetArmyType(ctx context.Context, req *pbarmytypes.GetArmyTypeRequest) (res *pbarmytypes.ArmyType, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) UpdateArmyType(ctx context.Context, req *pbarmytypes.UpdateArmyTypeRequest) (res *pbarmytypes.ArmyType, err error) {
	return nil, types.NewNotImplementedError()
}
