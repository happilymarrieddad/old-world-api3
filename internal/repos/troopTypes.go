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

//go:generate mockgen -source=./troopTypes.go -destination=./mocks/TroopTypesRepo.go -package=mock_repos TroopTypesRepo
type TroopTypesRepo interface {
	Find(ctx context.Context, gameID string, limit, offset int) ([]*types.TroopType, int64, error)
	FindOrCreate(ctx context.Context, at types.CreateTroopType) (*types.TroopType, error)
}

func NewTroopTypesRepo(db neo4j.DriverWithContext) TroopTypesRepo {
	return &troopTypesRepo{db}
}

type troopTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *troopTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.TroopType, int64, error) {
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
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(tt:TroopType)
			RETURN tt
			ORDER BY tt.name
			%s %s;
		`, offsetQry, limitQry), map[string]any{"game_id": gameID})
		if err != nil {
			return nil, err
		}

		tts := []*types.TroopType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			tts = append(tts, types.TroopTypeFromNode(node))
		}

		result, err = tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(n:TroopType)
			RETURN count(n) as count
		`, map[string]any{"game_id": gameID})
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

		return tts, nil
	})
	if err != nil {
		return nil, 0, err
	} else if res == nil {
		return []*types.TroopType{}, 0, nil
	}

	return res.([]*types.TroopType), count, nil
}

func (r *troopTypesRepo) FindOrCreate(ctx context.Context, at types.CreateTroopType) (*types.TroopType, error) {
	existingTroopType, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingTroopType, nil
}

func (r *troopTypesRepo) getByName(ctx context.Context, at types.CreateTroopType) (*types.TroopType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(tt:TroopType)
			WHERE tt.name = $name
			RETURN tt;
		`, map[string]any{"name": at.Name, "game_id": at.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.TroopTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("troopType")
	}

	return res.(*types.TroopType), nil
}

func (r *troopTypesRepo) create(ctx context.Context, at types.CreateTroopType) (*types.TroopType, error) {
	if err := types.Validate(at); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (tt:TroopType{
				name: 			$name
				,game_id: 		$game_id
			})
			ON CREATE
				SET tt.created_at = $created_at,
				tt.id = apoc.create.uuid()
			ON MATCH
				SET tt.updated_at = $updated_at
			MERGE (tt)-[:BELONGS_TO]->(g)
			RETURN tt;
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
			return types.TroopTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("troopType")
	}

	return res.(*types.TroopType), nil
}
