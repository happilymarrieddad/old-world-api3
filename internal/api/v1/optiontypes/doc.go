package optiontypes

import (
	"context"

	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pboptiontypes.UnimplementedV1OptionTypesServer
}

func InitRoutes() pboptiontypes.V1OptionTypesServer {
	return &grpcHandler{}
}

func (h *grpcHandler) GetOptionType(ctx context.Context, req *pboptiontypes.GetOptionTypeRequest) (res *pboptiontypes.OptionType, err error) {
	return nil, types.NewNotImplementedError()
}
