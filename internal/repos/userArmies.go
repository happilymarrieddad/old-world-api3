package repos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type FindUserArmyOpts struct {
	IDs    []string
	Limit  int
	Offset int
	Debug  bool
}

//go:generate mockgen -source=./userArmies.go -destination=./mocks/UserArmiesRepo.go -package=mock_repos UserArmiesRepo
type UserArmiesRepo interface {
	FindUserUnitIDsByUnitTypeID(ctx context.Context, unitTypeID string) ([]*types.UserArmyUnit, error)
	FindUserUnitIDsByUnitTypeIDTx(ctx context.Context, tx neo4j.ManagedTransaction, unitTypeID string) ([]*types.UserArmyUnit, error)
	Get(ctx context.Context, userID, userArmyID string) (*types.UserArmy, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, userID, userArmyID string) (*types.UserArmy, error)
	Find(ctx context.Context, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, int64, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, int64, error)
	Create(ctx context.Context, nua types.CreateUserArmy) (*types.UserArmy, error)
	CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, nua types.CreateUserArmy) (*types.UserArmy, error)
	Update(ctx context.Context, obj types.UpdateUserArmy) error
	UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, obj types.UpdateUserArmy) error
	GetUnit(ctx context.Context, userArmyUnitID string) (*types.UserArmyUnit, error)
	GetUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyUnitID string) (*types.UserArmyUnit, error)
	AddUnits(ctx context.Context, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error)
	AddUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error)
	RemoveUnits(ctx context.Context, userArmyID string, userArmyUnitIDs ...string) error
	RemoveUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, userArmyUnitIDs ...string) error
	UpdateUnit(ctx context.Context, opts types.UpdateUserArmyUnit) error
	UpdateUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, opts types.UpdateUserArmyUnit) error
}

func NewUserArmiesRepo(db neo4j.DriverWithContext) UserArmiesRepo {
	return &userArmyRepo{db}
}

type userArmyRepo struct {
	db neo4j.DriverWithContext
}

func (r *userArmyRepo) FindUserUnitIDsByUnitTypeID(ctx context.Context, unitTypeID string) ([]*types.UserArmyUnit, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.FindUserUnitIDsByUnitTypeIDTx(ctx, tx, unitTypeID)
	})

	if err != nil {
		return nil, err
	}

	return res.([]*types.UserArmyUnit), nil
}

func (r *userArmyRepo) FindUserUnitIDsByUnitTypeIDTx(ctx context.Context, tx neo4j.ManagedTransaction, unitTypeID string) ([]*types.UserArmyUnit, error) {
	cmd := `
		MATCH (ut:UnitType{ id: $unitTypeId })
		MATCH (uau:UserArmyUnit)-[:IS_UNIT_TYPE]->(ut)
		RETURN uau
		`
	params := map[string]any{
		"unitTypeId": unitTypeID,
	}

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	uaus := []*types.UserArmyUnit{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user army")
		}
		uua := types.UserArmyUnitFromNode(node)

		uau, err := r.GetUnitTx(ctx, tx, uua.ID)
		if err != nil {
			return nil, err
		}

		uaus = append(uaus, uau)
	}

	return uaus, nil
}

func (r *userArmyRepo) Get(ctx context.Context, userID, userArmyID string) (*types.UserArmy, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetTx(ctx, tx, userID, userArmyID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.UserArmy), nil
}

func (r *userArmyRepo) GetTx(ctx context.Context, tx neo4j.ManagedTransaction, userID, userArmyID string) (*types.UserArmy, error) {
	uas, _, err := r.FindTx(ctx, tx, userID, &FindUserArmyOpts{IDs: []string{userArmyID}, Limit: 1})
	if err != nil {
		return nil, err
	} else if len(uas) == 0 {
		return nil, types.NewNotFoundError("unable to get user army")
	}
	return uas[0], nil
}

func (r *userArmyRepo) Find(ctx context.Context, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, int64, error) {
	var count int64
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		res, c, e := r.FindTx(ctx, tx, userID, opts)
		count = c
		return res, e
	})
	if err != nil {
		return nil, 0, err
	}
	return res.([]*types.UserArmy), count, nil
}

