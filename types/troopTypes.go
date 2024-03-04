package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type TroopType struct {
	ID        string     `json:"id"`
	Name      string     `validate:"required" json:"name"`
	GameID    string     `validate:"required" json:"game_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func TroopTypeFromNode(node dbtype.Node) *TroopType {
	obj := &TroopType{
		ID:     node.Props["id"].(string),
		Name:   node.Props["name"].(string),
		GameID: node.Props["game_id"].(string),
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

type CreateTroopType struct {
	Name   string `validate:"required" json:"name"`
	GameID string `validate:"required" json:"game_id"`
}
