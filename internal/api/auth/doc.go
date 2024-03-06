package auth

import (
	"context"
	"fmt"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type grpcHandler struct {
	pbauth.UnimplementedV1AuthServer
}

func InitRoutes() pbauth.V1AuthServer {
	return &grpcHandler{}
}

func (s *grpcHandler) Login(ctx context.Context, req *pbauth.LoginRequest) (reply *pbauth.LoginReply, err error) {
	fmt.Println("Inside login")
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := gr.Users().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !usr.PasswordMatches(req.Password) {
		return nil, types.NewUnauthorizedError("unauthorized")
	}

	reply = new(pbauth.LoginReply)
	reply.Bearer = jwt.NewToken(usr)

	return reply, nil
}

func (s *grpcHandler) Refresh(ctx context.Context, req *pbauth.JWTRequest) (reply *pbauth.LoginReply, err error) {
	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	reply = new(pbauth.LoginReply)
	reply.Bearer = jwt.NewToken(usr)

	return reply, nil
}

func (s *grpcHandler) Validate(ctx context.Context, req *pbauth.JWTRequest) (reply *pbauth.EmptyReply, err error) {
	return
}