func (r *userArmyRepo) FindTx(
	ctx context.Context, tx neo4j.ManagedTransaction, userID string, opts *FindUserArmyOpts,
) ([]*types.UserArmy, int64, error) {
	if opts == nil {
		opts = &FindUserArmyOpts{}
	}

	params := map[string]any{
		"user_id": userID,
	}

	var limitOffsetQry string
	if opts.Offset > 0 {
		limitOffsetQry += fmt.Sprintf(" SKIP %d", opts.Offset)
	}
	if opts.Limit > 0 {
		limitOffsetQry += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	comp := `WHERE`
	var whrQry string
	if len(opts.IDs) > 0 {
		whrQry += fmt.Sprintf(" %s ua.id IN $user_army_ids", comp)
		params["user_army_ids"] = opts.IDs
	}

	cmd := fmt.Sprintf(`
	MATCH (u:User{ id: $user_id })
	MATCH (g:Game)<-[:BELONGS_TO]-(ua:UserArmy)-[:BELONGS_TO]->(u)
	MATCH (ua)-[:BELONGS_TO]->(at:ArmyType)
	%s
	RETURN ua, g, u, at
	ORDER BY ua.name
	%s;
	`, whrQry, limitOffsetQry)

	if opts.Debug {
		fmt.Println("??????? UserArmies Find Query")
		fmt.Println(cmd)
		spew.Dump(params)
	}

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, 0, err
	}

	mpUnits := make(map[string][]*types.UserArmyUnit)
	innerResult, err := tx.Run(ctx, `
		MATCH (u:User{ id: $user_id })<-[:BELONGS_TO]-(ua:UserArmy)
		MATCH (uau:UserArmyUnit)-[:BELONGS_TO]->(ua)
		MATCH (uau)-[:IS_UNIT_TYPE]->(ut)-[:IS_COMPOSITION_TYPE]->(ct:CompositionType)
		RETURN uau
		ORDER BY ct.position, ut.created_at;
	`, map[string]any{"user_id": userID})
	if err != nil {
		return nil, 0, err
	}

	for innerResult.Next(ctx) {
		node, ok := innerResult.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type user army")
		}
		uua := types.UserArmyUnitFromNode(node)

		uua, err = r.GetUnitTx(ctx, tx, uua.ID)
		if err != nil {
			return nil, 0, errors.New("unable to get full unit data")
		}

		mpUnits[uua.UserArmyID] = append(mpUnits[uua.UserArmyID], uua)
	}

	uas := []*types.UserArmy{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type user army")
		}
		ua := types.UserArmyFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type game")
		}
		gm := types.GameFromNode(node)
		ua.GameName = gm.Name

		node, ok = result.Record().Values[2].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type user")
		}
		usr := types.UserFromNode(node)
		ua.UserFirstName = usr.FirstName
		ua.UserLastName = usr.LastName
		ua.UserEmail = usr.Email

		node, ok = result.Record().Values[3].(dbtype.Node)
		if !ok {
			return nil, 0, errors.New("unable to convert database type army type")
		}
		at := types.ArmyTypeFromNode(node)
		ua.ArmyTypeName = at.Name

		// Fill out the units
		unts, exists := mpUnits[ua.ID]
		if exists && len(unts) > 0 {
			ua.Units = unts
		}

		uas = append(uas, ua)
	}

	var count int64
	result, err = tx.Run(ctx, `
		MATCH (u:User{ id: $user_id })
		MATCH (ua:UserArmy)-[:BELONGS_TO]->(u)
		RETURN count(ua) as count
	`, map[string]any{"user_id": userID})
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

	return uas, count, nil
}

func (r *userArmyRepo) Create(ctx context.Context, nua types.CreateUserArmy) (*types.UserArmy, error) {
	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.CreateTx(ctx, tx, nua)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.UserArmy), nil
}

