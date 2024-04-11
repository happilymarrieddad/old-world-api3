package repos

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FindUnitTypesOpts struct {
	ArmyTypeID             string
	UnitTypeIDs            []string
	ParentUnitTypeID       string
	Name                   string
	Names                  []string
	Limit                  int
	Offset                 int
	IsChildRequest         bool
	Debug                  bool
	IncludeUnitTypeOptions bool
}

//go:generate mockgen -source=./unitTypes.go -destination=./mocks/UnitTypesRepo.go -package=mock_repos UnitTypesRepo
type UnitTypesRepo interface {
	GetNamesByArmyTypeID(ctx context.Context, armyTypeID string) ([]*types.UnitType, error)
	Find(ctx context.Context, opts *FindUnitTypesOpts) ([]*types.UnitType, int64, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindUnitTypesOpts) ([]*types.UnitType, int64, error)
	Get(ctx context.Context, id string) (*types.UnitType, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.UnitType, error)
	FindOrCreate(ctx context.Context, at types.CreateUnitType) (*types.UnitType, error)
	EnsureChildUnitTypeExists(ctx context.Context, ncut types.CreateChildUnitType) error
	EnsureChildUnitTypeExistsTx(ctx context.Context, tx neo4j.ManagedTransaction, ncut types.CreateChildUnitType) error
	Update(ctx context.Context, ut types.UpdateUnitType) error
	UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, ut types.UpdateUnitType) error
}

func NewUnitTypesRepo(db neo4j.DriverWithContext) UnitTypesRepo {
	return &unitTypesRepo{db}
}

type unitTypesRepo struct {
	db neo4j.DriverWithContext
}

func (r *unitTypesRepo) Update(ctx context.Context, ut types.UpdateUnitType) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.UpdateTx(ctx, tx, ut)
	})
	return err
}

func (r *unitTypesRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, ut types.UpdateUnitType) error {
	if err := types.Validate(ut); err != nil {
		return err
	}

	result, err := tx.Run(ctx, `
			MATCH (ut:UnitType{id: $id})
			SET ut.name = $name
				,ut.points_per_model = $pointsPerModel
				,ut.min_models = $minModels
				,ut.max_models = $maxModels
			RETURN ut
		`, map[string]interface{}{
		"id":             ut.ID,
		"name":           ut.Name,
		"pointsPerModel": ut.PointsPerModel,
		"minModels":      ut.MinModels,
		"maxModels":      ut.MaxModels,
	})
	if err != nil {
		return err
	}
	return result.Err()
}

func (r *unitTypesRepo) Get(ctx context.Context, id string) (*types.UnitType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetTx(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, types.NewNotFoundError("unit type not found Get")
	}

	return res.(*types.UnitType), nil
}

func (r *unitTypesRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, id string) (*types.UnitType, error) {
	uts, _, err := r.FindTx(ctx, tx, &FindUnitTypesOpts{UnitTypeIDs: []string{id}, IncludeUnitTypeOptions: true})
	if err != nil {
		return nil, err
	} else if len(uts) == 0 {
		return nil, types.NewNotFoundError("unit type not found GetTx")
	}
	return uts[0], nil
}

func (r *unitTypesRepo) Find(ctx context.Context, opts *FindUnitTypesOpts) ([]*types.UnitType, int64, error) {
	var count int64
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		uts, c, e := r.FindTx(ctx, tx, opts)
		count = c
		return uts, e
	})
	if err != nil {
		return nil, 0, err
	} else if res == nil {
		return []*types.UnitType{}, 0, nil
	}

	return res.([]*types.UnitType), count, nil
}

