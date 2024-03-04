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

type FindUserArmyOpts struct {
	IDs    []string
	Limit  int
	Offset int
	Debug  bool
}

//go:generate mockgen -source=./userArmies.go -destination=./mocks/UserArmiesRepo.go -package=mock_repos UserArmiesRepo
type UserArmiesRepo interface {
	Get(ctx context.Context, userID, userArmyID string) (*types.UserArmy, error)
	GetTx(ctx context.Context, tx neo4j.ManagedTransaction, userID, userArmyID string) (*types.UserArmy, error)
	Find(ctx context.Context, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, error)
	FindTx(ctx context.Context, tx neo4j.ManagedTransaction, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, error)
	Create(ctx context.Context, nua types.CreateUserArmy) (*types.UserArmy, error)
	CreateTx(ctx context.Context, tx neo4j.ManagedTransaction, nua types.CreateUserArmy) (*types.UserArmy, error)
	AddUnits(ctx context.Context, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error)
	AddUnitsTx(ctx context.Context, tx neo4j.ManagedTransaction, userArmyID string, uaus ...*types.CreateUserArmyUnit) ([]*types.UserArmyUnit, error)
}

func NewUserArmiesRepo(db neo4j.DriverWithContext) UserArmiesRepo {
	return &userArmyRepo{db}
}

type userArmyRepo struct {
	db neo4j.DriverWithContext
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
	uas, err := r.FindTx(ctx, tx, userID, &FindUserArmyOpts{IDs: []string{userArmyID}, Limit: 1})
	if err != nil {
		return nil, err
	} else if len(uas) == 0 {
		return nil, types.NewNotFoundError("unable to get user army")
	}
	return uas[0], nil
}

func (r *userArmyRepo) Find(ctx context.Context, userID string, opts *FindUserArmyOpts) ([]*types.UserArmy, error) {
	res, err := db.ReadData(ctx, r.db, func(tx neo4j.ManagedTransaction) (any, error) {
		return r.FindTx(ctx, tx, userID, opts)
	})
	if err != nil {
		return nil, err
	}
	return res.([]*types.UserArmy), nil
}

func (r *userArmyRepo) FindTx(
	ctx context.Context, tx neo4j.ManagedTransaction, userID string, opts *FindUserArmyOpts,
) ([]*types.UserArmy, error) {
	if opts == nil {
		opts = &FindUserArmyOpts{}
	}

	params := map[string]any{
		"user_id": userID,
	}

	var limitOffsetQry string
	if opts.Limit > 0 {
		limitOffsetQry += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	if opts.Offset > 0 {
		limitOffsetQry += fmt.Sprintf(" OFFSET %d", opts.Offset)
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
		return nil, err
	}

	mpUnits := make(map[string][]*types.UserArmyUnit)
	innerResult, err := tx.Run(ctx, `
	MATCH (u:User{ id: $user_id })<-[:BELONGS_TO]-(ua:UserArmy)
	MATCH (uau:UserArmyUnit)-[:BELONGS_TO]->(ua)
	RETURN uau, ua
	`, map[string]any{"user_id": userID})
	if err != nil {
		return nil, err
	}
	unitTypeRepo := NewUnitTypesRepo(r.db)
	for innerResult.Next(ctx) {
		node, ok := innerResult.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user army")
		}
		uua := types.UserArmyUnitFromNode(node)

		node, ok = innerResult.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user army")
		}
		ua := types.UserArmyFromNode(node)
		uua.UserArmyName = ua.Name

		if _, exists := mpUnits[ua.ID]; !exists {
			mpUnits[ua.ID] = []*types.UserArmyUnit{}
		}

		ut, err := unitTypeRepo.GetTx(ctx, tx, uua.UnitTypeID)
		if err != nil {
			return nil, err
		}
		uua.UnitType = ut

		mpUnits[ua.ID] = append(mpUnits[ua.ID], uua)
	}

	uas := []*types.UserArmy{}
	for result.Next(ctx) {
		node, ok := result.Record().Values[0].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user army")
		}
		ua := types.UserArmyFromNode(node)

		node, ok = result.Record().Values[1].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type game")
		}
		gm := types.GameFromNode(node)
		ua.GameName = gm.Name

		node, ok = result.Record().Values[2].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type user")
		}
		usr := types.UserFromNode(node)
		ua.UserFirstName = usr.FirstName
		ua.UserLastName = usr.LastName
		ua.UserEmail = usr.Email

		node, ok = result.Record().Values[3].(dbtype.Node)
		if !ok {
			return nil, errors.New("unable to convert database type army type")
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

	return uas, nil
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

	currentTimestamp := time.Now().UTC().Unix()
	for _, cuau := range cuaus {
		params := map[string]any{
			"user_army_id": cuau.UserArmyID,
			"unit_type_id": cuau.UnitTypeID,
			"points":       cuau.Points,
			"created_at":   currentTimestamp,
			"updated_at":   currentTimestamp,
		}

		cmd := `
			MATCH (ua:UserArmy{ id: $user_army_id })
			MATCH (ut:UnitType{ id: $unit_type_id })

			MERGE (uau:UserArmyUnit{
				user_army_id: 			$user_army_id
				,unit_type_id:			$unit_type_id
				,quantity:				ut.min_models
			})
			ON CREATE
				SET uau.created_at 	= $created_at
					,uau.id 		= apoc.create.uuid()
					,uau.points 	= ut.min_models * ut.points_per_model
			ON MATCH
				SET uau.updated_at 	= $updated_at
					,uau.points 	= $points
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

			// TODO: this is a possible SLOW QUERY part of the app
			ut, err := unitTypesRepo.GetTx(ctx, tx, eut.ID)
			if err != nil {
				return nil, err
			}
			newUau.UnitType = ut

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
