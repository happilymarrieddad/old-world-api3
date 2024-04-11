package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type Game struct {
	ID        string     `json:"id"`
	Name      string     `validate:"required" json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func GameFromNode(node dbtype.Node) *Game {
	obj := &Game{
		ID:   node.Props["id"].(string),
		Name: node.Props["name"].(string),
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

type CreateGame struct {
	Name string `validate:"required" json:"name"`
}

type UpdateGame struct {
	ID   string `validate:"required" json:"id"`
	Name string `validate:"required" json:"name"`
}
