package repos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FindItemTypeOpts struct {
	GameID string
	Limit  int
	Offset int
	Name   []string
	ID     []string
}

//go:generate mockgen -source=./itemTypes.go -destination=./mocks/ItemTypesRepo.go -package=mock_repos ItemTypesRepo
type ItemTypesRepo interface {
	Find(ctx context.Context, opts *FindItemTypeOpts) ([]*types.ItemType, int64, error)
	FindOrCreate(ctx context.Context, at types.CreateItemType) (*types.ItemType, error)
}

func NewItemTypesRepo(db neo4j.DriverWithContext) ItemTypesRepo {
	return &itemTypesRepo{db}
}

type itemTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *itemTypesRepo) Find(ctx context.Context, opts *FindItemTypeOpts) ([]*types.ItemType, int64, error) {
	if opts == nil {
		opts = &FindItemTypeOpts{}
	}

	var count int64
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		var limitQry string
		var offsetQry string

		if opts.Limit > 0 {
			limitQry = fmt.Sprintf("LIMIT %d", opts.Limit)
		}

		if opts.Offset > 0 {
			offsetQry = fmt.Sprintf("SKIP %d", opts.Offset)
		}

		params := make(map[string]any)
		params["game_id"] = opts.GameID
		var whereQry string
		comp := "WHERE"
		if len(opts.Name) > 0 {
			whereQry += fmt.Sprintf(" %s it.name IN $names", comp)
			params["names"] = opts.Name
			comp = "AND"
		}

		if len(opts.ID) > 0 {
			whereQry += fmt.Sprintf(" %s it.id IN $ids", comp)
			params["ids"] = opts.ID
			comp = "AND"
		}

		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(it:ItemType)
			%s
			RETURN it
			ORDER BY it.name
			%s %s;
		`, whereQry, offsetQry, limitQry), params)
		if err != nil {
			return nil, err
		}

		its := []*types.ItemType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			its = append(its, types.ItemTypeFromNode(node))
		}

		result, err = tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(n:ItemType)
			RETURN count(n) as count
		`, map[string]any{"game_id": opts.GameID})
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

		return its, nil
	})
	if err != nil {
		return nil, 0, err
	} else if res == nil {
		return []*types.ItemType{}, 0, nil
	}

	return res.([]*types.ItemType), count, nil
}

func (r *itemTypesRepo) FindOrCreate(ctx context.Context, at types.CreateItemType) (*types.ItemType, error) {
	existingItemType, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingItemType, nil
}

func (r *itemTypesRepo) getByName(ctx context.Context, at types.CreateItemType) (*types.ItemType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(it:ItemType)
			WHERE it.name = $name
			RETURN it;
		`, map[string]any{"name": at.Name, "game_id": at.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.ItemTypeFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("itemType")
	}

	return res.(*types.ItemType), nil
}

func (r *itemTypesRepo) create(ctx context.Context, at types.CreateItemType) (*types.ItemType, error) {
	if err := types.Validate(at); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		var position int64 = 1
		result, err := tx.Run(ctx, `MATCH (it:ItemType) RETURN size(collect(it))`, make(map[string]any))
		if err != nil {
			log.Printf("unable to get position for item type err: %s\n", err.Error())
		}
		if result != nil && result.Next(ctx) {
			position = result.Record().Values[0].(int64) + 1
		}

		result2, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (it:ItemType{
				name: 			$name
				,game_id: 		$game_id
			})
			ON CREATE
				SET it.created_at = $created_at,
				it.id = apoc.create.uuid(),
				it.position = $position
			ON MATCH
				SET it.updated_at = $updated_at
			MERGE (it)-[:BELONGS_TO]->(g)
			RETURN it;
		`, map[string]any{
			"name":       at.Name,
			"game_id":    at.GameID,
			"position":   position,
			"created_at": time.Now().UTC().Unix(),
			"updated_at": time.Now().UTC().Unix(),
		})
		if err != nil {
			return nil, err
		}

		if result2.Next(ctx) {
			node, ok := result2.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.ItemTypeFromNode(node), nil
		}

		return nil, result2.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("itemType")
	}

	return res.(*types.ItemType), nil
}
