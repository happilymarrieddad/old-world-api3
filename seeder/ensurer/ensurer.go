package ensurer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
)

func EnsureData(ctx context.Context, gr repos.GlobalRepo, ad Games) error {
	t := time.Now()
	log.Println("Starting loading data")

	if _, err := gr.Users().FindOrCreate(ctx, types.CreateUser{
		FirstName:       "Nick",
		LastName:        "Kotenberg",
		Email:           "nick@mail.com",
		Password:        "1234",
		PasswordConfirm: "1234",
	}); err != nil {
		return err
	}

	if _, err := gr.Users().FindOrCreate(ctx, types.CreateUser{
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@mail.com",
		Password:        "1234",
		PasswordConfirm: "1234",
	}); err != nil {
		return err
	}

	gameIDX := 1
	for gameName, gameData := range ad {
		fmt.Printf("Running game: %s (%d of %d)\n", gameName, gameIDX, len(ad))
		gm, err := gr.Games().FindOrCreate(ctx, gameName)
		if err != nil {
			return err
		}
		gameIDX++

		for _, stat := range gameData.StatisticNames {
			if _, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: stat.Name, Display: stat.Display, GameID: gm.ID}); err != nil {
				return err
			}
		}

		for _, tt := range gameData.TroopTypes {
			if _, err = gr.TroopTypes().FindOrCreate(ctx, types.CreateTroopType{Name: tt, GameID: gm.ID}); err != nil {
				return err
			}
		}

		ctMap := make(map[string]*types.CompositionType)
		for _, ct := range gameData.CompositionTypes {
			newCt, err := gr.CompositionTypes().FindOrCreate(ctx, types.CreateCompositionType{Name: ct, GameID: gm.ID})
			if err != nil {
				return err
			}
			ctMap[ct] = newCt
		}

		uotMap := make(map[string]*types.UnitOptionType)
		for _, uot := range gameData.UnitOptionTypes {
			newUot, err := gr.UnitOptionTypes().FindOrCreate(ctx, types.CreateUnitOptionType{Name: uot, GameID: gm.ID})
			if err != nil {
				return err
			}
			uotMap[newUot.Name] = newUot
		}

		itMap := make(map[string]*types.ItemType)
		for _, it := range gameData.ItemTypes {
			newIt, err := gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{Name: it, GameID: gm.ID})
			if err != nil {
				return err
			}
			itMap[it] = newIt
		}

		for _, itm := range gameData.Items {
			it, ok := itMap[itm.Type]
			if !ok {
				return fmt.Errorf("unable to find item type Name: %s Type: %s", itm.Name, itm.Type)
			}

			if _, err := gr.Items().Create(ctx, types.CreateItem{
				Name:        itm.Name,
				Points:      int(itm.Points),
				GameID:      gm.ID,
				ItemTypeID:  it.ID,
				Description: itm.Description,
				Story:       itm.Story,
			}); err != nil {
				return err
			}
		}

		for armyIDX, armyData := range gameData.Armies {
			fmt.Printf("Doing army: %s (%d of %d)\n", armyData.Name, armyIDX+1, len(gameData.Armies))
			at, err := gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: armyData.Name, GameID: gm.ID})
			if err != nil {
				return err
			}

			for _, aitm := range armyData.Items {
				it, ok := itMap[aitm.Type]
				if !ok {
					return fmt.Errorf("unable to find item type 2: Name: %s Type: %s", aitm.Name, aitm.Type)
				}

				if _, err := gr.Items().Create(ctx, types.CreateItem{
					Name:        aitm.Name,
					Points:      int(aitm.Points),
					GameID:      gm.ID,
					ItemTypeID:  it.ID,
					ArmyTypeID:  &at.ID,
					Description: aitm.Description,
					Story:       aitm.Story,
				}); err != nil {
					return err
				}
			}

			unitTypeRepo := gr.UnitTypes()
			for _, unitType := range armyData.UnitTypes {
				tt, err := gr.TroopTypes().FindOrCreate(ctx, types.CreateTroopType{Name: unitType.TroopTypeName, GameID: gm.ID})
				if err != nil {
					return err
				}

				ct, ok := ctMap[unitType.CompositionTypeName]
				if !ok {
					return errors.New("unable to find composition type")
				}

				nut := types.CreateUnitType{
					Name:              unitType.Name,
					GameID:            gm.ID,
					ArmyTypeID:        at.ID,
					TroopTypeID:       tt.ID,
					CompositionTypeID: ct.ID,
					PointsPerModel:    int(unitType.PointsPerModel),
					MinModels:         int(unitType.MinModels),
					MaxModels:         int(unitType.MaxModels),
					// TODO: add these in neo4j query
					UnitOptions: []*types.UnitTypeOption{},
				}

				for _, opt := range unitType.Options {
					uot, ok := uotMap[opt.UnitOptionTypeName]
					if !ok {
						return fmt.Errorf("unable to find unit option type name: %s", opt.UnitOptionTypeName)
					}

					nuo := &types.UnitTypeOption{
						UnitTypeID:         "", // Should get in the create func
						UnitTypeName:       unitType.Name,
						UnitOptionTypeID:   uot.ID,
						UnitOptionTypeName: uot.Name,
						Txt:                opt.Txt,
						Points:             int(opt.Points),
						PerModel:           opt.PerModel,
						MaxPoints:          int(opt.MaxPts),
						Items:              []*types.Item{},
					}

					if opt.IsMagicStandards {
						existingItemsData, err := gr.Items().Find(ctx, &repos.FindItemsOpts{
							GameID:      gm.ID,
							ArmyTypeID:  &at.ID,
							ItemTypeIDs: []string{itMap["Magic Standards"].ID},
						})
						if err != nil {
							return err
						}
						for _, existingItemData := range existingItemsData {
							if existingItemData.Points <= nuo.MaxPoints {
								nuo.Items = append(nuo.Items, existingItemData)
							}
						}
					} else if opt.IsItems {
						existingItemsData, err := gr.Items().Find(ctx, &repos.FindItemsOpts{
							GameID:     gm.ID,
							ArmyTypeID: &at.ID,
							ItemTypeIDs: []string{
								itMap["Magic Weapons"].ID,
								itMap["Magic Armour"].ID,
								itMap["Talismans"].ID,
								itMap["Enchanted Items"].ID,
								itMap["Arcane Items"].ID,
							},
						})
						if err != nil {
							return err
						}
						for _, existingItemData := range existingItemsData {
							if existingItemData.Points <= nuo.MaxPoints {
								nuo.Items = append(nuo.Items, existingItemData)
							}
						}
					} else {
						for _, optItem := range opt.PossibleValues {
							it, ok := itMap[optItem.Type]
							if !ok {
								return fmt.Errorf("unable to find item type 3: Name: %s Type: %s", optItem.Txt, optItem.Type)
							}
							// TODO: find the item and add it to the unit options
							existingItems, err := gr.Items().Find(ctx, &repos.FindItemsOpts{
								GameID:     gm.ID,
								ArmyTypeID: &at.ID,
								Names:      []string{optItem.Txt},
								Limit:      1,
							})
							if err != nil {
								return err
							} else if len(existingItems) > 0 {
								nuo.Items = append(nuo.Items, existingItems[0])
							} else {
								if nuo.MaxPoints == 0 || optItem.Points <= int64(nuo.MaxPoints) {
									// Item does not exist so we must create it
									newItem, err := gr.Items().Create(ctx, types.CreateItem{
										Name:        optItem.Txt,
										Points:      int(optItem.Points),
										GameID:      gm.ID,
										ArmyTypeID:  &at.ID,
										ItemTypeID:  it.ID,
										Description: optItem.Description,
										Story:       optItem.Story,
									})
									if err != nil {
										return err
									}
									nuo.Items = append(nuo.Items, newItem)
								}
							}
						}
					}

					nut.UnitOptions = append(nut.UnitOptions, nuo)
				}

				if len(unitType.ParentUnitName) > 0 {
					parent, _, err := unitTypeRepo.Find(ctx, &repos.FindUnitTypesOpts{
						Name: unitType.ParentUnitName, ArmyTypeID: at.ID, Limit: 1,
					})
					if err != nil {
						return err
					} else if len(parent) == 0 {
						spew.Dump(unitType)
						fmt.Println(at.ID)
						return errors.New("somehow parent was created then can't be found? JSON is most certainly bad somewhere")
					}

					nut.UnitTypeID = parent[0].ID
				}

				for display, value := range unitType.UnitStatistics {
					nut.Statistics = append(nut.Statistics, &types.CreateUnitStatistics{
						Display: display,
						Value:   value,
					})
				}

				if _, err := unitTypeRepo.FindOrCreate(ctx, nut); err != nil {
					return err
				}
			}
		}
	}

	log.Printf("Finished loading data which took %s\n", time.Since(t).String())
	return nil
}

