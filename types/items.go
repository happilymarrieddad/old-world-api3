package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	pbitems "github.com/happilymarrieddad/old-world/api3/pb/proto/items"
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
	Description  string     `json:"description"`
	Story        string     `json:"story"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func PbItemFromItem(itm *Item) *pbitems.Item {
	pbItem := new(pbitems.Item)

	pbItem.Id = itm.ID
	pbItem.Name = itm.Name
	pbItem.Points = int64(itm.Points)
	pbItem.ItemTypeId = itm.ItemTypeID
	pbItem.ItemTypeName = itm.ItemTypeName
	pbItem.GameId = itm.GameID
	pbItem.Description = itm.Description
	pbItem.Story = itm.Story
	if itm.ArmyTypeID != nil {
		pbItem.ArmyTypeId = *itm.ArmyTypeID
	}

	return pbItem
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

	desStr, ok := node.Props["description"]
	if ok {
		obj.Description = desStr.(string)
	}

	storStr, ok := node.Props["story"]
	if ok {
		obj.Story = storStr.(string)
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
	Name        string  `validate:"required" json:"name"`
	Points      int     `validate:"required" json:"points"`
	GameID      string  `validate:"required" json:"game_id"`
	ArmyTypeID  *string `json:"army_type_id"`
	ItemTypeID  string  `validate:"required" json:"item_type_id"`
	Description string  `json:"description"`
	Story       string  `json:"story"`
	Debug       bool
}

func NewItemFromPbCreateItem(itm *pbitems.CreateArmyItem) CreateItem {
	newItm := CreateItem{
		Name:        itm.GetName(),
		Points:      int(itm.GetPoints()),
		GameID:      itm.GetGameId(),
		ItemTypeID:  itm.GetItemTypeId(),
		Description: itm.GetDescription(),
		Story:       itm.GetStory(),
	}

	if len(itm.ArmyTypeId) > 0 {
		newItm.ArmyTypeID = utils.Ref(itm.GetArmyTypeId())
	}

	return newItm
}

type UpdateItem struct {
	ID          string `validate:"required" json:"id"`
	Name        string `validate:"required" json:"name"`
	Points      int64  `validate:"required" json:"points"`
	Description string `validate:"required" json:"description"`
	Story       string `validate:"required" json:"story"`
}

func NewItemFromPbUpdateItem(itm *pbitems.UpdateItem) UpdateItem {
	return UpdateItem{
		ID:          itm.GetId(),
		Name:        itm.GetName(),
		Points:      itm.GetPoints(),
		Description: itm.GetDescription(),
		Story:       itm.GetStory(),
	}
}
