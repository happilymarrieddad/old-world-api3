package types

import (
	"time"

	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UnitStatistic struct {
	ID          string     `json:"id"`
	Value       string     `validate:"required" json:"value"`
	UnitTypeID  string     `json:"unit_type_id"`
	StatisticID string     `validate:"required" json:"statistic_id"`
	Statistic   Statistic  `json:"statistic"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func UnitStatisticFromNode(node dbtype.Node) *UnitStatistic {
	obj := &UnitStatistic{
		ID:          node.Props["id"].(string),
		Value:       node.Props["value"].(string),
		UnitTypeID:  node.Props["unit_type_id"].(string),
		StatisticID: node.Props["statistic_id"].(string),
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

type UnitType struct {
	ID                  string            `json:"id"`
	Name                string            `validate:"required" json:"name"`
	GameID              string            `validate:"required" json:"game_id"`
	ArmyTypeID          string            `validate:"required" json:"army_type_id"`
	TroopTypeID         string            `validate:"required" json:"troop_type_id"`
	TroopTypeName       string            `json:"troop_type_name"`
	CompositionTypeID   string            `validate:"required" json:"composition_type_id"`
	CompositionTypeName string            `json:"composition_type_name"`
	PointsPerModel      int               `validate:"required" json:"points_per_model"`
	MinModels           int               `validate:"required,min=1" json:"min_models"`
	MaxModels           int               `validate:"required,min=1" json:"max_models"`
	Statistics          []*UnitStatistic  `json:"statistics"`
	Children            []*UnitType       `json:"children"`
	Options             []*UnitTypeOption `json:"options"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           *time.Time        `json:"updated_at"`
}

func UnitTypeFromNode(node dbtype.Node) *UnitType {
	obj := &UnitType{
		ID:             node.Props["id"].(string),
		Name:           node.Props["name"].(string),
		GameID:         node.Props["game_id"].(string),
		ArmyTypeID:     node.Props["army_type_id"].(string),
		PointsPerModel: GetIntFromNodeProps(node.Props["points_per_model"]),
		MinModels:      GetIntFromNodeProps(node.Props["min_models"]),
		MaxModels:      GetIntFromNodeProps(node.Props["max_models"]),
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

func (ut *UnitType) GetPbUnitType() *pbunittypes.UnitType {
	pbUt := &pbunittypes.UnitType{
		Id:                  ut.ID,
		Name:                ut.Name,
		GameId:              ut.GameID,
		ArmyTypeId:          ut.ArmyTypeID,
		TroopTypeId:         ut.TroopTypeID,
		TroopTypeName:       ut.TroopTypeName,
		CompositionTypeId:   ut.CompositionTypeID,
		CompositionTypeName: ut.CompositionTypeName,
		PointsPerModel:      int64(ut.PointsPerModel),
		MinModels:           int64(ut.MinModels),
		MaxModels:           int64(ut.MaxModels),
		CreatedAt:           timestamppb.New(ut.CreatedAt.UTC()),
	}

	for _, stat := range ut.Statistics {
		pbStat := &pbunittypes.UnitStatistic{
			Id:          stat.ID,
			Value:       stat.Value,
			UnitTypeId:  ut.ID,
			StatisticId: stat.StatisticID,
			Statistic: &pbunittypes.Statistic{
				Id:        stat.StatisticID,
				Name:      stat.Statistic.Name,
				Display:   stat.Statistic.Display,
				GameId:    stat.Statistic.GameID,
				CreatedAt: timestamppb.New(stat.CreatedAt.UTC()),
			},
		}

		pbUt.Statistics = append(pbUt.Statistics, pbStat)
	}

	for _, child := range ut.Children {
		pbChild := &pbunittypes.UnitTypeChild{
			Id:                  child.ID,
			Name:                child.Name,
			GameId:              child.GameID,
			ArmyTypeId:          child.ArmyTypeID,
			TroopTypeId:         child.TroopTypeID,
			TroopTypeName:       child.TroopTypeName,
			CompositionTypeId:   child.CompositionTypeID,
			CompositionTypeName: child.CompositionTypeName,
			PointsPerModel:      int64(child.PointsPerModel),
			MinModels:           int64(child.MinModels),
			MaxModels:           int64(child.MaxModels),
			CreatedAt:           timestamppb.New(child.CreatedAt.UTC()),
		}

		for _, childStat := range child.Statistics {
			pbChild.Statistics = append(pbChild.Statistics, &pbunittypes.UnitStatistic{
				Id:          childStat.ID,
				Value:       childStat.Value,
				UnitTypeId:  ut.ID,
				StatisticId: childStat.StatisticID,
				Statistic: &pbunittypes.Statistic{
					Id:        childStat.StatisticID,
					Name:      childStat.Statistic.Name,
					Display:   childStat.Statistic.Display,
					GameId:    childStat.Statistic.GameID,
					CreatedAt: timestamppb.New(childStat.CreatedAt.UTC()),
				},
			})
		}

		pbUt.Children = append(pbUt.Children, pbChild)
	}

	for _, opt := range ut.Options {
		pbOption := &pbunittypes.UnitTypeOption{
			Id:                 opt.ID,
			UnitTypeId:         opt.UnitTypeID,
			UnitTypeName:       opt.UnitTypeName,
			UnitOptionTypeId:   opt.UnitOptionTypeID,
			UnitOptionTypeName: opt.UnitOptionTypeName,
			Txt:                opt.Txt,
			Points:             int64(opt.Points),
			PerModel:           opt.PerModel,
			MaxPoints:          int64(opt.MaxPoints),
			CreatedAt:          timestamppb.New(opt.CreatedAt.UTC()),
		}

		for _, itm := range opt.Items {
			var atId string
			if itm.ArmyTypeID != nil { // very unlikely but we can't just grab a pointers value
				atId = *itm.ArmyTypeID
			}
			pbOption.Items = append(pbOption.Items, &pbunittypes.Item{
				Id:           itm.ID,
				Name:         itm.Name,
				Points:       int64(itm.Points),
				ItemTypeId:   itm.ItemTypeID,
				ItemTypeName: itm.ItemTypeName,
				GameId:       itm.GameID,
				ArmyTypeId:   atId,
				CreatedAt:    timestamppb.New(itm.CreatedAt.UTC()),
			})
		}

		pbUt.Options = append(pbUt.Options, pbOption)
	}

	return pbUt
}

type CreateUnitStatistics struct {
	Display string `validate:"required"`
	Value   string `validate:"required"`
}

type CreateUnitType struct {
	Name              string                  `validate:"required" json:"name"`
	GameID            string                  `validate:"required" json:"game_id"`
	ArmyTypeID        string                  `validate:"required" json:"army_type_id"`
	TroopTypeID       string                  `validate:"required" json:"troop_type_id"`
	CompositionTypeID string                  `validate:"required" json:"composition_type_id"`
	UnitTypeID        string                  `json:"unit_type_id"`
	PointsPerModel    int                     `json:"points_per_model"`
	MinModels         int                     `json:"min_models"`
	MaxModels         int                     `json:"max_models"`
	Statistics        []*CreateUnitStatistics `json:"statistics"`
	UnitOptions       []*UnitTypeOption       `json:"unit_options"`
}

type CreateChildUnitType struct {
	UnitTypeID string                  `json:"unit_type_id"`
	Name       string                  `validate:"required" json:"name"`
	Statistics []*CreateUnitStatistics `json:"statistics"`
}
