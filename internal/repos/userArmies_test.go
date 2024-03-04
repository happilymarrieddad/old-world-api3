package repos_test

import (
	"encoding/json"
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

			usr, err = gr.Users().Create(ctx, types.CreateUser{
				FirstName:       "Nick",
				LastName:        "Kotenberg",
				Email:           "nick@mail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			})
			Expect(err).To(BeNil())
			Expect(usr.ID).NotTo(HaveLen(0))

			runOnce = true
		}

		gms, err := gr.Games().Find(ctx, 1, 0)
		Expect(err).To(BeNil())
		Expect(gms).NotTo(HaveLen(0))
		gm = gms[0]

		ats, err := gr.ArmyTypes().Find(ctx, gm.ID, 2, 0)
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
			uas, err := repo.Find(ctx, uuid.New().String(), nil)
			Expect(err).To(BeNil())
			Expect(uas).To(HaveLen(0))

			uas, err = repo.Find(ctx, usr.ID, nil)
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
			uts, err := gr.UnitTypes().Find(ctx, &repos.FindUnitTypesOpts{
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
				"ID":    Equal(ua.ID),
				"Name":  Equal(ua.Name),
				"Units": HaveLen(2),
			})))
			// TODO: ensure units have everything
		})
	})
})
