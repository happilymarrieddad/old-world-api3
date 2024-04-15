package unittypes

import (
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
)

type grpcHandler struct {
	pbunittypes.UnimplementedV1UnitTypesServer
}

func InitRoutes() pbunittypes.V1UnitTypesServer {
	return &grpcHandler{}
}
