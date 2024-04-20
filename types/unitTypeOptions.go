package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type UnitTypeOption struct {
	ID                 string     `json:"id"`
	UnitTypeID         string     `validate:"required" json:"unit_type_id"`
	UnitTypeName       string     `json:"unit_type_name"`
	UnitOptionTypeID   string     `validate:"required" json:"unit_option_type_id"`
	UnitOptionTypeName string     `json:"unit_option_type_name"`
	Txt                string     `validate:"required" json:"txt"`
	DisplayOption      string     `json:"display_option"`
	DisplaySpecialRule string     `json:"display_special_rule"`
	Points             int        `validate:"required" json:"points"`
	PerModel           bool       `json:"per_model"`
	MaxPoints          int        `json:"max_points"`
	Items              []*Item    `json:"items"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

func UnitTypeOptionFromNode(node dbtype.Node) *UnitTypeOption {
	obj := &UnitTypeOption{
		ID:               node.Props["id"].(string),
		UnitTypeID:       node.Props["unit_type_id"].(string),
		UnitOptionTypeID: node.Props["unit_option_type_id"].(string),
		Txt:              node.Props["txt"].(string),
		Points:           GetIntFromNodeProps(node.Props["points"]),
		PerModel:         node.Props["per_model"].(bool),
		MaxPoints:        GetIntFromNodeProps(node.Props["max_pts"]),
		Items:            []*Item{},
	}

	doRaw, ok := node.Props["display_option"].(string)
	if ok {
		obj.DisplayOption = doRaw
	}

	dsrRaw, ok := node.Props["display_special_rule"].(string)
	if ok {
		obj.DisplaySpecialRule = dsrRaw
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