func (r *unitTypesRepo) GetNamesByArmyTypeID(ctx context.Context, armyTypeID string) ([]*types.UnitType, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		cmd := `
			MATCH (at:ArmyType{id: $army_type_id})<-[:BELONGS_TO]-(ut:UnitType)
			MATCH (ut)-[:IS_COMPOSITION_TYPE]->(ct:CompositionType)
			RETURN ut, ct
			ORDER BY ct.position, ut.created_at;
		`
		params := map[string]any{
			"army_type_id": armyTypeID,
		}
		result, e := tx.Run(ctx, cmd, params)
		if e != nil {
			return nil, e
		}

		uts := []*types.UnitType{}
		for result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, types.NewInternalServerError("unable to convert unit type database object")
			}
			ut := types.UnitTypeFromNode(node)

			node, ok = result.Record().Values[1].(dbtype.Node)
			if !ok {
				return nil, types.NewInternalServerError("unable to convert unit type database object")
			}
			ct := types.CompositionTypeFromNode(node)
			ut.CompositionTypeName = ct.Name

			uts = append(uts, ut)
		}

		return uts, nil
	})
	if err != nil {
		return nil, err
	}

	return res.([]*types.UnitType), nil
}

func (r *unitTypesRepo) FindTx(ctx context.Context, tx neo4j.ManagedTransaction, opts *FindUnitTypesOpts) (
	[]*types.UnitType, int64, error,
) {
	if opts == nil {
		opts = &FindUnitTypesOpts{Limit: 25}
	}

	var limitQry string
	var offsetQry string

	if opts.Limit > 0 {
		limitQry = fmt.Sprintf("LIMIT %d", opts.Limit)
	}

	if opts.Offset > 0 {
		offsetQry = fmt.Sprintf("SKIP %d", opts.Offset)
	}

	params := make(map[string]any)

	atQry := "(ut:UnitType)"
	if len(opts.ArmyTypeID) > 0 {
		atQry = `(at:ArmyType{id: $army_type_id})<-[:BELONGS_TO]-(ut:UnitType{army_type_id: at.id})`
		params["army_type_id"] = opts.ArmyTypeID
	}

	comp := "WHERE"

	var idQry string
	if len(opts.UnitTypeIDs) > 0 {
		idQry = fmt.Sprintf("%s ut.id IN $ids", comp)
		params["ids"] = opts.UnitTypeIDs
		comp = "AND"
	}

	var nQry string
	if len(opts.Name) > 0 {
		nQry = fmt.Sprintf("%s ut.name = $name", comp)
		params["name"] = opts.Name
		comp = "AND"
	}
	if len(opts.Names) > 0 {
		nQry += fmt.Sprintf("%s ut.name IN $names", comp)
		params["names"] = opts.Names
		comp = "AND"
	}

	childAry := comp
	comp = "AND"
	if !opts.IsChildRequest {
		childAry += " NOT"
	}
	childAry += " EXISTS((ut)-[:HAS_PARENT]->(:UnitType"
	var insideChildAry string
	if len(opts.ArmyTypeID) > 0 {
		if len(insideChildAry) > 0 {
			insideChildAry += ","
		} else {
			insideChildAry += "{"
		}
		insideChildAry += "army_type_id: $army_type_id"
	}

	if len(opts.ParentUnitTypeID) > 0 {
		params["parent_unit_type_id"] = opts.ParentUnitTypeID
		if len(insideChildAry) > 0 {
			insideChildAry += ","
		} else {
			insideChildAry += "{"
		}
		insideChildAry += "id: $parent_unit_type_id"
	}

	if len(insideChildAry) > 0 {
		insideChildAry += "}"
	}

	childAry += insideChildAry + "))"

	cmd := fmt.Sprintf(`
		MATCH %s-[:HAS_STATISTIC]->(us:UnitStatistic)-[:IS_STATISTIC]->(s:Statistic)
		%s %s %s
		MATCH (ut)-[:IS_TROOP_TYPE]->(tt:TroopType)
		MATCH (ut)-[:IS_COMPOSITION_TYPE]->(ct:CompositionType)
		UNWIND (us) as unitStats
		UNWIND (s) as stats
		RETURN ut, tt, ct, collect(stats), collect(unitStats)
		ORDER BY ut.name
		%s %s;
	`, atQry, idQry, nQry, childAry, offsetQry, limitQry)

	if opts.Debug {
		fmt.Println("??????? UnitType Find Query")
		fmt.Println(cmd)
		spew.Dump(params)
	}

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, 0, err
	}

	uts := []*types.UnitType{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type")
		}

		ut := types.UnitTypeFromNode(node)
		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert trool type node from database type")
		}
		tt := types.TroopTypeFromNode(node)
		ut.TroopTypeID = tt.ID
		ut.TroopTypeName = tt.Name

		node, ok = result.Record().Values[2].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type composition type")
		}
		ct := types.CompositionTypeFromNode(node)
		ut.CompositionTypeID = ct.ID
		ut.CompositionTypeName = ct.Name

		stats := []*types.Statistic{}
		node2, ok := result.Record().Values[3].([]interface{})
		if ok {
			for _, n := range node2 {
				stats = append(stats, types.StatisticFromNode(n.(dbtype.Node)))
			}
		}
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].Position < stats[j].Position
		})

		node3, ok := result.Record().Values[4].([]interface{})
		if ok {
			// TODO: VERY inefficient... need to do this in neo4j at some point
			// do an ORDER BY with multiple queries in NEO4j at some point
			for _, st := range stats {
				for _, n := range node3 {
					us := types.UnitStatisticFromNode(n.(dbtype.Node))
					if us.StatisticID == st.ID {
						us.Statistic = *st
						ut.Statistics = append(ut.Statistics, us)
						continue
					}
				}
			}
		}

		if !opts.IsChildRequest {
			ut.Children, err = r.getChildUnitTypes(ctx, tx, ut.ID)
			if err != nil {
				return nil, 0, err
			}
		}

		if opts.IncludeUnitTypeOptions {
			ut.Options, err = r.getUnitOptionsFromUnitTypeByID(ctx, tx, ut.ID)
			if err != nil {
				return nil, 0, err
			}
		}

		uts = append(uts, ut)
	}

	var count int64

	if len(opts.ArmyTypeID) > 0 {
		result, err = tx.Run(ctx, `
			MATCH (at:ArmyType{ id: $army_type_id })<-[:BELONGS_TO]-(n:UnitOptionType)
			RETURN count(n) as count
		`, map[string]any{"army_type_id": opts.ArmyTypeID})
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
	} else {
		result, err = tx.Run(ctx, `
			MATCH (n:UnitOptionType)
			RETURN count(n) as count
		`, map[string]any{})
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
	}

	return uts, count, nil
}

