package auth

import (
	"context"
	"log"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
)

func (s *grpcHandler) Validate(ctx context.Context, req *pbauth.JWTRequest) (reply *pbauth.ValidateReply, err error) {
	userID, err := jwt.IsTokenValid(req.JWT)
	if err != nil {
		log.Printf("jwt.IsTokenValid err: %s\n", err.Error())
		return nil, err
	}

	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		log.Printf("interceptors.GetGlobalRepoFromContext err: %s\n", err.Error())
		return nil, err
	}

	user, err := gr.Users().GetByID(ctx, *userID)
	if err != nil {
		log.Printf("gr.Users().GetByID err: %s\n", err.Error())
		return nil, err
	}

	reply = new(pbauth.ValidateReply)
	reply.User = &pbauth.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return reply, nil
}
