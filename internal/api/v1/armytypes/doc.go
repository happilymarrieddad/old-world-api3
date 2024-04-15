package armytypes

import (
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
)

type grpcHandler struct {
	pbarmytypes.UnimplementedV1ArmyTypesServer
}

func InitRoutes() pbarmytypes.V1ArmyTypesServer {
	return &grpcHandler{}
}