func (r *unitTypesRepo) getUnitOptionsFromUnitTypeByID(
	ctx context.Context, tx neo4j.ManagedTransaction, unitTypeID string,
) ([]*types.UnitTypeOption, error) {
	itRepo := NewItemTypesRepo(r.db)

	params := map[string]any{
		"unit_type_id": unitTypeID,
	}

	cmd := `
		MATCH (ut:UnitType{ id: $unit_type_id })
		MATCH (ut)-[:IS_UNIT_OPTION]->(nuo:UnitOption)
		MATCH (nuo)-[:IS_UNIT_OPTION_TYPE]->(uot:UnitOptionType)

		UNWIND (nuo) AS unitOptions
		UNWIND (uot) AS unitOptionTypes

		RETURN ut, collect(unitOptions), collect(uot)
	`

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	utos := []*types.UnitTypeOption{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert unit option unit type")
		}
		ut := types.UnitTypeFromNode(node)

		nodes, ok := result.Record().Values[2].([]interface{})
		if !ok {
			return nil, errors.New("unable to convert unit option type slice")
		}
		uotMap := map[string]*types.UnitOptionType{}
		for _, nd := range nodes {
			node, ok := nd.(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database unit option type")
			}
			uot := types.UnitOptionTypeFromNode(node)
			uotMap[uot.ID] = uot
		}

		nodes, ok = result.Record().Values[1].([]interface{})
		if !ok {
			return nil, errors.New("unable to convert unit option slice")
		}
		for _, nd := range nodes {
			node, ok := nd.(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database unit option")
			}
			uto := types.UnitTypeOptionFromNode(node)
			uto.UnitTypeID = ut.ID
			uto.UnitTypeName = ut.Name

			uot, ok := uotMap[uto.UnitOptionTypeID]
			if !ok {
				return nil, errors.New("unable to get unit option type from unit type option")
			}
			uto.UnitOptionTypeName = uot.Name

			result2, err := tx.Run(ctx, `
			MATCH (uo:UnitOption{ id: $unit_type_option_id })
			MATCH (it:ItemType)<-[:IS_ITEM_TYPE]-(i:Item)<-[rel:IS_OPTION_ITEM]-(uo)
			WITH i
			ORDER BY it.position, i.points DESC
			UNWIND (i) AS items
			RETURN collect(items)
			`, map[string]any{
				"unit_type_option_id": uto.ID,
			})
			if err != nil {
				return nil, err
			}

			if result2.Next(ctx) {
				nodes, ok := result2.Record().Values[0].([]interface{})
				if !ok {
					return nil, errors.New("unable to convert database unit option items")
				}

				for _, nd := range nodes {
					node, ok := nd.(dbtype.Node)
					if !ok {
						return nil, errors.New("unable to get unit unit type option item")
					}
					it := types.ItemFromNode(node)

					its, _, err := itRepo.Find(ctx, &FindItemTypeOpts{
						GameID: it.GameID,
						IDs:    []string{it.ItemTypeID},
					})
					if err != nil {
						return nil, err
					} else if len(its) > 0 {
						it.ItemTypeName = its[0].Name
					}

					uto.Items = append(uto.Items, it)
				}
			}

			utos = append(utos, uto)
		}
	}

	return utos, nil
}