func (r *userArmyRepo) CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, nua types.CreateUserArmy) (*types.UserArmy, error) {
	if err := types.Validate(nua); err != nil {
		return nil, err
	}

	currentTimestamp := time.Now().UTC().Unix()
	params := map[string]any{
		"name":         nua.Name,
		"game_id":      nua.GameID,
		"user_id":      nua.UserID,
		"army_type_id": nua.ArmyTypeID,
		"points":       nua.Points,
		"created_at":   currentTimestamp,
		"updated_at":   currentTimestamp,
	}
	cmd := `
	MATCH (g:Game{ id: $game_id })
	MATCH (u:User{ id: $user_id })
	MATCH (at:ArmyType{ id: $army_type_id })

	MERGE (ua:UserArmy{
		name: 			$name
		,game_id:		$game_id
		,user_id:		$user_id
		,army_type_id:	$army_type_id
	})
	ON CREATE
		SET ua.created_at = $created_at
			,ua.id = apoc.create.uuid()
			,ua.points = $points
	ON MATCH
		set ua.updated_at = $updated_at
			,ua.points = $points
	MERGE (ua)-[:BELONGS_TO]->(g)
	MERGE (ua)-[:BELONGS_TO]->(u)
	MERGE (ua)-[:BELONGS_TO]->(at)
	RETURN ua, g, u, at;
	`

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user army")
		}
		newUserArmy := types.UserArmyFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type game")
		}
		gm := types.GameFromNode(node)
		newUserArmy.GameName = gm.Name

		node, ok = result.Record().Values[2].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user")
		}
		usr := types.UserFromNode(node)
		newUserArmy.UserFirstName = usr.FirstName
		newUserArmy.UserLastName = usr.LastName
		newUserArmy.UserEmail = usr.Email

		node, ok = result.Record().Values[3].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type army type")
		}
		at := types.ArmyTypeFromNode(node)
		newUserArmy.ArmyTypeName = at.Name

		return newUserArmy, nil
	}

	return nil, types.NewNotFoundError("unable to create user army")
}

func (r *userArmyRepo) RemoveUnits(ctx context.Context, userArmyID string, userArmyUnitIDs ...string) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.RemoveUnitsTx(ctx, tx, userArmyID, userArmyUnitIDs...)
	})
	return err
}

func (r *userArmyRepo) RemoveUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, userArmyUnitIDs ...string) error {
	params := map[string]any{
		"user_army_id": userArmyID,
		"ids":          userArmyUnitIDs,
	}
	cmd := `
	MATCH (ua:UserArmy{ id: $user_army_id })
	MATCH (uau:UserArmyUnit)-[:BELONGS_TO]->(ua)
	WHERE uau.id IN $ids
	DETACH DELETE (uau)
	`

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return err
	}

	return result.Err()
}

func (r *userArmyRepo) Update(ctx context.Context, obj types.UpdateUserArmy) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.UpdateTx(ctx, tx, obj)
	})
	return err
}

func (r *userArmyRepo) UpdateTx(ctx context.Context, tx neo4j.ManagedTransaction, obj types.UpdateUserArmy) error {
	if obj.Name == nil && obj.Points == nil {
		return nil
	}

	params := map[string]any{
		"user_army_id": obj.ID,
	}
	cmd := `
	MATCH (ua:UserArmy{ id: $user_army_id })
	`

	var setQry string
	if obj.Name != nil {
		setQry += "SET ua.name = $name"
		params["name"] = *obj.Name
	}

	if obj.Points != nil {
		params["points"] = *obj.Points
		if len(setQry) == 0 {
			setQry = "SET ua.points = $points"
		} else {
			setQry += " ,ua.points = $points"
		}
	}
	cmd += " " + setQry + ` RETURN ua;`

	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return err
	}

	return result.Err()
}

func (r *userArmyRepo) AddUnits(ctx context.Context, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error) {
	res, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.AddUnitsTx(ctx, tx, userArmyID, uaus...)
	})
	if err != nil {
		return nil, err
	}
	return res.([]*types.UserArmyUnit), nil
}

