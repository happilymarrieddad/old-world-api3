package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (ua *UserArmy) GetPbUserArmy() *pbuserarmies.UserArmy {
	pbUa := &pbuserarmies.UserArmy{
		Id:           ua.ID,
		Name:         ua.Name,
		GameId:       ua.GameID,
		GameName:     ua.GameName,
		ArmyTypeId:   ua.ArmyTypeID,
		ArmyTypeName: ua.ArmyTypeName,
		Points:       int64(ua.Points),
		CreatedAt:    timestamppb.New(ua.CreatedAt.UTC()),
	}

	for _, unit := range ua.Units {
		pbUnit := &pbuserarmies.ArmyUnit{
			Id:           unit.ID,
			UserArmyId:   ua.ID,
			UserArmyName: ua.Name,
			UnitTypeId:   unit.UnitTypeID,
			Points:       int64(unit.Points),
			Quantity:     int64(unit.Quantity),
			CreatedAt:    timestamppb.New(unit.CreatedAt.UTC()),
		}

		if unit.UnitType != nil {
			pbUnit.UnitType = unit.UnitType.GetPbUnitType()
		}

		pbUa.Units = append(pbUa.Units, pbUnit)
	}

	return pbUa
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