func (r *unitTypesRepo) FindOrCreate(ctx context.Context, ut types.CreateUnitType) (*types.UnitType, error) {
	existingUnitType, err := r.getByName(ctx, ut)
	if types.IsNotFoundError(err) {
		// Important to not "GET" here because child items will fail the query
		// Possible TODO: investigate way to query without the parent thing?
		// 				  	not an easy way to tell if the query is a child unless
		//					inspecting CreateUnitType
		return r.create(ctx, ut)
	} else if err != nil {
		return nil, err
	}

	return existingUnitType, nil
}

func (r *unitTypesRepo) EnsureChildUnitTypeExists(ctx context.Context, ncut types.CreateChildUnitType) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.EnsureChildUnitTypeExistsTx(ctx, tx, ncut)
	})
	return err
}

func (r *unitTypesRepo) EnsureChildUnitTypeExistsTx(ctx context.Context, tx neo4j.ManagedTransaction, ncut types.CreateChildUnitType) error {
	if _, err := r.getUnitTypeChildByName(ctx, tx, ncut.UnitTypeID, ncut.Name); err != nil {
		if types.IsNotFoundError(err) {
			// Get the parent for the other data
			parent, err := r.GetTx(ctx, tx, ncut.UnitTypeID)
			if err != nil {
				return err
			}

			if _, err = r.create(ctx, types.CreateUnitType{
				Name:              ncut.Name,
				GameID:            parent.GameID,
				ArmyTypeID:        parent.ArmyTypeID,
				TroopTypeID:       parent.TroopTypeID,
				CompositionTypeID: parent.CompositionTypeID,
				Statistics:        ncut.Statistics,
				UnitTypeID:        parent.ID,
			}); err != nil {
				return err
			}

			return nil
		}
		return err
	}

	return nil
}

func (r *unitTypesRepo) getByName(ctx context.Context, ut types.CreateUnitType) (*types.UnitType, error) {
	uts, _, err := r.Find(ctx, &FindUnitTypesOpts{Name: ut.Name, ArmyTypeID: ut.ArmyTypeID})
	if err != nil {
		return nil, err
	} else if len(uts) == 0 {
		return nil, types.NewNotFoundError("unit type not found getByName")
	}
	return uts[0], nil
}

func (r *unitTypesRepo) getUnitTypeChildByName(ctx context.Context, tx neo4j.ManagedTransaction, unitTypeID, name string) (*types.UnitType, error) {
	uts, _, err := r.FindTx(ctx, tx, &FindUnitTypesOpts{
		UnitTypeIDs:    []string{unitTypeID},
		Name:           name,
		IsChildRequest: true,
		Limit:          1,
	})
	if err != nil {
		return nil, err
	} else if len(uts) == 0 {
		return nil, types.NewNotFoundError("unable to find unit type child")
	}

	return uts[0], nil
}

