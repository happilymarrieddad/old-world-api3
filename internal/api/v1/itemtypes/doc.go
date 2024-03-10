package itemtypes

import (
	"context"

	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbitemtypes.UnimplementedV1ItemTypesServer
}

func InitRoutes() pbitemtypes.V1ItemTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetItemType(ctx context.Context, req *pbitemtypes.GetItemTypeRequest) (res *pbitemtypes.ItemType, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) UpdateItemType(ctx context.Context, req *pbitemtypes.UpdateItemTypeRequest) (res *pbitemtypes.ItemType, err error) {
	return nil, types.NewNotImplementedError()
}
