package auth

import (
	"context"
	"log"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (s *grpcHandler) Login(ctx context.Context, req *pbauth.LoginRequest) (reply *pbauth.LoginReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		log.Println("unable to get global repo from context")
		return nil, err
	}

	usr, err := gr.Users().GetByEmail(ctx, req.Email)
	if err != nil {
		log.Println("login GetByEmail failed with err: ", err.Error())
		return nil, err
	}

	if !usr.PasswordMatches(req.Password) {
		log.Println("login password does not match")
		return nil, types.NewUnauthorizedError("unauthorized")
	}

	reply = new(pbauth.LoginReply)
	reply.Bearer = jwt.NewToken(usr)
	reply.User = &pbauth.User{
		Id:        usr.ID,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
	}

	return reply, nil
}
