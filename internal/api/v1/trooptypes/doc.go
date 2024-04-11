package trooptypes

import (
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
)

type grpcHandler struct {
	pbtrooptypes.UnimplementedV1TroopTypesServer
}

func InitRoutes() pbtrooptypes.V1TroopTypesServer {
	return &grpcHandler{}
}