func (r *userArmyRepo) AddUnitsTx(
	ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, cuaus ...*types.CreateUserArmyUnit,
) ([]*types.UserArmyUnit, error) {
	unitTypesRepo := NewUnitTypesRepo(r.db)
	uau := []*types.UserArmyUnit{}

	var armyPoints int

	utRepo := NewUnitTypesRepo(r.db)
	currentTimestamp := time.Now().UTC().Unix()
	for _, cuau := range cuaus {
		ut, err := utRepo.Get(ctx, cuau.UnitTypeID)
		if err != nil {
			return nil, err
		}
		// Do not worry about options here because any new army unit won't have any options selected
		cuau.Points = ut.MinModels * ut.PointsPerModel

		armyPoints += cuau.Points

		params := map[string]any{
			"user_army_id": userArmyID,
			"unit_type_id": cuau.UnitTypeID,
			"quantity":     ut.MinModels,
			"points":       cuau.Points,
			"created_at":   currentTimestamp,
			"updated_at":   currentTimestamp,
		}

		cmd := `
			MATCH (ua:UserArmy{ id: $user_army_id })
			MATCH (ut:UnitType{ id: $unit_type_id })

			CREATE (uau:UserArmyUnit{
				id: 					apoc.create.uuid()
				,user_army_id: 			$user_army_id
				,unit_type_id:			$unit_type_id
				,quantity:				$quantity
				,points:				$points
				,created_at:			$created_at
			})

			MERGE (uau)-[:BELONGS_TO]->(ua)
			MERGE (uau)-[:IS_UNIT_TYPE]->(ut)
			RETURN uau, ut, ua;
		`

		result, err := tx.Run(ctx, cmd, params)
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			node, ok := result.Record().Values[0].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type")
			}
			newUau := types.UserArmyUnitFromNode(node)

			node, ok = result.Record().Values[1].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type troop type")
			}
			eut := types.UnitTypeFromNode(node)

			ut, err := unitTypesRepo.GetTx(ctx, tx, eut.ID)
			if err != nil {
				return nil, err
			}
			newUau.UnitType = ut

			for _, opt := range ut.Options {
				params2 := map[string]any{
					"user_army_unit_id": newUau.ID,
					"unit_option_id":    opt.ID,
					"is_selected":       false,
					"index_selected":    "",
					"ids_selected":      []string{},
					"qty_selected":      0,
					"created_at":        currentTimestamp,
				}
				switch opt.UnitOptionTypeName {
				case "Single", "One Of", "Many Of", "Many To":
					params2["is_selected"] = false
					cmd2 := `
						MATCH (uau:UserArmyUnit{ id: $user_army_unit_id })
						MATCH (uto:UnitOption{ id: $unit_option_id })

						CREATE (uauo:UserArmyUnitOptionValue{
							id: 					apoc.create.uuid()
							,user_army_unit_id: 	$user_army_unit_id
							,unit_option_id:		$unit_option_id
							,is_selected:			$is_selected
							,index_selected:		$index_selected
							,ids_selected:			$ids_selected
							,qty_selected:			$qty_selected
							,created_at:			$created_at
						})

						MERGE (uauo)-[:IS_USER_ARMY_UNIT_OPTION_VALUE]->(uau)
						MERGE (uauo)<-[:IS_UNIT_TYPE_OPTION]-(uto)

						RETURN uauo;
					`
					result2, err := tx.Run(ctx, cmd2, params2)
					if err != nil {
						return nil, err
					}

					if result2.Next(ctx) {
						node, ok := result2.Record().Values[0].(dbtype.Node)
						if !ok {
							return nil, types.NewInternalServerError("unable to get user army unit option val")
						}
						uaov := types.UserArmyUnitOptionValueFromNode(node)
						uaov.UserArmyUnitName = newUau.UnitType.Name
						uaov.UnitOptionName = opt.UnitOptionTypeName
						uaov.UnitOption = opt

						newUau.OptionValues = append(newUau.OptionValues, uaov)
					}
				default:
					return nil, types.NewInternalServerError(fmt.Sprintf("invalid unit option type which should never happen err: %s", opt.UnitOptionTypeName))
				}
			}

			node, ok = result.Record().Values[2].(dbtype.Node)
			if !ok {
				return nil, errors.New("unable to convert database type troop type")
			}
			ua := types.UserArmyFromNode(node)
			newUau.UserArmyName = ua.Name

			uau = append(uau, newUau)
		}
	}

	return uau, nil
}

func (r *userArmyRepo) GetUnit(ctx context.Context, userArmyUnitID string) (*types.UserArmyUnit, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.GetUnitTx(ctx, tx, userArmyUnitID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*types.UserArmyUnit), nil
}