type ItemData struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Points      int64  `json:"points"`
	Description string `json:"description"`
	Story       string `json:"story"`
}

type StatisticNameData struct {
	Name    string `json:"name"`
	Display string `json:"display"`
}

type UnitStatisticsValue map[string]string

type UnitOptionsPossibleValuesData struct {
	Txt         string `json:"txt"`
	Points      int64  `json:"points"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Story       string `json:"story"`
}

type UnitOptionData struct {
	Txt                string                           `json:"txt"`
	UnitOptionTypeName string                           `json:"unit_option_type_name"`
	Points             int64                            `json:"points"`
	PerModel           bool                             `json:"per_model"`
	PossibleValues     []*UnitOptionsPossibleValuesData `json:"possible_values"`
	IsItems            bool                             `json:"is_magic_items"`
	IsMagicStandards   bool                             `json:"is_magic_standards"`
	MaxPts             int64                            `json:"max_pts"`
	MustChooseOne      bool                             `json:"must_choose_one"`
}

type UnitTypeData struct {
	Name                string              `json:"name"`
	UnitStatistics      UnitStatisticsValue `json:"unit_statistics"`
	PointsPerModel      int64               `json:"points_per_model"`
	MinModels           int64               `json:"min_models"`
	MaxModels           int64               `json:"max_models"`
	ParentUnitName      string              `json:"parent_unit_type_name"`
	TroopTypeName       string              `json:"troop_type_name"`
	CompositionTypeName string              `json:"composition_type_name"`
	Options             []*UnitOptionData   `json:"options"`
}

type ArmyTypeData struct {
	Name      string          `json:"name"`
	UnitTypes []*UnitTypeData `json:"unit_types"`
	Items     []*ItemData     `json:"items"`
}

type GameData struct {
	StatisticNames   []*StatisticNameData `json:"statistic_names"`
	TroopTypes       []string             `json:"troop_types"`
	CompositionTypes []string             `json:"composition_types"`
	Armies           []*ArmyTypeData      `json:"army_types"`
	ItemTypes        []string             `json:"item_types"`
	Items            []*ItemData          `json:"items"`
	UnitOptionTypes  []string             `json:"unit_option_types"`
}

type Games map[string]*GameData
