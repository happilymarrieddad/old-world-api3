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

//go:generate mockgen -source=./compositionTypes.go -destination=./mocks/CompositionTypesRepo.go -package=mock_repos CompositionTypesRepo
type CompositionTypesRepo interface {
	Find(ctx context.Context, gameID string, limit, offset int) ([]*types.CompositionType, error)
	FindOrCreate(ctx context.Context, at types.CreateCompositionType) (*types.CompositionType, error)
}

func NewCompositionTypesRepo(db neo4j.DriverWithContext) CompositionTypesRepo {
	return &compositionTypesRepo{db}
}

type compositionTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *compositionTypesRepo) Find(ctx context.Context, gameID string, limit, offset int) ([]*types.CompositionType, error) {
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
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(ct:CompositionType)
			RETURN ct
			ORDER BY ct.name
			%s %s;
		`, limitQry, offsetQry), map[string]any{"game_id": gameID})
		if err != nil {
			return nil, err
		}

		ats := []*types.CompositionType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			ats = append(ats, types.CompositionTypeFromNode(node))
		}

		return ats, nil
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return []*types.CompositionType{}, nil
	}

	return res.([]*types.CompositionType), nil
}

func (r *compositionTypesRepo) FindOrCreate(ctx context.Context, at types.CreateCompositionType) (*types.CompositionType, error) {
	existingCompositionType, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingCompositionType, nil
}

func (r *compositionTypesRepo) getByName(ctx context.Context, at types.CreateCompositionType) (*types.CompositionType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(ct:CompositionType)
			WHERE ct.name = $name
			RETURN ct;
		`, map[string]any{"name": at.Name, "game_id": at.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.CompositionTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("compositionType")
	}

	return res.(*types.CompositionType), nil
}

func (r *compositionTypesRepo) create(ctx context.Context, at types.CreateCompositionType) (*types.CompositionType, error) {
	if err := types.Validate(at); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (ct:CompositionType{
				name: 			$name
				,game_id: 		$game_id
			})
			ON CREATE
				SET ct.created_at = $created_at,
				ct.id = apoc.create.uuid()
			ON MATCH
				SET ct.updated_at = $updated_at
			MERGE (ct)-[:BELONGS_TO]->(g)
			RETURN ct;
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
			return types.CompositionTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("compositionType")
	}

	return res.(*types.CompositionType), nil
}
