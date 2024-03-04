package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Statistic struct {
	ID        string     `json:"id"`
	Name      string     `validate:"required" json:"name"`
	Display   string     `validate:"required" json:"display"`
	GameID    string     `validate:"required" json:"game_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func StatisticFromNode(node dbtype.Node) *Statistic {
	obj := &Statistic{
		ID:      node.Props["id"].(string),
		Name:    node.Props["name"].(string),
		Display: node.Props["display"].(string),
		GameID:  node.Props["game_id"].(string),
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

type CreateStatistic struct {
	Name    string `validate:"required" json:"name"`
	Display string `validate:"required" json:"display"`
	GameID  string `validate:"required" json:"game_id"`
}
