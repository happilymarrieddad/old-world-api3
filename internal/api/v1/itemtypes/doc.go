package itemtypes

import (
	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
)

type grpcHandler struct {
	pbitemtypes.UnimplementedV1ItemTypesServer
}

func InitRoutes() pbitemtypes.V1ItemTypesServer {
	return &grpcHandler{}
}
