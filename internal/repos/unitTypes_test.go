package repos_test

import (
	"encoding/json"
	"os"

	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/seeder/ensurer"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
)

var _ = Describe("repo:UnitTypes", func() {
	var repo repos.UnitTypesRepo
	var gm *types.Game
	var at *types.ArmyType
	var tt *types.TroopType
	var ct *types.CompositionType

	BeforeEach(func() {
		repo = gr.UnitTypes()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())

		at, err = gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "Bretonnia",
			GameID: gm.ID})
		Expect(err).To(BeNil())

		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Movement", Display: "M", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Weapon Skill", Display: "WS", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Ballistic Skill", Display: "BS", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Strength", Display: "S", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Toughnous", Display: "T", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Wounds", Display: "W", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Initiative", Display: "I", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Attacks", Display: "A", GameID: at.GameID})
		Expect(err).To(BeNil())
		_, err = gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{Name: "Leadership", Display: "Ld", GameID: at.GameID})
		Expect(err).To(BeNil())

		tt, err = gr.TroopTypes().FindOrCreate(ctx, types.CreateTroopType{Name: "Regular Infantry", GameID: at.GameID})
		Expect(err).To(BeNil())

		ct, err = gr.CompositionTypes().FindOrCreate(ctx, types.CreateCompositionType{Name: "Core", GameID: at.GameID})
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should just merge the same unit statistics", func() {
			ut, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Baron",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut.ID).NotTo(HaveLen(0))
			Expect(ut.Statistics).To(HaveLen(1))
		})

		It("should successfully create a unit type", func() {
			ut, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Baron",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
					{Display: "WS", Value: "7"},
					{Display: "BS", Value: "5"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut.ID).NotTo(HaveLen(0))
			Expect(ut.Statistics).To(HaveLen(3))
			Expect(ut.Statistics).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Movement"),
						"Display": Equal("M"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("7"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Weapon Skill"),
						"Display": Equal("WS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("5"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Ballistic Skill"),
						"Display": Equal("BS"),
					}),
				})),
			))

			ut2, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Baron",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
			})
			Expect(err).To(BeNil())
			Expect(ut2.ID).NotTo(HaveLen(0))
			Expect(ut.ID).To(Equal(ut2.ID))
			// This is a FIND so it SHOULD have the statistics already done
			Expect(ut.Statistics).To(HaveLen(3))
			Expect(ut.Statistics).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Movement"),
						"Display": Equal("M"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("7"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Weapon Skill"),
						"Display": Equal("WS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("5"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Ballistic Skill"),
						"Display": Equal("BS"),
					}),
				})),
			))

			ut3, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Men at Arms",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut3.ID).NotTo(HaveLen(0))
			Expect(ut3.ID).NotTo(Equal(at.ID))
			Expect(ut3.ID).NotTo(Equal(ut2.ID))
			Expect(ut3.Statistics).To(HaveLen(1))
			Expect(ut.Statistics).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Movement"),
						"Display": Equal("M"),
					}),
				})),
			))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			ut, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Baron",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut.ID).NotTo(HaveLen(0))
			Expect(ut.PointsPerModel).To(BeNumerically("==", 10))
			Expect(ut.MinModels).To(BeNumerically("==", 1))
			Expect(ut.MaxModels).To(BeNumerically("==", 1))

			ut2, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Men at Arms",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut2.ID).NotTo(HaveLen(0))

			ut3, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Duke",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of unit types", func() {
			ats, _, err := repo.Find(ctx, &repos.FindUnitTypesOpts{ArmyTypeID: at.ID, Limit: 10})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(3))

			for _, at := range ats {
				Expect(at.Statistics).To(HaveLen(1))
				Expect(at.Statistics).To(ContainElements(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Value": Equal("4"),
						"Statistic": MatchFields(IgnoreExtras, Fields{
							"Name":    Equal("Movement"),
							"Display": Equal("M"),
						}),
					})),
				))
			}

			ats, _, err = repo.Find(ctx, &repos.FindUnitTypesOpts{ArmyTypeID: at.ID, Limit: 1})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))
		})
	})

	Context("GetNamesByArmyTypeID", func() {
		BeforeEach(func() {
			ut, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Baron",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut.ID).NotTo(HaveLen(0))
			Expect(ut.PointsPerModel).To(BeNumerically("==", 10))
			Expect(ut.MinModels).To(BeNumerically("==", 1))
			Expect(ut.MaxModels).To(BeNumerically("==", 1))

			ut2, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Men at Arms",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut2.ID).NotTo(HaveLen(0))

			ut3, err := repo.FindOrCreate(ctx, types.CreateUnitType{
				Name:              "Duke",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    10,
				MinModels:         1,
				MaxModels:         1,
				Statistics: []*types.CreateUnitStatistics{
					{Display: "M", Value: "4"},
				},
			})
			Expect(err).To(BeNil())
			Expect(ut3.ID).NotTo(HaveLen(0))
		})

		It("should successfully get the unit types by army type id", func() {
			uts, err := repo.GetNamesByArmyTypeID(ctx, at.ID)
			Expect(err).To(BeNil())
			Expect(uts).To(HaveLen(3))
		})
	})

	Context("Heavy Testing", func() {
		BeforeEach(func() {
			if os.Getenv("DEBUG") != "true" {
				Skip("skipping heavy testing")
			}

			bts, err := os.ReadFile("data/example.json")
			if err != nil {
				panic(err)
			}

			ad := make(ensurer.Games)
			if err := json.Unmarshal(bts, &ad); err != nil {
				panic(errors.WithMessage(err, "unable to marshal data"))
			}

			clearAllData()

			Expect(ensurer.EnsureData(ctx, gr, ad)).To(Succeed())
		})

		It("should check for data that's been 'ensured'", func() {
			var err error
			gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
			Expect(err).To(BeNil())

			at, err = gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "Beastmen Brayherds", GameID: gm.ID})
			Expect(err).To(BeNil())

			uts, _, err := repo.Find(ctx, &repos.FindUnitTypesOpts{Limit: 20, ArmyTypeID: at.ID, Name: "Gor Herds"})
			Expect(err).To(BeNil())
			Expect(uts).To(HaveLen(1))

			ut := uts[0]

			Expect(ut.Name).To(Equal("Gor Herds"))
			Expect(ut.TroopTypeName).To(Equal("Regular Infantry"))
			Expect(ut.CompositionTypeName).To(Equal("Core"))

			Expect(ut.Statistics).To(HaveLen(9))
			Expect(ut.Statistics).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("5"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Movement"),
						"Display": Equal("M"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Weapon Skill"),
						"Display": Equal("WS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("2"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Ballistic Skill"),
						"Display": Equal("BS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("3"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Strength"),
						"Display": Equal("S"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Toughness"),
						"Display": Equal("T"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("1"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Wounds"),
						"Display": Equal("W"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("3"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Initiative"),
						"Display": Equal("I"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("1"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Attacks"),
						"Display": Equal("A"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("7"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Leadership"),
						"Display": Equal("Ld"),
					}),
				})),
			))

			Expect(ut.Children).To(HaveLen(1))
			Expect(ut.Children[0].Name).To(Equal("True-horn"))
			Expect(ut.Children[0].TroopTypeName).To(Equal("Regular Infantry"))
			Expect(ut.Children[0].CompositionTypeName).To(Equal("Core"))

			Expect(ut.Children[0].Statistics).To(HaveLen(9))
			Expect(ut.Children[0].Statistics).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("5"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Movement"),
						"Display": Equal("M"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Weapon Skill"),
						"Display": Equal("WS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("2"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Ballistic Skill"),
						"Display": Equal("BS"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("3"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Strength"),
						"Display": Equal("S"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("4"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Toughness"),
						"Display": Equal("T"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("1"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Wounds"),
						"Display": Equal("W"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("3"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Initiative"),
						"Display": Equal("I"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("2"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Attacks"),
						"Display": Equal("A"),
					}),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Value": Equal("8"),
					"Statistic": MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("Leadership"),
						"Display": Equal("Ld"),
					}),
				})),
			))
		})
	})
})
