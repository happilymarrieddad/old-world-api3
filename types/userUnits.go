package types

import (
	"log"
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (ua *UserArmyUnit) CalculatePoints() int {
	pts := ua.Quantity * ua.UnitType.PointsPerModel
	OptionValues := []*UpdateUserArmyUnitOptionValue{}

	for _, uauo := range ua.OptionValues {
		multiplier := 1
		if uauo.UnitOption.PerModel {
			multiplier = ua.Quantity
		}
		switch uauo.UnitOption.UnitOptionTypeName {
		case "Single":
			if uauo.IsSelected {
				pts += multiplier * uauo.UnitOption.Points
			}
		case "One Of":
			if uauo.IndexSelected != "" {
				for _, itm := range uauo.UnitOption.Items {
					if itm.ID == uauo.IndexSelected {
						pts += multiplier * itm.Points
					}
				}
			}
		case "Many Of":
			for _, id := range uauo.IDsSelected {
				for _, itm := range uauo.UnitOption.Items {
					if itm.ID == id {
						pts += multiplier * itm.Points
					}
				}
			}
		case "Many To":
			pts += multiplier * uauo.QtySelected
		default:
			log.Printf("unable to find unit option name '%s' for id: '%s'\n", uauo.UnitOptionName, uauo.ID)
		}

		OptionValues = append(OptionValues, &UpdateUserArmyUnitOptionValue{
			ID:            uauo.ID,
			IsSelected:    uauo.IsSelected,
			IndexSelected: uauo.IndexSelected,
			IDsSelected:   uauo.IDsSelected,
			QtySelected:   uauo.QtySelected,
		})
	}

	return pts
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

func (uau *UserArmyUnit) GetPbUserArmyType() *pbuserarmies.ArmyUnit {
	pbUt := &pbuserarmies.ArmyUnit{
		Id:           uau.ID,
		UserArmyId:   uau.UserArmyID,
		UserArmyName: uau.UserArmyName,
		UnitTypeId:   uau.UnitTypeID,
		UnitType:     uau.UnitType.GetPbUnitType(),
		Points:       int64(uau.Points),
		Quantity:     int64(uau.Quantity),
		CreatedAt:    timestamppb.New(uau.CreatedAt.UTC()),
	}

	for idx, opt := range pbUt.UnitType.Options {
		for _, uopt := range uau.OptionValues {
			if uopt.UnitOptionID == opt.Id {
				pbUt.UnitType.Options[idx].Value = &unittypes.UnitOptionValue{
					Id:               uopt.ID,
					UserArmyUnitId:   uau.ID,
					UserArmyUnitName: uopt.UserArmyUnitName,
					UnitOptionId:     uopt.UnitOptionID,
					UnitOptionName:   uopt.UnitOptionName,
					IsSelected:       uopt.IsSelected,
					IndexSelected:    uopt.IndexSelected,
					IdsSelected:      uopt.IDsSelected,
					QtySelected:      int64(uopt.QtySelected),
					CreatedAt:        timestamppb.New(uopt.CreatedAt.UTC()),
				}
			}
		}
	}

	return pbUt
}

type CreateUserArmyUnit struct {
	UserArmyID string `validate:"required" json:"user_army_id"`
	UnitTypeID string `validate:"required" json:"unit_type_id"`
	Points     int    `json:"points"`
}

type UpdateUserArmyUnitOptionValue struct {
	ID            string   `json:"id"`
	IsSelected    bool     `json:"is_selected"`
	IndexSelected string   `json:"index_selected"`
	IDsSelected   []string `json:"ids_selected"`
	QtySelected   int      `json:"qty_selected"`
}

type UpdateUserArmyUnit struct {
	ID           string                           `validate:"required" json:"id"`
	Qty          *int                             `json:"qty"`
	Points       *int                             `json:"points"`
	OptionValues []*UpdateUserArmyUnitOptionValue `json:"option_values"`
}

type UserArmyUnitOptionValue struct {
	ID               string          `json:"id"`
	UserArmyUnitID   string          `json:"user_army_unit_id"`
	UserArmyUnitName string          `json:"user_army_unit_name"`
	UnitOptionID     string          `json:"unit_option_id"`
	UnitOptionName   string          `json:"unit_option_name"`
	IsSelected       bool            `json:"is_selected"`
	IndexSelected    string          `json:"index_selected"`
	IDsSelected      []string        `json:"ids_selected"`
	QtySelected      int             `json:"qty_selected"`
	UnitOption       *UnitTypeOption `json:"unit_option"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at"`
}

func UserArmyUnitOptionValueFromNode(node dbtype.Node) *UserArmyUnitOptionValue {
	obj := &UserArmyUnitOptionValue{
		ID:             node.Props["id"].(string),
		UserArmyUnitID: node.Props["user_army_unit_id"].(string),
		UnitOptionID:   node.Props["unit_option_id"].(string),
		IsSelected:     node.Props["is_selected"].(bool),
		IndexSelected:  node.Props["index_selected"].(string),
		QtySelected:    GetIntFromNodeProps(node.Props["qty_selected"]),
	}

	// IDsSelected:    node.Props["ids_selected"].([]string),
	idsSelectedIface, ok := node.Props["ids_selected"].([]interface{})
	if ok {
		for _, idIface := range idsSelectedIface {
			obj.IDsSelected = append(obj.IDsSelected, idIface.(string))
		}
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
