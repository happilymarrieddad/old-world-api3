package games

import (
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
)

type grpcHandler struct {
	pbgames.UnimplementedV1GamesServer
}

func InitRoutes() pbgames.V1GamesServer {
	return &grpcHandler{}
}
