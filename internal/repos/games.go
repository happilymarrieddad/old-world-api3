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

type FindGameOpts struct {
	IDs    []string
	Names  []string
	Limit  int
	Offset int
}

//go:generate mockgen -source=./games.go -destination=./mocks/GamesRepo.go -package=mock_repos GamesRepo
type GamesRepo interface {
	Get(ctx context.Context, id string) (*types.Game, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.Game, error)
	Find(ctx context.Context, opts *FindGameOpts) ([]*types.Game, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindGameOpts) ([]*types.Game, error)
	FindOrCreate(ctx context.Context, name string) (*types.Game, error)
	Update(ctx context.Context, gm types.UpdateGame) error
	UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, gm types.UpdateGame) error
}

func NewGamesRepo(db neo4j.DriverWithContext) GamesRepo {
	return &gamesRepo{db}
}

type gamesRepo struct {
	db neo4j.DriverWithContext
}

func (r *gamesRepo) Get(ctx context.Context, id string) (*types.Game, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetTx(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.Game), nil
}

func (r *gamesRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.Game, error) {
	gms, err := r.FindTx(ctx, tx, &FindGameOpts{
		IDs: []string{id}, Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(gms) == 0 {
		return nil, types.NewNotFoundError("game not found")
	}

	return gms[0], nil
}

func (r *gamesRepo) Update(ctx context.Context, gm types.UpdateGame) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.UpdateTx(ctx, tx, gm)
	})
	return err
}

func (r *gamesRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, gm types.UpdateGame) error {
	if err := types.Validate(gm); err != nil {
		return err
	}

	result, err := tx.Run(ctx, `
			MATCH (gm:Game{id: $id})
			SET gm.name = $name
			RETURN gm
		`, map[string]interface{}{
		"id":   gm.ID,
		"name": gm.Name,
	})
	if err != nil {
		return err
	}
	return result.Err()
}

func (r *gamesRepo) Find(ctx context.Context, opts *FindGameOpts) ([]*types.Game, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.FindTx(ctx, tx, opts)
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return []*types.Game{}, nil
	}

	return res.([]*types.Game), nil
}

func (r *gamesRepo) FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindGameOpts) ([]*types.Game, error) {
	if opts == nil {
		opts = &FindGameOpts{}
	}

	var limitQry string
	var offsetQry string
	var whereQry string
	params := map[string]any{}
	comp := "WHERE"

	if len(opts.IDs) > 0 {
		whereQry += fmt.Sprintf(" %s g.id IN $ids", comp)
		params["ids"] = opts.IDs
		comp = "AND"
	}

	if len(opts.Names) > 0 {
		whereQry += fmt.Sprintf(" %s g.name IN $names", comp)
		params["names"] = opts.Names
		comp = "AND"
	}

	if opts.Limit > 0 {
		limitQry = fmt.Sprintf("LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		offsetQry = fmt.Sprintf("SKIP %d", opts.Offset)
	}

	result, err := tx.Run(ctx, fmt.Sprintf(`
		MATCH (g:Game)
		%s
		RETURN g
		ORDER BY g.name
		%s %s;
	`, whereQry, offsetQry, limitQry), params)
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
