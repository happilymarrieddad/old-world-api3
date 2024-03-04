package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type UnitStatistic struct {
	ID          string     `json:"id"`
	Value       string     `validate:"required" json:"value"`
	UnitTypeID  string     `json:"unit_type_id"`
	StatisticID string     `validate:"required" json:"statistic_id"`
	Statistic   Statistic  `json:"statistic"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func UnitStatisticFromNode(node dbtype.Node) *UnitStatistic {
	obj := &UnitStatistic{
		ID:          node.Props["id"].(string),
		Value:       node.Props["value"].(string),
		UnitTypeID:  node.Props["unit_type_id"].(string),
		StatisticID: node.Props["statistic_id"].(string),
	}

	timeRaw, ok := node.Props["created_at"].(int64)
	if ok {
		obj.CreatedAt = time.Unix(timeRaw, 0)
	}

	timeRaw, ok = node.Props["updated_at"].(int64)
	if ok {
		obj.UpdatedAt = utils.Ref(time.Unix(timeRaw, 0))
	}

	return obj
}

type UnitType struct {
	ID                  string            `json:"id"`
	Name                string            `validate:"required" json:"name"`
	GameID              string            `validate:"required" json:"game_id"`
	ArmyTypeID          string            `validate:"required" json:"army_type_id"`
	TroopTypeID         string            `validate:"required" json:"troop_type_id"`
	TroopTypeName       string            `json:"troop_type_name"`
	CompositionTypeID   string            `validate:"required" json:"composition_type_id"`
	CompositionTypeName string            `json:"composition_type_name"`
	PointsPerModel      int               `validate:"required" json:"points_per_model"`
	MinModels           int               `validate:"required,min=1" json:"min_models"`
	MaxModels           int               `validate:"required,min=1" json:"max_models"`
	Statistics          []*UnitStatistic  `json:"statistics"`
	Children            []*UnitType       `json:"children"`
	Options             []*UnitTypeOption `json:"options"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           *time.Time        `json:"updated_at"`
}

func UnitTypeFromNode(node dbtype.Node) *UnitType {
	obj := &UnitType{
		ID:             node.Props["id"].(string),
		Name:           node.Props["name"].(string),
		GameID:         node.Props["game_id"].(string),
		ArmyTypeID:     node.Props["army_type_id"].(string),
		PointsPerModel: GetIntFromNodeProps(node.Props["points_per_model"]),
		MinModels:      GetIntFromNodeProps(node.Props["min_models"]),
		MaxModels:      GetIntFromNodeProps(node.Props["max_models"]),
	}

	timeRaw, ok := node.Props["created_at"].(int64)
	if ok {
		obj.CreatedAt = time.Unix(timeRaw, 0)
	}

	timeRaw, ok = node.Props["updated_at"].(int64)
	if ok {
		obj.UpdatedAt = utils.Ref(time.Unix(timeRaw, 0))
	}

	return obj
}

type CreateUnitStatistics struct {
	Display string `validate:"required"`
	Value   string `validate:"required"`
}

type CreateUnitType struct {
	Name              string                  `validate:"required" json:"name"`
	GameID            string                  `validate:"required" json:"game_id"`
	ArmyTypeID        string                  `validate:"required" json:"army_type_id"`
	TroopTypeID       string                  `validate:"required" json:"troop_type_id"`
	CompositionTypeID string                  `validate:"required" json:"composition_type_id"`
	UnitTypeID        string                  `json:"unit_type_id"`
	PointsPerModel    int                     `json:"points_per_model"`
	MinModels         int                     `json:"min_models"`
	MaxModels         int                     `json:"max_models"`
	Statistics        []*CreateUnitStatistics `json:"statistics"`
	UnitOptions       []*UnitTypeOption       `json:"unit_options"`
}

type CreateChildUnitType struct {
	UnitTypeID string                  `json:"unit_type_id"`
	Name       string                  `validate:"required" json:"name"`
	Statistics []*CreateUnitStatistics `json:"statistics"`
}
