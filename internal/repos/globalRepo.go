package repos

import (
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var singular GlobalRepo

//go:generate mockgen -source=./globalRepo.go -destination=./mocks/GlobalRepo.go -package=mock_repos GlobalRepo
type GlobalRepo interface {
	DB() neo4j.DriverWithContext
	ArmyTypes() ArmyTypesRepo
	CompositionTypes() CompositionTypesRepo
	Games() GamesRepo
	Items() ItemsRepo
	ItemTypes() ItemTypesRepo
	Statistics() StatisticsRepo
	TroopTypes() TroopTypesRepo
	UnitOptionTypes() UnitOptionTypesRepo
	UnitTypes() UnitTypesRepo
	UserArmies() UserArmiesRepo
	Users() UsersRepo
}

func NewGlobalRepo(db neo4j.DriverWithContext) (GlobalRepo, error) {
	if singular == nil {
		singular = &globalRepo{
			db:    db,
			mutex: &sync.RWMutex{},
			repos: make(map[string]interface{}),
		}
	}

	return singular, nil
}

type globalRepo struct {
	db    neo4j.DriverWithContext
	repos map[string]interface{}
	mutex *sync.RWMutex
}

func (gr *globalRepo) DB() neo4j.DriverWithContext {
	return gr.db
}

func (gr *globalRepo) factory(key string, fn func(db neo4j.DriverWithContext) interface{}) interface{} {
	gr.mutex.Lock()
	defer gr.mutex.Unlock()

	val, exists := gr.repos[key]
	if exists {
		return val
	}

	newFac := fn(gr.db)
	gr.repos[key] = newFac

	return newFac
}

func (gr *globalRepo) ArmyTypes() ArmyTypesRepo {
	return gr.factory("ArmyTypes", func(db neo4j.DriverWithContext) interface{} { return NewArmyTypesRepo(db) }).(ArmyTypesRepo)
}

func (gr *globalRepo) CompositionTypes() CompositionTypesRepo {
	return gr.factory("CompositionTypes", func(db neo4j.DriverWithContext) interface{} { return NewCompositionTypesRepo(db) }).(CompositionTypesRepo)
}

func (gr *globalRepo) Games() GamesRepo {
	return gr.factory("Games", func(db neo4j.DriverWithContext) interface{} { return NewGamesRepo(db) }).(GamesRepo)
}

func (gr *globalRepo) Items() ItemsRepo {
	return gr.factory("Items", func(db neo4j.DriverWithContext) interface{} { return NewItemsRepo(db) }).(ItemsRepo)
}

func (gr *globalRepo) ItemTypes() ItemTypesRepo {
	return gr.factory("ItemTypes", func(db neo4j.DriverWithContext) interface{} { return NewItemTypesRepo(db) }).(ItemTypesRepo)
}

func (gr *globalRepo) Statistics() StatisticsRepo {
	return gr.factory("Statistics", func(db neo4j.DriverWithContext) interface{} { return NewStatisticsRepo(db) }).(StatisticsRepo)
}

func (gr *globalRepo) TroopTypes() TroopTypesRepo {
	return gr.factory("TroopTypes", func(db neo4j.DriverWithContext) interface{} { return NewTroopTypesRepo(db) }).(TroopTypesRepo)
}

func (gr *globalRepo) UnitOptionTypes() UnitOptionTypesRepo {
	return gr.factory("UnitOptionTypes", func(db neo4j.DriverWithContext) interface{} { return NewUnitOptionTypesRepo(db) }).(UnitOptionTypesRepo)
}

func (gr *globalRepo) UnitTypes() UnitTypesRepo {
	return gr.factory("UnitTypes", func(db neo4j.DriverWithContext) interface{} { return NewUnitTypesRepo(db) }).(UnitTypesRepo)
}

func (gr *globalRepo) UserArmies() UserArmiesRepo {
	return gr.factory("UserArmies", func(db neo4j.DriverWithContext) interface{} { return NewUserArmiesRepo(db) }).(UserArmiesRepo)
}

func (gr *globalRepo) Users() UsersRepo {
	return gr.factory("Users", func(db neo4j.DriverWithContext) interface{} { return NewUsersRepo(db) }).(UsersRepo)
}
