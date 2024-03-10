package repos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

//go:generate mockgen -source=./games.go -destination=./mocks/GamesRepo.go -package=mock_repos GamesRepo
type GamesRepo interface {
	Find(ctx context.Context, limit, offset int) ([]*types.Game, error)
	FindOrCreate(ctx context.Context, name string) (*types.Game, error)
}

func NewGamesRepo(db neo4j.DriverWithContext) GamesRepo {
	return &gamesRepo{db}
}

type gamesRepo struct {
	db neo4j.DriverWithContext
}

func (r *gamesRepo) Find(ctx context.Context, limit, offset int) ([]*types.Game, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		var limitQry string
		var offsetQry string

		if limit > 0 {
			limitQry = fmt.Sprintf("LIMIT %d", limit)
		}

		if offset > 0 {
			offsetQry = fmt.Sprintf("SKIP %d", offset)
		}

		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (g:Game)
			RETURN g
			ORDER BY g.name
			%s %s;
		`, offsetQry, limitQry), map[string]any{})
		if err != nil {
			return nil, err
		}

		gms := []*types.Game{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			gms = append(gms, types.GameFromNode(node))
		}

		return gms, nil
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return []*types.Game{}, nil
	}

	return res.([]*types.Game), nil
}

func (r *gamesRepo) FindOrCreate(ctx context.Context, name string) (*types.Game, error) {
	existingUser, err := r.getByName(ctx, name)
	if types.IsNotFoundError(err) {
		return r.create(ctx, types.CreateGame{Name: name})
	} else if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (r *gamesRepo) getByName(ctx context.Context, name string) (*types.Game, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ name: $name })
			RETURN g;
		`, map[string]any{"name": name})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.GameFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("game")
	}

	return res.(*types.Game), nil
}

func (r *gamesRepo) create(ctx context.Context, usr types.CreateGame) (*types.Game, error) {
	if err := types.Validate(usr); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MERGE (g:Game{
				name: 			$name
			})
			ON CREATE
				SET g.created_at = $created_at,
				g.id = apoc.create.uuid()
			ON MATCH
				SET g.updated_at = $updated_at
			RETURN g;
		`, map[string]any{
			"name":       usr.Name,
			"created_at": time.Now().UTC().Unix(),
			"updated_at": time.Now().UTC().Unix(),
		})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			gm := types.GameFromNode(node)

			// The four min unit option type
			uots := []string{"Single", "One Of", "Many Of", "Many To"}
			uotRepo := NewUnitOptionTypesRepo(r.db)
			for _, uot := range uots {
				if _, err := uotRepo.FindOrCreateTx(ctx, tx, types.CreateUnitOptionType{
					GameID: gm.ID,
					Name:   uot,
				}); err != nil {
					return nil, err
				}
			}

			return gm, nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("game")
	}

	return res.(*types.Game), nil
}