func (r *userArmyRepo) GetUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyUnitID string) (*types.UserArmyUnit, error) {
	params := map[string]any{
		"user_army_unit_id": userArmyUnitID,
	}
	cmd := `
	MATCH (uau:UserArmyUnit{ id: $user_army_unit_id })
	MATCH (uau)-[:BELONGS_TO]->(ua:UserArmy)
	MATCH (uau)-[:IS_UNIT_TYPE]->(ut)-[:IS_COMPOSITION_TYPE]->(ct:CompositionType)
	RETURN uau, ua, ct
	ORDER BY ct.position, ut.created_at
	`
	result, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user unit type")
		}
		uua := types.UserArmyUnitFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user unit army type")
		}
		ua := types.UserArmyFromNode(node)
		uua.UserArmyName = ua.Name

		unitTypeRepo := NewUnitTypesRepo(r.db)
		ut, err := unitTypeRepo.GetTx(ctx, tx, uua.UnitTypeID)
		if err != nil {
			return nil, err
		}
		uua.UnitType = ut

		result2, err2 := tx.Run(ctx, `
			MATCH (uau:UserArmyUnit{ id: $id })

			OPTIONAL MATCH (uau)<-[:IS_USER_ARMY_UNIT_OPTION_VALUE]-(uauo:UserArmyUnitOptionValue)
			OPTIONAL MATCH (uauo)<-[:IS_UNIT_TYPE_OPTION]-(uto:UnitOption)
			UNWIND (uauo) AS optionVals
			UNWIND (uto) AS unitTypeOptions
			RETURN collect(unitTypeOptions), collect(optionVals)
		`, map[string]any{
			"id": uua.ID,
		})
		if err2 != nil {
			return nil, err2
		}

		for result2.Next(ctx) {
			nodes, ok := result2.Record().Values[0].([]interface{})
			if !ok {
				return nil, errors.New("unable to convert database type user unit unit option types slice")
			}
			utoMap := make(map[string]*types.UnitTypeOption)
			for _, node := range nodes {
				nd, ok := node.(dbtype.Node)
				if !ok {
					return nil, errors.New("unable to get unit type option from database")
				}
				uot := types.UnitTypeOptionFromNode(nd)
				utoMap[uot.ID] = uot
			}

			nodes, ok = result2.Record().Values[1].([]interface{})
			if !ok {
				return nil, errors.New("unable to convert database type user unit option possible values slice")
			}
			for _, node := range nodes {
				nd, ok := node.(dbtype.Node)
				if !ok {
					return nil, errors.New("unable to get unit option possible value from database")
				}
				uaov := types.UserArmyUnitOptionValueFromNode(nd)
				uaov.UserArmyUnitName = uua.UnitType.Name
				for _, uo := range ut.Options {
					if uaov.UnitOptionID == uo.ID {
						uaov.UnitOptionName = uo.UnitOptionTypeName
						uaov.UnitOption = uo
					}
				}

				uua.OptionValues = append(uua.OptionValues, uaov)
			}
		}

		return uua, nil
	}

	log.Println("unable to get user army unit '" + userArmyUnitID + "'")
	return nil, types.NewNotFoundError("unable to get user army unit")
}

func (r *userArmyRepo) UpdateUnit(ctx context.Context, opts types.UpdateUserArmyUnit) error {
	_, err := db.WriteData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.UpdateUnitTx(ctx, tx, opts)
	})
	return err
}

func (r *userArmyRepo) UpdateUnitTx(ctx context.Context, tx neo4j.ManagedTransaction, opts types.UpdateUserArmyUnit) error {
	if err := types.Validate(opts); err != nil {
		return err
	}

	existingUserUnit, err := r.GetUnitTx(ctx, tx, opts.ID)
	if err != nil {
		return err
	}

	// |||||||| Validations |||||||||
	// Validation check for Quantity
	if opts.Qty != nil {
		qty := *opts.Qty
		if existingUserUnit.UnitType.MinModels > qty {
			qty = existingUserUnit.UnitType.MinModels
		}
		if existingUserUnit.UnitType.MaxModels < qty {
			qty = existingUserUnit.UnitType.MaxModels
		}
		existingUserUnit.Quantity = qty
	}
	points := existingUserUnit.Quantity * existingUserUnit.UnitType.PointsPerModel
	if opts.Points != nil {
		points = *opts.Points
	}
	// |||||||| END Validations |||||||||

	params := map[string]any{
		"id":     opts.ID,
		"qty":    existingUserUnit.Quantity,
		"points": points,
	}
	cmd := `
		MATCH (uau:UserArmyUnit{ id: $id })
		SET
			uau.quantity = 	$qty
			,uau.points = 	$points
		RETURN uau;
	`

	results, err := tx.Run(ctx, cmd, params)
	if err != nil {
		return err
	}

	for _, uto := range opts.OptionValues {
		params2 := map[string]any{
			"id":             uto.ID,
			"is_selected":    uto.IsSelected,
			"index_selected": uto.IndexSelected,
			"ids_selected":   uto.IDsSelected,
			"qty_selected":   uto.QtySelected,
		}
		cmd2 := `
			MATCH (uauo:UserArmyUnitOptionValue{ id: $id })
			SET
				uauo.is_selected =			$is_selected
				,uauo.index_selected =		$index_selected
				,uauo.ids_selected =		$ids_selected
				,uauo.qty_selected =		$qty_selected
			RETURN uauo;
		`

		if _, err := tx.Run(ctx, cmd2, params2); err != nil {
			return err
		}
	}

	return results.Err()
}
