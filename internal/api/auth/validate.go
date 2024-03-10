package auth

import (
	"context"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
)

func (s *grpcHandler) Validate(ctx context.Context, req *pbauth.JWTRequest) (reply *pbauth.ValidateReply, err error) {
	reply = new(pbauth.ValidateReply)
	userID, err := jwt.IsTokenValid(req.JWT)
	if err != nil {
		return nil, err
	}

	gr, err := interceptors.GetGlobalRepoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := gr.Users().GetByID(ctx, *userID)
	if err != nil {
		return nil, err
	}

	reply.User = &pbauth.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return reply, nil
}
