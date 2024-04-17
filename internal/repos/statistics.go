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

type FindStatisticsOpts struct {
	GameIDs  []string
	Names    []string
	Displays []string
	IDs      []string
	Limit    int
	Offset   int
}

//go:generate mockgen -source=./statistics.go -destination=./mocks/StatisticsRepo.go -package=mock_repos StatisticsRepo
type StatisticsRepo interface {
	Get(ctx context.Context, id string) (*types.Statistic, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.Statistic, error)
	Find(ctx context.Context, opts *FindStatisticsOpts) ([]*types.Statistic, int64, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindStatisticsOpts) ([]*types.Statistic, int64, error)
	FindOrCreate(ctx context.Context, at types.CreateStatistic) (*types.Statistic, error)
	Update(ctx context.Context, stat types.UpdateStatistic) error
	UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, stat types.UpdateStatistic) error
}

func NewStatisticsRepo(db neo4j.DriverWithContext) StatisticsRepo {
	return &statisticsRepo{db}
}

type statisticsRepo struct {
	db neo4j.DriverWithContext
}

func (r *statisticsRepo) Update(ctx context.Context, stat types.UpdateStatistic) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.UpdateTx(ctx, tx, stat)
	})
	return err
}

func (r *statisticsRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, stat types.UpdateStatistic) error {
	if err := types.Validate(stat); err != nil {
		return err
	}

	result, err := tx.Run(ctx, `
			MATCH (stat:Statistic{id: $id})
			SET stat.name = $name
				,stat.display = $display
			RETURN stat
		`, map[string]interface{}{
		"id":      stat.ID,
		"name":    stat.Name,
		"display": stat.Display,
	})
	if err != nil {
		return err
	}
	return result.Err()
}

func (r *statisticsRepo) Get(ctx context.Context, id string) (*types.Statistic, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetTx(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.Statistic), nil
}

func (r *statisticsRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.Statistic, error) {
	stats, _, err := r.FindTx(ctx, tx, &FindStatisticsOpts{
		IDs: []string{id}, Limit: 1,
	})
	if err != nil {
		return nil, err
	} else if len(stats) == 0 {
		return nil, types.NewNotFoundError("statistic not found")
	}

	return stats[0], nil
}

func (r *statisticsRepo) Find(ctx context.Context, opts *FindStatisticsOpts) ([]*types.Statistic, int64, error) {
	var count int64
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, c, err := r.FindTx(ctx, tx, opts)
		if err != nil {
			return nil, err
		}
		count = c
		return result, nil
	})
	if err != nil {
		return nil, 0, err
	} else if res == nil {
		return []*types.Statistic{}, 0, nil
	}

	return res.([]*types.Statistic), count, nil
}

func (r *statisticsRepo) FindTx(
	ctx context.Context, tx neo4j.ManagedTransaction, opts *FindStatisticsOpts,
) ([]*types.Statistic, int64, error) {
	var count int64

	params := map[string]any{}
	var limitQry string
	var offsetQry string
	var whereQry string
	var matchQry string

	if opts.Limit > 0 {
		limitQry = fmt.Sprintf("LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		offsetQry = fmt.Sprintf("SKIP %d", opts.Offset)
	}

	comp := "WHERE"
	if len(opts.GameIDs) > 0 {
		params["gameIds"] = opts.GameIDs
		whereQry += fmt.Sprintf(" %s g.id IN $gameIds", comp)
		comp = "AND"
		matchQry = "(g:Game)<-[:BELONGS_TO]-"
	}

	if len(opts.IDs) > 0 {
		params["ids"] = opts.IDs
		whereQry += fmt.Sprintf(" %s stat.id IN $ids", comp)
		comp = "AND"
	}

	if len(opts.Names) > 0 {
		params["names"] = opts.Names
		whereQry += fmt.Sprintf(" %s stat.name IN $names", comp)
		comp = "AND"
	}

	if len(opts.Displays) > 0 {
		params["displays"] = opts.Displays
		whereQry += fmt.Sprintf(" %s stat.display IN $displays", comp)
		comp = "AND"
	}

	result, err := tx.Run(ctx, fmt.Sprintf(`
		MATCH %s(stat:Statistic)
		%s
		RETURN stat
		ORDER BY stat.position
		%s %s;
	`, matchQry, whereQry, offsetQry, limitQry), params)
	if err != nil {
		return nil, 0, err
	}

	ats := []*types.Statistic{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type")
		}
		ats = append(ats, types.StatisticFromNode(node))
	}

	result, err = tx.Run(ctx, fmt.Sprintf(`
		MATCH %s(stat:Statistic)
		%s
		RETURN count(stat) as count
	`, matchQry, whereQry), params)
	if err != nil {
		return nil, 0, err
	}

	if result.Next(ctx) {
		var ok bool
		count, ok = result.Record().Values[0].(int64)
		if !ok {
			return nil, 0, errors.New("unable to convert database count to int64")
		}
	}

	return ats, count, nil
}

func (r *statisticsRepo) FindOrCreate(ctx context.Context, at types.CreateStatistic) (*types.Statistic, error) {
	existingStatistic, err := r.getByName(ctx, at)
	if types.IsNotFoundError(err) {
		return r.create(ctx, at)
	} else if err != nil {
		return nil, err
	}

	return existingStatistic, nil
}

func (r *statisticsRepo) getByName(ctx context.Context, stat types.CreateStatistic) (*types.Statistic, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(stat:Statistic{ name: $name, display: $display })
			RETURN stat;
		`, map[string]any{"name": stat.Name, "display": stat.Display, "game_id": stat.GameID})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			return types.StatisticFromNode(node), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("statistic")
	}

	return res.(*types.Statistic), nil
}

func (r *statisticsRepo) create(ctx context.Context, stat types.CreateStatistic) (*types.Statistic, error) {
	if err := types.Validate(stat); err != nil {
		return nil, err
	}

	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		var position int64 = 1
		result, err := tx.Run(ctx, `
			MATCH (g:Game{ id: $game_id })<-[:BELONGS_TO]-(stat:Statistic) RETURN size(collect(stat))
		`, map[string]any{
			"game_id": stat.GameID,
		})
		if err != nil {
			log.Printf("unable to get position for statistic err: %s\n", err.Error())
		}
		if result != nil && result.Next(ctx) {
			position = result.Record().Values[0].(int64) + 1
		}

		result2, err := tx.Run(ctx, `
			MATCH (g:Game{ id:$game_id })
			MERGE (stat:Statistic{
				name: 			$name
				,display:		$display
				,game_id: 		g.id
			})
			ON CREATE
				SET stat.created_at = $created_at,
				stat.id = apoc.create.uuid(),
				stat.position = $position
			ON MATCH
				SET stat.updated_at = $updated_at
			MERGE (stat)-[:BELONGS_TO]->(g)
			RETURN stat;
		`, map[string]any{
			"name":       stat.Name,
			"display":    stat.Display,
			"game_id":    stat.GameID,
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
			return types.StatisticFromNode(node), nil
		}

		return nil, result2.Err()
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("statistic")
	}

	return res.(*types.Statistic), nil
}
