package db

import (
	"context"
	"strings"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/inconshreveable/log15"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDB() (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(utils.GetEnv("OLDWORLD_API3_HOST", "bolt://localhost:7687"), neo4j.BasicAuth(
		utils.GetEnv("OLDWORLD_API3_USER", "username"), utils.GetEnv("OLDWORLD_API3_PASS", "password"), ""))
	if err != nil {
		log15.Crit("testutils.NewTestDriver: error connecting to database", "err", err)
		return nil, err
	}

	ctx := context.Background()

	qrys := []string{
		// Indexes
		"CREATE INDEX index_user_all IF NOT EXISTS FOR (u:User) ON (u.id,u.email,u.first_name,u.last_name);",
		"CREATE INDEX index_game_all IF NOT EXISTS FOR (g:Game) ON (g.id,g.name);",
		"CREATE INDEX index_statistics_all IF NOT EXISTS FOR (s:Statistic) ON (s.id,s.name,s.display,s.game_id);",
		"CREATE INDEX index_army_type_all IF NOT EXISTS FOR (at:ArmyType) ON (at.id,at.name,at.game_id);",
		"CREATE INDEX index_unit_type_all IF NOT EXISTS FOR (ut:UnitType) ON (ut.id,ut.name,ut.army_type_id,ut.game_id);",
		"CREATE INDEX index_unit_statistics_type_all IF NOT EXISTS FOR (us:UnitStatistic) ON (us.id,us.value,us.unit_type_id,us.statistic_id);",
		"CREATE INDEX index_composition_type_all IF NOT EXISTS FOR (ct:CompositionType) ON (ct.id,ct.name,ct.game_id);",
		"CREATE INDEX index_item_type_all IF NOT EXISTS FOR (it:ItemType) ON (it.id,it.name,it.game_id);",
		"CREATE INDEX index_troop_type_all IF NOT EXISTS FOR (tt:TroopType) ON (tt.id,tt.name,tt.game_id);",
		"CREATE INDEX index_user_armies_all IF NOT EXISTS FOR (ua:UserArmy) ON (ua.id,ua.name,ua.game_id,ua.user_id,ua.army_type_id);",
		"CREATE INDEX index_user_army_units_all IF NOT EXISTS FOR (uau:UserArmyUnit) ON (uau.id,uau.user_army_id,uau.unit_type_id,uau.quantity);",

		// Unique
		"CREATE CONSTRAINT constraint_user_unq_email IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE;",
		"CREATE CONSTRAINT constraint_game_uniq IF NOT EXISTS FOR (g:Game) REQUIRE g.name IS UNIQUE;",
	}

	for _, qry := range qrys {
		if _, err = WriteData(ctx, driver, func(tx neo4j.ManagedTransaction) (any, error) {
			res, e := tx.Run(ctx, qry, make(map[string]any))
			if e != nil {
				return nil, e
			} else if res != nil && res.Err() != nil {
				return nil, res.Err()
			}
			return nil, nil
		}); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return nil, err
			}
		}
	}

	return driver, nil
}

func ReadData(ctx context.Context, driver neo4j.DriverWithContext, fn func(transaction neo4j.ManagedTransaction) (any, error)) (any, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return session.ExecuteWrite(ctx, fn)
}

func WriteData(ctx context.Context, driver neo4j.DriverWithContext, fn func(transaction neo4j.ManagedTransaction) (any, error)) (any, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	return session.ExecuteWrite(ctx, fn)
}
