package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type UserArmy struct {
	ID            string          `json:"id"`
	Name          string          `validate:"required" json:"value"`
	UserID        string          `validate:"required" json:"user_id"`
	UserFirstName string          `json:"user_first_name"`
	UserLastName  string          `json:"user_last_name"`
	UserEmail     string          `json:"user_email"`
	GameID        string          `validate:"required" json:"game_id"`
	GameName      string          `json:"game_name"`
	ArmyTypeID    string          `validate:"required" json:"army_type_id"`
	ArmyTypeName  string          `json:"army_type_name"`
	Points        int             `validate:"required" json:"points"`
	Units         []*UserArmyUnit `json:"units"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
}

type CreateUserArmy struct {
	Name       string `validate:"required" json:"value"`
	UserID     string `validate:"required" json:"user_id"`
	GameID     string `validate:"required" json:"game_id"`
	ArmyTypeID string `validate:"required" json:"army_type_id"`
	Points     int    `validate:"required" json:"points"`
}

func UserArmyFromNode(node dbtype.Node) *UserArmy {
	obj := &UserArmy{
		ID:         node.Props["id"].(string),
		Name:       node.Props["name"].(string),
		UserID:     node.Props["user_id"].(string),
		GameID:     node.Props["game_id"].(string),
		ArmyTypeID: node.Props["army_type_id"].(string),
		Points:     GetIntFromNodeProps(node.Props["points"]),
		Units:      []*UserArmyUnit{},
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
