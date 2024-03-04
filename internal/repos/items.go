package repos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FindItemsOpts struct {
	IDs         []string
	Names       []string
	GameID      string `validate:"required"`
	ArmyTypeID  *string
	ItemTypeIDs []string
	Limit       int
	Offset      int
	Debug       bool
}

//go:generate mockgen -source=./items.go -destination=./mocks/ItemsRepo.go -package=mock_repos ItemsRepo
type ItemsRepo interface {
	Get(ctx context.Context, id, gameID string) (*types.Item, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id, gameID string) (*types.Item, error)
	Find(ctx context.Context, opts *FindItemsOpts) ([]*types.Item, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindItemsOpts) ([]*types.Item, error)
	Create(ctx context.Context, itm types.CreateItem) (*types.Item, error)
	CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, itm types.CreateItem) (*types.Item, error)
}

func NewItemsRepo(db neo4j.DriverWithContext) ItemsRepo {
	return &itemsRepo{db}
}

type itemsRepo struct {
	db neo4j.DriverWithContext
}

func (r *itemsRepo) Get(ctx context.Context, id, gameID string) (*types.Item, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetTx(ctx, tx, id, gameID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.Item), nil
}

func (r *itemsRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id, gameID string) (*types.Item, error) {
	res, err := r.FindTx(ctx, tx, &FindItemsOpts{IDs: []string{id}, Limit: 1, GameID: gameID})
	if err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, types.NewNotFoundError("items")
	}

	return res[0], nil
}

func (r *itemsRepo) Find(ctx context.Context, opts *FindItemsOpts) ([]*types.Item, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.FindTx(ctx, tx, opts)
	})
	if err != nil {
		return nil, err
	}
	return res.([]*types.Item), nil
}

func (r *itemsRepo) FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindItemsOpts) ([]*types.Item, error) {
	if err := types.Validate(opts); err != nil {
		return nil, err
	}

	var limitQry string
	var offsetQry string

	if opts.Limit > 0 {
		limitQry = fmt.Sprintf("LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		offsetQry = fmt.Sprintf("OFFSET %d", opts.Offset)
	}

	params := map[string]any{
		"game_id": opts.GameID,
	}

	var atQry string
	whrQry := "WHERE EXISTS((i)-[:BELONGS_TO]->(g)) AND NOT (i)-[:BELONGS_TO]->(:ArmyType)"
	if opts.ArmyTypeID != nil {
		params["army_type_id"] = *opts.ArmyTypeID
		atQry = `
		MATCH (at:ArmyType{ id: $army_type_id })
		`
		whrQry = `
			WHERE ((
				(i)-[:BELONGS_TO]->(g) AND NOT (i)-[:BELONGS_TO]->(:ArmyType)
			) OR (
				(i)-[:BELONGS_TO]->(at) AND (i)-[:BELONGS_TO]->(g)
			))`
	}

	var idQry string
	if len(opts.IDs) > 0 {
		idQry = " AND i.id IN $ids"
		params["ids"] = opts.IDs
	}

	var namesQry string
	if len(opts.Names) > 0 {
		namesQry = " AND i.name IN $names"
		params["names"] = opts.Names
	}

	var typeIdsQry string
	if len(opts.ItemTypeIDs) > 0 {
		idQry += " AND it.id IN $type_ids"
		params["type_ids"] = opts.ItemTypeIDs
	}

	cmd := fmt.Sprintf(`
		MATCH (g:Game{ id: $game_id })
		%s
		MATCH (g)<-[:BELONGS_TO]-(i:Item)-[:IS_ITEM_TYPE]->(it:ItemType)
		%s %s %s %s
		RETURN i,it
		ORDER BY it.position, i.points DESC
		%s %s;
	`, atQry, whrQry, idQry, namesQry, typeIdsQry, limitQry, offsetQry)

	if opts.Debug {
		fmt.Println("??????? Items Find Query")
		fmt.Println(cmd)
		spew.Dump(params)
	}

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	itms := []*types.Item{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database object to item")
		}
		itm := types.ItemFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database object to item type")
		}
		it := types.ItemTypeFromNode(node)
		itm.ItemTypeID = it.ID
		itm.ItemTypeName = it.Name

		itms = append(itms, itm)
	}

	return itms, nil
}

func (r *itemsRepo) Create(ctx context.Context, itm types.CreateItem) (*types.Item, error) {
	newItm, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		i, e := r.CreateTx(ctx, tx, itm)
		return i, e
	})
	if err != nil {
		return nil, err
	}
	return newItm.(*types.Item), nil
}

func (r *itemsRepo) CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, itm types.CreateItem) (*types.Item, error) {
	params := map[string]any{
		"name":         itm.Name,
		"points":       itm.Points,
		"game_id":      itm.GameID,
		"item_type_id": itm.ItemTypeID,
		"created_at":   time.Now().UTC().Unix(),
		"updated_at":   time.Now().UTC().Unix(),
	}

	var mergeQryAdd string
	var atMatchQry string
	var atQry string
	if itm.ArmyTypeID != nil {
		params["army_type_id"] = *itm.ArmyTypeID
		atMatchQry = "MATCH (at:ArmyType{ id: $army_type_id })"
		atQry = "MERGE (i)-[:BELONGS_TO]->(at)"
		mergeQryAdd = ",i.army_type_id = at.id"
	}

	cmd := fmt.Sprintf(`
		MATCH (g:Game{ id: $game_id })
		MATCH (it:ItemType{ id: $item_type_id })
		%s
		MERGE (i:Item{
			name:			$name
			,points:		$points
			,game_id:		$game_id
			,item_type_id:  $item_type_id
		})
		ON CREATE
			SET i.created_at = $created_at,
			i.id = apoc.create.uuid()
			%s
		ON MATCH
			SET i.updated_at = $updated_at
		MERGE (i)-[:BELONGS_TO]->(g)
		MERGE (it)<-[:IS_ITEM_TYPE]-(i)
		%s
		RETURN i,it;
	`, atMatchQry, mergeQryAdd, atQry)

	if itm.Debug {
		fmt.Println("???????")
		fmt.Println(cmd)
		spew.Dump(itm)
	}

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, types.NewInternalServerError("unable to convert item database object")
		}
		ni := types.ItemFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, types.NewInternalServerError("unable to convert item type database object")
		}
		nit := types.ItemTypeFromNode(node)
		ni.ItemTypeID = nit.ID
		ni.ItemTypeName = nit.Name

		return ni, nil
	}

	if result.Err() == nil {
		return nil, types.NewNotFoundError("unable to find/match item")
	}

	return nil, result.Err()
}
