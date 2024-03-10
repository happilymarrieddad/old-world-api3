package auth

import (
	"context"
	"fmt"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (s *grpcHandler) Login(ctx context.Context, req *pbauth.LoginRequest) (reply *pbauth.LoginReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, err := gr.Users().GetByEmail(ctx, req.Email)
	if err != nil {
		fmt.Println("login GetByEmail failed with err: ", err.Error())
		return nil, err
	}

	if !usr.PasswordMatches(req.Password) {
		fmt.Println("login password does not match")
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
