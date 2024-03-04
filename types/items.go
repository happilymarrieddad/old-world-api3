package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Item struct {
	ID           string     `json:"id"`
	Name         string     `validate:"required" json:"name"`
	Points       int        `validate:"required" json:"points"`
	ItemTypeID   string     `validate:"required" json:"item_type_id"`
	ItemTypeName string     `json:"item_type_name"`
	GameID       string     `validate:"required" json:"game_id"`
	ArmyTypeID   *string    `json:"army_type_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func ItemFromNode(node dbtype.Node) *Item {
	obj := &Item{
		ID:         node.Props["id"].(string),
		Name:       node.Props["name"].(string),
		GameID:     node.Props["game_id"].(string),
		ItemTypeID: node.Props["item_type_id"].(string),
		Points:     GetIntFromNodeProps(node.Props["points"]),
	}

	atID, ok := node.Props["army_type_id"]
	if ok {
		obj.ArmyTypeID = utils.Ref(atID.(string))
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

type CreateItem struct {
	Name       string  `validate:"required" json:"name"`
	Points     int     `validate:"required" json:"points"`
	GameID     string  `validate:"required" json:"game_id"`
	ArmyTypeID *string `json:"army_type_id"`
	ItemTypeID string  `validate:"required" json:"item_type_id"`
	Debug      bool
}
