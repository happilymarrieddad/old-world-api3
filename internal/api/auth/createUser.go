package auth

import (
	"context"
	"errors"
	"log"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func (h *grpcHandler) CreateUser(ctx context.Context, req *pbauth.CreateUserRequest) (reply *pbauth.LoginReply, err error) {
	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	usr, _ := gr.Users().GetByEmail(ctx, req.Email)
	if usr != nil {
		return nil, errors.New("user already exists")
	}

	if len(req.Password) == 0 || len(req.Email) == 0 || len(req.FirstName) == 0 || len(req.LastName) == 0 {
		return nil, errors.New("must include all necessary fields")
	}

	if req.Password != req.PasswordConfirm {
		return nil, errors.New("password must match password confirm")
	}

	usr, err = gr.Users().Create(ctx, types.CreateUser{
		FirstName:       req.GetFirstName(),
		LastName:        req.GetLastName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	})
	if err != nil {
		log.Println("unable to create user with err: " + err.Error())
		return nil, err
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
