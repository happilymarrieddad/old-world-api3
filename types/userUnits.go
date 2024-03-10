package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type UserArmyUnit struct {
	ID           string                     `json:"id"`
	UserArmyID   string                     `validate:"required" json:"user_army_id"`
	UserArmyName string                     `json:"user_army_name"`
	UnitTypeID   string                     `json:"unit_type_id"`
	UnitType     *UnitType                  `json:"unit_type"`
	Points       int                        `validate:"required" json:"points"`
	Quantity     int                        `validate:"required" json:"quantity"`
	OptionValues []*UserArmyUnitOptionValue `json:"option_values"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    *time.Time                 `json:"updated_at"`
}

func UserArmyUnitFromNode(node dbtype.Node) *UserArmyUnit {
	obj := &UserArmyUnit{
		ID:         node.Props["id"].(string),
		UserArmyID: node.Props["user_army_id"].(string),
		UnitTypeID: node.Props["unit_type_id"].(string),
		UnitType:   &UnitType{},
		Points:     GetIntFromNodeProps(node.Props["points"]),
		Quantity:   GetIntFromNodeProps(node.Props["quantity"]),
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

type CreateUserArmyUnit struct {
	UserArmyID string `validate:"required" json:"user_army_id"`
	UnitTypeID string `validate:"required" json:"unit_type_id"`
	Points     int    `json:"points"`
}

type UserArmyUnitOptionValue struct {
	ID                 string          `json:"id"`
	UserArmyUnitID     string          `json:"user_army_unit_id"`
	UserArmyUnitName   string          `json:"user_army_unit_name"`
	UnitTypeOptionID   string          `json:"unit_type_option_id"`
	UnitTypeOptionName string          `json:"unit_type_option_name"`
	IsSelected         bool            `json:"is_selected"`
	IndexSelected      string          `json:"index_selected"`
	IDsSelected        []string        `json:"ids_selected"`
	QtySelected        int             `json:"qty_selected"`
	UnitTypeOption     *UnitTypeOption `json:"unit_type_option"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          *time.Time      `json:"updated_at"`
}

func UserArmyUnitOptionValueFromNode(node dbtype.Node) *UserArmyUnitOptionValue {
	obj := &UserArmyUnitOptionValue{
		ID:               node.Props["id"].(string),
		UserArmyUnitID:   node.Props["user_army_unit_id"].(string),
		UnitTypeOptionID: node.Props["unit_type_option_id"].(string),
		IsSelected:       node.Props["is_selected"].(bool),
		IndexSelected:    node.Props["index_selected"].(string),
		IDsSelected:      node.Props["ids_selected"].([]string),
		QtySelected:      GetIntFromNodeProps(node.Props["qty_selected"]),
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
