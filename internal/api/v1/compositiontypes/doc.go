package compositiontypes

import (
	"context"

	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositionTypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbcompositiontypes.UnimplementedV1CompositionTypesServer
}

func InitRoutes() pbcompositiontypes.V1CompositionTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetCompositionType(ctx context.Context, req *pbcompositiontypes.GetCompositionTypeRequest) (res *pbcompositiontypes.CompositionType, err error) {
	return nil, types.NewNotImplementedError()
}

func (h *grpcHandler) UpdateCompositionType(ctx context.Context, req *pbcompositiontypes.UpdateCompositionTypeRequest) (res *pbcompositiontypes.CompositionType, err error) {
	return nil, types.NewNotImplementedError()
}
