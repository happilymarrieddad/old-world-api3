package repos

import (
	"context"
	"errors"
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

//go:generate mockgen -source=./users.go -destination=./mocks/UsersRepo.go -package=mock_repos UsersRepo
type UsersRepo interface {
	Create(context.Context, types.CreateUser) (*types.User, error)
	GetByID(ctx context.Context, id string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	FindOrCreate(ctx context.Context, usr types.CreateUser) (*types.User, error)
	Delete(ctx context.Context, id ...string) error
}

func NewUsersRepo(db neo4j.DriverWithContext) UsersRepo {
	return &usersRepo{db: db}
}

type usersRepo struct {
	db neo4j.DriverWithContext
}

func (r *usersRepo) Create(ctx context.Context, nu types.CreateUser) (*types.User, error) {
	if err := types.Validate(nu); err != nil {
		return nil, err
	}

	pwhash, err := nu.GetPasswordHash()
	if err != nil {
		return nil, errors.New("unable to hash password")
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			CREATE (u:User{
				id: 			apoc.create.uuid()
				,first_name: 	$first_name
				,last_name: 		$last_name
				,email: 		$email
				,password: 		$password
				,created_at: 	$created_at
			})
			RETURN u;
		`, map[string]any{
			"first_name": nu.FirstName,
			"last_name":  nu.LastName,
			"email":      nu.Email,
			"password":   pwhash,
			"created_at": time.Now().UTC().Unix(),
		})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.UserFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	return res.(*types.User), nil
}

func (r *usersRepo) GetByID(ctx context.Context, id string) (*types.User, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User{ id: $id })
			RETURN u;
		`, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.UserFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("user")
	}

	return res.(*types.User), nil
}

func (r *usersRepo) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User{ email: $email })
			RETURN u;
		`, map[string]any{"email": email})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.UserFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("user")
	}

	return res.(*types.User), nil
}

func (r *usersRepo) FindOrCreate(ctx context.Context, usr types.CreateUser) (*types.User, error) {
	if err := types.Validate(usr); err != nil {
		return nil, err
	}

	existingUser, err := r.GetByEmail(ctx, usr.Email)
	if types.IsNotFoundError(err) {
		return r.Create(ctx, usr)
	} else if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (r *usersRepo) Delete(ctx context.Context, ids ...string) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		_, e := tx.Run(ctx, `
			MATCH (u:User) WHERE u.id IN $ids
			DETACH DELETE u;
		`, map[string]any{
			"ids": ids,
		})
		return nil, e
	})
	return err
}
