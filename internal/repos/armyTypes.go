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

//go:generate mockgen -source=./armyTypes.go -destination=./mocks/ArmyTypesRepo.go -package=mock_repos ArmyTypesRepo
type ArmyTypesRepo interface {
	Find(ctx context.Context, gameID string, limit, offset int) ([]*types.ArmyType, int64, error)
	FindOrCreate(ctx context.Context, at types.CreateArmyType) (*types.ArmyType, error)
}

func NewArmyTypesRepo(db neo4j.DriverWithContext) ArmyTypesRepo {
	return &armyTypesRepo{db}
}

type armyTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *armyTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.ArmyType, int64, error) {
	var count int64
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
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(at:ArmyType)
			WHERE at.name <> 'All Armies'
			RETURN at
			ORDER BY at.name
			%s %s;
		`, offsetQry, limitQry), map[string]any{"game_id": gameID})
		if err != nil {
			return nil, err
		}

		ats := []*types.ArmyType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			ats = append(ats, types.ArmyTypeFromNode(node))
		}

		result, err = tx.Run(ctx, `
			MATCH (n:ArmyType)
			RETURN count(n) as count
		`, map[string]any{})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			var ok bool
			count, ok = result.Record().Values[0].(int64)
			if !ok {
				return nil, errors.New("unable to convert database count to int64")
			}
		}

		return ats, nil
	})
	if err != nil {
		return nil, 0, err
	} else if res == nil {
		return []*types.ArmyType{}, 0, nil
	}

	return res.([]*types.ArmyType), count, nil
}

func (r *armyTypesRepo) FindOrCreate(ctx context.Context, at types.CreateArmyType) (*types.ArmyType, error) {
	existingArmyType, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingArmyType, nil
}

func (r *armyTypesRepo) getByName(ctx context.Context, at types.CreateArmyType) (*types.ArmyType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(at:ArmyType)
			WHERE at.name = $name
			RETURN at;
		`, map[string]any{"name": at.Name, "game_id": at.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.ArmyTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("armyType")
	}

	return res.(*types.ArmyType), nil
}

func (r *armyTypesRepo) create(ctx context.Context, at types.CreateArmyType) (*types.ArmyType, error) {
	if err := types.Validate(at); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (at:ArmyType{
				name: 			$name
				,game_id: 		$game_id
			})
			ON CREATE
				SET at.created_at = $created_at,
				at.id = apoc.create.uuid()
			ON MATCH
				SET at.updated_at = $updated_at
			MERGE (at)-[:BELONGS_TO]->(g)
			RETURN at;
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
			return types.ArmyTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("armyType")
	}

	return res.(*types.ArmyType), nil
}
