package interceptors

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/happilymarrieddad/old-world/api3/internal/jwt"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	"google.golang.org/grpc"
)

var unauthorizedRoutes = []string{
	"/auth.Auth/Login",
	"/auth.Auth/CreateUser",
	"/auth.Auth/Validate",
}

func GlobalRepoInjector(gr repos.GlobalRepo) grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("route: %s", info.FullMethod)
		ctx = SetGlobalRepoOnContext(ctx, gr)

		// check unauthorized routes
		for _, route := range unauthorizedRoutes {
			if route == info.FullMethod {
				return handler(ctx, req)
			}
		}

		// BEFORE the request
		v := reflect.Indirect(reflect.ValueOf(req))
		vField := reflect.Indirect(v.FieldByName("JWT"))

		// No JWT field so throw an error
		if !vField.IsValid() {
			return nil, types.NewUnauthorizedError("unauthorized")
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

		// Make the actual request
		res, err := handler(SetUserOnContext(ctx, usr), req)
		if err != nil {
			fmt.Println(err.Error())
		}
		return res, err
	})
}
