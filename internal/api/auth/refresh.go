package auth

import (
	"context"
	"log"

	"github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
)

func (s *grpcHandler) Refresh(ctx context.Context, req *pbauth.JWTRequest) (reply *pbauth.LoginReply, err error) {
	usr, err := interceptors.GetUserFromContext(ctx)
	if err != nil {
		log.Printf("interceptors.GetUserFromContext err: %s\n", err.Error())
		return nil, err
	}

	reply = new(pbauth.LoginReply)
	reply.Bearer = jwt.NewToken(usr)

	return reply, nil
}
