package compositiontypes

import (
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositiontypes"
)

type grpcHandler struct {
	pbcompositiontypes.UnimplementedV1CompositionTypesServer
}

func InitRoutes() pbcompositiontypes.V1CompositionTypesServer {
	return &grpcHandler{}
}
