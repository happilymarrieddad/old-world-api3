package repos_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/seeder/ensurer"
	"github.com/happilymarrieddad/old-world/api3/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
)

var _ = Describe("repo:UserArmies", func() {
	Context("Simple Integ Tests", func() {
		var (
			repo repos.UserArmiesRepo
			gm   *types.Game
			usr  *types.User
			at   *types.ArmyType
			ut   *types.UnitType
			uot  *types.UnitOptionType
			uot2 *types.UnitOptionType
		)

		BeforeEach(func() {
			clearAllData()
			repo = gr.UserArmies()

			var err error
			usr, err = gr.Users().Create(ctx, types.CreateUser{
				FirstName:       "Nick",
				LastName:        "Kotenberg",
				Email:           "nick@mail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			})
			Expect(err).To(BeNil())
			Expect(usr.ID).NotTo(HaveLen(0))

			gm, err = gr.Games().FindOrCreate(ctx, "Some Game")
			Expect(err).To(BeNil())

			at, err = gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{
				GameID: gm.ID,
				Name:   "some army 1",
			})
			Expect(err).To(BeNil())

			tt, err := gr.TroopTypes().FindOrCreate(ctx, types.CreateTroopType{
				Name:   "troop type 1",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())

			ct, err := gr.CompositionTypes().FindOrCreate(ctx, types.CreateCompositionType{
				Name:   "composition type 1",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())

			stat1, err := gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{
				Name:    "Movement",
				Display: "M",
				GameID:  gm.ID,
			})
			Expect(err).To(BeNil())

			stat2, err := gr.Statistics().FindOrCreate(ctx, types.CreateStatistic{
				Name:    "Attacks",
				Display: "A",
				GameID:  gm.ID,
			})
			Expect(err).To(BeNil())

			uot, err = gr.UnitOptionTypes().FindOrCreate(ctx, types.CreateUnitOptionType{
				Name:   "Many Of",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())
			Expect(uot).NotTo(BeNil())

			uot2, err = gr.UnitOptionTypes().FindOrCreate(ctx, types.CreateUnitOptionType{
				Name:   "Single",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())
			Expect(uot2).NotTo(BeNil())

			it, err := gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{
				Name:   "Magic Weapons",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())

			it2, err := gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{
				Name:   "Magic Armour",
				GameID: gm.ID,
			})
			Expect(err).To(BeNil())

			itm, err := gr.Items().Create(ctx, types.CreateItem{
				Name:       "Some Blade",
				Points:     40,
				GameID:     gm.ID,
				ArmyTypeID: &at.ID,
				ItemTypeID: it.ID,
			})
			Expect(err).To(BeNil())

			itm2, err := gr.Items().Create(ctx, types.CreateItem{
				Name:       "Some Armour",
				Points:     15,
				GameID:     gm.ID,
				ArmyTypeID: &at.ID,
				ItemTypeID: it2.ID,
			})
			Expect(err).To(BeNil())

			ut, err = gr.UnitTypes().FindOrCreate(ctx, types.CreateUnitType{
				Name:              "unit type 1",
				GameID:            gm.ID,
				ArmyTypeID:        at.ID,
				TroopTypeID:       tt.ID,
				CompositionTypeID: ct.ID,
				PointsPerModel:    7,
				MinModels:         5,
				MaxModels:         100,
				Statistics: []*types.CreateUnitStatistics{
					{Display: stat1.Display, Value: "4"},
					{Display: stat2.Display, Value: "1"},
				},
				UnitOptions: []*types.UnitTypeOption{
					{
						UnitOptionTypeID: uot.ID,
						Txt:              "Pick any of the following",
						Points:           0,
						PerModel:         false,
						MaxPoints:        50,
						Items:            []*types.Item{itm, itm2},
					},
					{
						UnitOptionTypeID: uot2.ID,
						Txt:              "May add a shield",
						Points:           2,
						PerModel:         true,
						MaxPoints:        0,
					},
				},
			})
			Expect(err).To(BeNil())
		})

		It("should successfully create an army and verify it's contents", func() {
			eua, err := repo.Create(ctx, types.CreateUserArmy{
				Name:       "Army 1",
				UserID:     usr.ID,
				GameID:     gm.ID,
				ArmyTypeID: at.ID,
				Points:     2000,
			})
			Expect(err).To(BeNil())
			Expect(eua.Name).To(Equal("Army 1"))
			Expect(eua.Units).To(HaveLen(0))

			euaUnits, err := repo.AddUnits(ctx, eua.ID, &types.CreateUserArmyUnit{
				UserArmyID: eua.ID,
				UnitTypeID: ut.ID,
			}, &types.CreateUserArmyUnit{
				UserArmyID: eua.ID,
				UnitTypeID: ut.ID,
			})
			Expect(err).To(BeNil())
			Expect(euaUnits).To(HaveLen(2))

			existingUa, err := repo.Get(ctx, usr.ID, eua.ID)
			Expect(err).To(BeNil())
			Expect(existingUa.Name).To(Equal("Army 1"))
			Expect(existingUa.Units).To(HaveLen(2))

			Expect(repo.RemoveUnits(ctx, existingUa.ID, euaUnits[0].ID)).To(Succeed())

			existingUa, err = repo.Get(ctx, usr.ID, eua.ID)
			Expect(err).To(BeNil())
			Expect(existingUa.Name).To(Equal("Army 1"))
			Expect(existingUa.Units).To(HaveLen(1))

			existingUserArmyUnit := existingUa.Units[0]
			Expect(existingUserArmyUnit).NotTo(BeNil())

			// TODO: add option values to the unit
			bts, err := json.Marshal(existingUserArmyUnit)
			Expect(err).To(BeNil())
			fmt.Println(string(bts))
		})
	})

	Context("Full Tests for User Armies", func() {
		var repo repos.UserArmiesRepo
		var gm *types.Game
		var usr *types.User
		var at *types.ArmyType
		var at2 *types.ArmyType
		var runOnce bool

		BeforeEach(func() {
			if os.Getenv("DEBUG") != "true" {
				Skip("skipping heavy testing")
			}

			repo = gr.UserArmies()

			// Only need the data once
			if !runOnce {
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

				usr, err = gr.Users().GetByEmail(ctx, "nick@mail.com")
				if err != nil || usr == nil {
					usr, err = gr.Users().Create(ctx, types.CreateUser{
						FirstName:       "Nick",
						LastName:        "Kotenberg",
						Email:           "nick@mail.com",
						Password:        "1234",
						PasswordConfirm: "1234",
					})
					Expect(err).To(BeNil())
					Expect(usr.ID).NotTo(HaveLen(0))
				}

				runOnce = true
			}

			gms, err := gr.Games().Find(ctx, 1, 0)
			Expect(err).To(BeNil())
			Expect(gms).NotTo(HaveLen(0))
			gm = gms[0]

			ats, _, err := gr.ArmyTypes().Find(ctx, gm.ID, 2, 0)
			Expect(err).To(BeNil())
			Expect(ats).NotTo(HaveLen(0))
			at = ats[0]
			at2 = ats[1]
		})

		Context("Create", func() {
			It("should return an err when an invalid user army is passed in", func() {
				ua, err := repo.Create(ctx, types.CreateUserArmy{})
				Expect(err).NotTo(BeNil())
				Expect(ua).To(BeNil())
			})

			It("should successfully create a basic user army", func() {
				ua, err := repo.Create(ctx, types.CreateUserArmy{
					Name:       "Some usr army",
					UserID:     usr.ID,
					GameID:     gm.ID,
					ArmyTypeID: at.ID,
					Points:     2000,
				})
				Expect(err).To(BeNil())
				Expect(ua).NotTo(BeNil())
				Expect(ua).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":          Equal("Some usr army"),
					"Points":        BeNumerically("==", 2000),
					"UserFirstName": Equal(usr.FirstName),
					"UserLastName":  Equal(usr.LastName),
					"UserEmail":     Equal(usr.Email),
					"GameName":      Equal(gm.Name),
					"ArmyTypeName":  Equal(at.Name),
				})))
			})
		})

		Context("Find/Get", func() {
			var ua1 *types.UserArmy
			var ua2 *types.UserArmy
			BeforeEach(func() {
				db.WriteData(ctx, driver, func(tx neo4j.ManagedTransaction) (any, error) {
					_, err := tx.Run(ctx, `
				MATCH (ua:UserArmy) DETACH DELETE ua
				`, map[string]any{})
					Expect(err).To(BeNil())
					return nil, nil
				})

				ua, err := repo.Create(ctx, types.CreateUserArmy{
					Name:       "Some usr army",
					UserID:     usr.ID,
					GameID:     gm.ID,
					ArmyTypeID: at.ID,
					Points:     2000,
				})
				Expect(err).To(BeNil())
				Expect(ua).NotTo(BeNil())
				Expect(ua).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":          Equal("Some usr army"),
					"Points":        BeNumerically("==", 2000),
					"UserFirstName": Equal(usr.FirstName),
					"UserLastName":  Equal(usr.LastName),
					"UserEmail":     Equal(usr.Email),
					"GameName":      Equal(gm.Name),
					"ArmyTypeName":  Equal(at.Name),
				})))
				ua1 = ua
				Expect(ua1).NotTo(BeNil())
				ua2, err = repo.Create(ctx, types.CreateUserArmy{
					Name:       "Some usr army 2",
					UserID:     usr.ID,
					GameID:     gm.ID,
					ArmyTypeID: at2.ID,
					Points:     2000,
				})
				Expect(err).To(BeNil())
				Expect(ua2).NotTo(BeNil())
				Expect(ua2).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":          Equal("Some usr army 2"),
					"Points":        BeNumerically("==", 2000),
					"UserFirstName": Equal(usr.FirstName),
					"UserLastName":  Equal(usr.LastName),
					"UserEmail":     Equal(usr.Email),
					"GameName":      Equal(gm.Name),
					"ArmyTypeName":  Equal(at2.Name),
				})))
			})

			It("should return the existing user armies", func() {
				uas, _, err := repo.Find(ctx, uuid.New().String(), nil)
				Expect(err).To(BeNil())
				Expect(uas).To(HaveLen(0))

				uas, _, err = repo.Find(ctx, usr.ID, nil)
				Expect(err).To(BeNil())
				Expect(uas).To(HaveLen(2))

				for _, ua := range uas {
					switch ua.Name {
					case "Some usr army":
						Expect(ua).To(BeEquivalentTo(ua1))
					case "Some usr army 2":
						Expect(ua).To(BeEquivalentTo(ua2))
					default:
						spew.Dump(ua)
						Expect(false).To(BeTrue(), "invalid user army passed in")
					}
				}

				existingUserArmy, err := repo.Get(ctx, usr.ID, ua1.ID)
				Expect(err).To(BeNil())
				Expect(existingUserArmy).To(BeEquivalentTo(ua1))

				existingUserArmy, err = repo.Get(ctx, usr.ID, uuid.New().String())
				Expect(err).NotTo(BeNil())
				Expect(types.IsNotFoundError(err)).To(BeTrue())
				Expect(existingUserArmy).To(BeNil())
			})
		})

		Context("AddUnits", func() {
			var ua *types.UserArmy
			BeforeEach(func() {
				db.WriteData(ctx, driver, func(tx neo4j.ManagedTransaction) (any, error) {
					_, err := tx.Run(ctx, `
				MATCH (ua:UserArmy) DETACH DELETE ua
				`, map[string]any{})
					Expect(err).To(BeNil())
					return nil, nil
				})

				var err error
				ua, err = repo.Create(ctx, types.CreateUserArmy{
					Name:       "Some usr army",
					UserID:     usr.ID,
					GameID:     gm.ID,
					ArmyTypeID: at.ID,
					Points:     2000,
				})
				Expect(err).To(BeNil())
				Expect(ua).NotTo(BeNil())
				Expect(ua).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":          Equal("Some usr army"),
					"Points":        BeNumerically("==", 2000),
					"UserFirstName": Equal(usr.FirstName),
					"UserLastName":  Equal(usr.LastName),
					"UserEmail":     Equal(usr.Email),
					"GameName":      Equal(gm.Name),
					"ArmyTypeName":  Equal(at.Name),
				})))
			})

			It("should add units to the user army", func() {
				uts, _, err := gr.UnitTypes().Find(ctx, &repos.FindUnitTypesOpts{
					ArmyTypeID: ua.ArmyTypeID,
					Names:      []string{"Beastlord", "Gor Herds"},
				})
				Expect(err).To(BeNil())
				Expect(len(uts)).To(BeNumerically("==", 2))

				existingArmyUnits, err := repo.AddUnits(ctx, ua.ID, &types.CreateUserArmyUnit{
					UserArmyID: ua.ID,
					UnitTypeID: uts[0].ID,
				}, &types.CreateUserArmyUnit{
					UserArmyID: ua.ID,
					UnitTypeID: uts[1].ID,
				})
				Expect(err).To(BeNil())
				Expect(existingArmyUnits).To(HaveLen(2))

				existingUa, err := repo.Get(ctx, usr.ID, ua.ID)
				Expect(err).To(BeNil())

				for _, eau := range existingArmyUnits {
					Expect(eau.UserArmyID).To(Equal(existingUa.ID))
				}

				Expect(existingUa).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ID":     Equal(ua.ID),
					"Name":   Equal(ua.Name),
					"Units":  HaveLen(2),
					"Points": Equal(2000),
				})))
				var points int
				for _, unit := range existingUa.Units {
					points += unit.Points
				}
				// kind of a quick check if the units have points values
				Expect(points).To(Equal(115 + (5 * 7)))

				// TODO: ensure units have everything

				Expect(repo.RemoveUnits(ctx, existingUa.ID, existingUa.Units[0].ID)).To(Succeed())

				existingUa2, err := repo.Get(ctx, usr.ID, ua.ID)
				Expect(err).To(BeNil())
				Expect(existingUa2.Units).To(HaveLen(len(existingUa.Units) - 1))
			})
		})
	})
})
