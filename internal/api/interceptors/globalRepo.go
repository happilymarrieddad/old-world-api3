package interceptors

import (
	"context"
	"errors"

	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
)

type key string

const (
	globalRepoKey key = "globalRepo"
	userKey       key = "user"
)

// GetGlobalRepoFromContext - get the global repo from context
func GetGlobalRepoFromContext(
	ctx context.Context,
) (globalRepo repos.GlobalRepo, err error) {
	r, ok := ctx.Value(globalRepoKey).(repos.GlobalRepo)
	if ok {
		return r, nil
	}

	return nil, errors.New("unable to get global repo from context")
}

// SetGlobalRepoOnContext - set the global repo value in the context
func SetGlobalRepoOnContext(
	ctx context.Context,
	globalRepo repos.GlobalRepo,
) context.Context {
	return context.WithValue(ctx, globalRepoKey, globalRepo)
}

// SetUserOnContext - sets the users on the current context
func SetUserOnContext(
	ctx context.Context,
	user *types.User,
) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUserFromContext - get the user from context
func GetUserFromContext(
	ctx context.Context,
) (usr *types.User, err error) {
	usr, ok := ctx.Value(userKey).(*types.User)
	if ok {
		return usr, nil
	}

	return nil, errors.New("unable to get user from context")
}
