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

//go:generate mockgen -source=./unitOptionTypes.go -destination=./mocks/UnitOptionTypesRepo.go -package=mock_repos UnitOptionTypesRepo
type UnitOptionTypesRepo interface {
	Find(ctx context.Context, gameID string, limit, offset int) ([]*types.UnitOptionType, error)
	FindOrCreate(ctx context.Context, at types.CreateUnitOptionType) (*types.UnitOptionType, error)
}

func NewUnitOptionTypesRepo(db neo4j.DriverWithContext) UnitOptionTypesRepo {
	return &unitOptionTypesRepo{db}
}

type unitOptionTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *unitOptionTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.UnitOptionType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		var limitQry string
		var offsetQry string

		if limit > 0 {
			limitQry = fmt.Sprintf("LIMIT %d", limit)
		}

		if offset > 0 {
			offsetQry = fmt.Sprintf("OFFSET %d", offset)
		}

		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(uot:UnitOptionType)
			RETURN uot
			ORDER BY uot.name
			%s %s;
		`, limitQry, offsetQry), map[string]any{"game_id": gameID})
		if err != nil {
			return nil, err
		}

		ats := []*types.UnitOptionType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			ats = append(ats, types.UnitOptionTypeFromNode(node))
		}

		return ats, nil
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return []*types.UnitOptionType{}, nil
	}

	return res.([]*types.UnitOptionType), nil
}

func (r *unitOptionTypesRepo) FindOrCreate(ctx context.Context, at types.CreateUnitOptionType) (*types.UnitOptionType, error) {
	existingUnitOptionType, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingUnitOptionType, nil
}

func (r *unitOptionTypesRepo) getByName(ctx context.Context, at types.CreateUnitOptionType) (*types.UnitOptionType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(uot:UnitOptionType)
			WHERE uot.name = $name
			RETURN uot;
		`, map[string]any{"name": at.Name, "game_id": at.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.UnitOptionTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("unitOptionType")
	}

	return res.(*types.UnitOptionType), nil
}

func (r *unitOptionTypesRepo) create(ctx context.Context, at types.CreateUnitOptionType) (*types.UnitOptionType, error) {
	if err := types.Validate(at); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (uot:UnitOptionType{
				name: 			$name
				,game_id: 		$game_id
			})
			ON CREATE
				SET uot.created_at = $created_at,
				uot.id = apoc.create.uuid()
			ON MATCH
				SET uot.updated_at = $updated_at
			MERGE (uot)-[:BELONGS_TO]->(g)
			RETURN uot;
		`, map[string]any{
			"name":       at.Name,
			"game_id":    at.GameID,
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
			return types.UnitOptionTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("unitOptionType")
	}

	return res.(*types.UnitOptionType), nil
}