func (r *unitTypesRepo) getChildUnitTypes(ctx context.Context, tx neo4j.ManagedTransaction, unitTypeID string) ([]*types.UnitType, error) {
	uts, _, err := r.FindTx(ctx, tx, &FindUnitTypesOpts{
		ParentUnitTypeID: unitTypeID,
		IsChildRequest:   true,
	})
	if err != nil {
		return nil, err
	}

	return uts, nil
}

func (r *unitTypesRepo) create(ctx context.Context, ut types.CreateUnitType) (*types.UnitType, error) {
	if err := types.Validate(ut); err != nil {
		return nil, err
	} else if len(ut.Statistics) == 0 {
		return nil, errors.New("unit types statistics are required for a new unit type")
	}

	currentTimestamp := time.Now().UTC().Unix()

	newUt, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"name":                ut.Name,
			"game_id":             ut.GameID,
			"army_type_id":        ut.ArmyTypeID,
			"troop_type_id":       ut.TroopTypeID,
			"composition_type_id": ut.CompositionTypeID,
			"points_per_model":    ut.PointsPerModel,
			"min_models":          ut.MinModels,
			"max_models":          ut.MaxModels,
			"created_at":          currentTimestamp,
			"updated_at":          currentTimestamp,
		}

		pQry := ``
		relQry := "CREATE (ut)-[:BELONGS_TO]->(at)"
		if len(ut.UnitTypeID) > 0 {
			pQry = "MATCH (put:UnitType{ id: $parentUnitTypeId })"
			relQry = `CREATE (put)<-[:HAS_PARENT]-(ut)`
			params["parentUnitTypeId"] = ut.UnitTypeID
		}

		cmd := fmt.Sprintf(`
			MATCH (g:Game{id: $game_id})
			MATCH (at:ArmyType{id: $army_type_id, game_id: $game_id})
			MATCH (tt:TroopType{id: $troop_type_id, game_id: $game_id})
			MATCH (ct:CompositionType{id: $composition_type_id, game_id: $game_id})
			%s
			MERGE (ut:UnitType{
				name: 					$name
				,game_id: 				g.id
				,army_type_id:			at.id
				,troop_type_id: 		tt.id
				,composition_type_id: 	ct.id
				,points_per_model:		$points_per_model
				,min_models:			$min_models
				,max_models:			$max_models
			})
			ON CREATE
				SET ut.created_at 	= $created_at
				,ut.id 				= apoc.create.uuid()
			ON MATCH
				SET ut.updated_at 	= $updated_at
			%s
			MERGE (ut)-[:BELONGS_TO]->(g)
			MERGE (tt)<-[:IS_TROOP_TYPE]-(ut)
			MERGE (ct)<-[:IS_COMPOSITION_TYPE]-(ut)
			RETURN ut, tt, ct;
		`, pQry, relQry)

		result, err := tx.Run(ctx, cmd, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			newUt := types.UnitTypeFromNode(node)

			node, ok = result.Record().Values[1].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type troop type")
			}
			tt := types.TroopTypeFromNode(node)
			newUt.TroopTypeID = tt.ID
			newUt.TroopTypeName = tt.Name

			node, ok = result.Record().Values[2].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type composition type")
			}
			ct := types.CompositionTypeFromNode(node)
			newUt.CompositionTypeID = ct.ID
			newUt.CompositionTypeName = ct.Name

			for _, uts := range ut.Statistics {
				usParams := map[string]any{
					"unit_type_id": newUt.ID,
					"value":        uts.Value,
					"display":      uts.Display,
					"created_at":   currentTimestamp,
					"updated_at":   currentTimestamp,
				}

				result2, err := tx.Run(ctx, `
					MATCH (s:Statistic{ display: $display })
					MATCH (ut:UnitType{ id: $unit_type_id })
					MERGE (us:UnitStatistic{
						,value: 		$value
						,unit_type_id: 	ut.id
						,statistic_id:	s.id
					})
					ON CREATE
  						SET us.created_at = $created_at,
						  us.id = apoc.create.uuid()
					ON MATCH
						SET us.updated_at = $updated_at
					MERGE (s)<-[:IS_STATISTIC]-(us)<-[:HAS_STATISTIC]-(ut)
					RETURN us,s;
				`, usParams)
				if err != nil {
					return nil, err
				}

				if result2.Next(ctx) {
					node, ok := result2.Record().Values[0].(dbtype.Node)
					if !ok {
						return nil, errors.New("unable to convert database type unit statistics")
					}
					newUs := types.UnitStatisticFromNode(node)

					node, ok = result2.Record().Values[1].(dbtype.Node)
					if !ok {
						return nil, errors.New("unable to convert database type unit statistics")
					}
					newStat := types.StatisticFromNode(node)
					newUs.Statistic = *newStat

					var found bool
					for _, exStat := range newUt.Statistics {
						if exStat.ID == newUs.ID {
							found = true
						}
					}

					if !found {
						newUt.Statistics = append(newUt.Statistics, newUs)
					}
				}
			}

			if len(newUt.Statistics) == 0 {
				return nil, errors.New("statistics are required for unit type - likely no match for statistics found")
			}

			for idx, opt := range ut.UnitOptions {
				params3 := map[string]any{
					"unit_type_id":        newUt.ID,
					"unit_option_type_id": opt.UnitOptionTypeID,
					"position":            idx + 1,
					"txt":                 opt.Txt,
					"points":              opt.Points,
					"per_model":           opt.PerModel,
					"max_pts":             opt.MaxPoints,
					"created_at":          currentTimestamp,
					"updated_at":          currentTimestamp,
				}

				cmd3 := `
				MATCH (ut:UnitType{ id: $unit_type_id })
				MATCH (uot:UnitOptionType{ id: $unit_option_type_id })
				MERGE (nuo:UnitOption{
					unit_type_id:			$unit_type_id
					,unit_option_type_id:	$unit_option_type_id
					,position:				$position
					,txt:					$txt
					,points:				$points
					,per_model:				$per_model
					,max_pts: 				$max_pts
				})
				ON CREATE
					SET nuo.created_at 	= $created_at
					,nuo.id 				= apoc.create.uuid()
				ON MATCH
					SET nuo.updated_at 	= $updated_at
				MERGE (ut)-[:IS_UNIT_OPTION]->(nuo)
				MERGE (nuo)-[:IS_UNIT_OPTION_TYPE]->(uot)
				RETURN nuo, uot;
				`

				result3, err := tx.Run(ctx, cmd3, params3)
				if err != nil {
					return nil, err
				}

				if result3.Next(ctx) {
					node, ok := result3.Record().Values[0].(dbtype.Node)
					if !ok {
						return nil, errors.New("unable to convert database type")
					}
					newUto := types.UnitTypeOptionFromNode(node)

					node, ok = result3.Record().Values[1].(dbtype.Node)
					if !ok {
						return nil, errors.New("unable to convert database type")
					}
					newUot := types.UnitOptionTypeFromNode(node)
					newUto.UnitOptionTypeName = newUot.ID

					for _, uoit := range opt.Items {
						// Create the link to the unit options
						params4 := map[string]any{
							"item_id":        uoit.ID,
							"unit_option_id": newUto.ID,
						}

						cmd4 := `
						MATCH (i:Item{ id: $item_id })
						MATCH (nuo:UnitOption{ id: $unit_option_id })
						MERGE (i)<-[rel:IS_OPTION_ITEM]-(nuo)
						RETURN rel;
						`

						if _, err := tx.Run(ctx, cmd4, params4); err != nil {
							return nil, err
						}
					}

					newUto.Items = append(newUto.Items, opt.Items...)
					newUt.Options = append(newUt.Options, newUto)
				}
			}

			return newUt, nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	return newUt.(*types.UnitType), nil
}
