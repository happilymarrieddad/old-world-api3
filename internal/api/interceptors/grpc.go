package interceptors

import (
	"context"
	"fmt"
	"reflect"

	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"google.golang.org/grpc"
)

func GlobalRepoInjector(gr repos.GlobalRepo) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newContext := SetGlobalRepoOnContext(ctx, gr)

		// BEFORE the request
		v := reflect.Indirect(reflect.ValueOf(req))
		vField := reflect.Indirect(v.FieldByName("JWT"))

		fmt.Printf("Method: %s\n", info.FullMethod)

		// If there is NO jwt field on the request, then just continue
		if !vField.IsValid() {
			return handler(newContext, req)
		}

		jwtToken := vField.String()

		userID, err := jwt.IsTokenValid(jwtToken)
		if err != nil {
			return nil, err
		}

		usr, err := gr.Users().GetByID(ctx, *userID)
		if err != nil {
			return nil, err
		}

		newContext = SetUserOnContext(newContext, usr)

		// Make the actual request
		res, err := handler(newContext, req)
		return res, err
	})
}
