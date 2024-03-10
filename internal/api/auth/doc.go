package auth

import (
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
)

type grpcHandler struct {
	pbauth.UnimplementedAuthServer
}

func InitRoutes() pbauth.AuthServer {
	return &grpcHandler{}
}
