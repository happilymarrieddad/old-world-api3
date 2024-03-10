package repos_test

import (
	"encoding/json"
	"os"

	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/internal/utils"
	"github.com/happilymarrieddad/old-world/api3/seeder/ensurer"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
)

var _ = Describe("repos:Items", func() {
	var (
		repo repos.ItemsRepo
		gm   *types.Game
		gm2  *types.Game
		at   *types.ArmyType
		at2  *types.ArmyType
		it   *types.ItemType
		it2  *types.ItemType

		itemList []*types.Item

		createItems = func() {
			ei1, err := repo.Create(ctx, types.CreateItem{
				Name:       "Ogre Blade",
				Points:     65,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
			})
			Expect(err).To(BeNil())
			Expect(ei1.ID).NotTo(HaveLen(0))
			ei2, err := repo.Create(ctx, types.CreateItem{
				Name:       "Sword of Battle",
				Points:     60,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
			})
			Expect(err).To(BeNil())
			Expect(ei2.ID).NotTo(HaveLen(0))

			ei3, err := repo.Create(ctx, types.CreateItem{
				Name:       "Enchanted Shield",
				Points:     10,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
				ArmyTypeID: utils.Ref(at.ID),
			})
			Expect(err).To(BeNil())
			Expect(ei3.ID).NotTo(HaveLen(0))
			ei4, err := repo.Create(ctx, types.CreateItem{
				Name:       "Dawnstone",
				Points:     35,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
				ArmyTypeID: utils.Ref(at.ID),
			})
			Expect(err).To(BeNil())
			Expect(ei4.ID).NotTo(HaveLen(0))

			ei5, err := repo.Create(ctx, types.CreateItem{
				Name:       "Dawnstone",
				Points:     30,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
				ArmyTypeID: utils.Ref(at2.ID),
			})
			Expect(err).To(BeNil())
			Expect(ei5.ID).NotTo(HaveLen(0))
			ei6, err := repo.Create(ctx, types.CreateItem{
				Name:       "Charmed Shield",
				Points:     5,
				GameID:     gm.ID,
				ItemTypeID: it.ID,
				ArmyTypeID: utils.Ref(at2.ID),
			})
			Expect(err).To(BeNil())
			Expect(ei6.ID).NotTo(HaveLen(0))

			ei7, err := repo.Create(ctx, types.CreateItem{
				Name:       "Ogre Blade",
				Points:     65,
				GameID:     gm2.ID,
				ItemTypeID: it2.ID,
			})
			Expect(err).To(BeNil())
			Expect(ei7.ID).NotTo(HaveLen(0))

			itemList = []*types.Item{ei1, ei2, ei3, ei4, ei5, ei6, ei7}
			Expect(itemList).To(HaveLen(7))
		}
	)

	BeforeEach(func() {
		repo = gr.Items()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())

		gm2, err = gr.Games().FindOrCreate(ctx, "Warhammer 40k 10th Edition")
		Expect(err).To(BeNil())

		at, err = gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "Bretonnia", GameID: gm.ID})
		Expect(err).To(BeNil())

		at2, err = gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "Beastmen Herds", GameID: gm.ID})
		Expect(err).To(BeNil())

		it, err = gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{Name: "Magic Weapons", GameID: gm.ID})
		Expect(err).To(BeNil())

		it2, err = gr.ItemTypes().FindOrCreate(ctx, types.CreateItemType{Name: "Magic Armour", GameID: gm.ID})
		Expect(err).To(BeNil())
	})

	Context("Create", func() {
		It("should successfully create list of new items", func() {
			createItems()
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			createItems()
			createItems()
		})

		It("should return the num of items per game/army type", func() {
			oldWorldItems, err := repo.Find(ctx, &repos.FindItemsOpts{GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(oldWorldItems).To(HaveLen(2))

			fourtyKItems, err := repo.Find(ctx, &repos.FindItemsOpts{GameID: gm2.ID})
			Expect(err).To(BeNil())
			Expect(fourtyKItems).To(HaveLen(1))

			bretonnianItems, err := repo.Find(ctx, &repos.FindItemsOpts{
				GameID:     gm.ID,
				ArmyTypeID: &at.ID,
			})
			Expect(err).To(BeNil())
			Expect(bretonnianItems).To(HaveLen(4))
			Expect(bretonnianItems).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Ogre Blade"),
					"Points": BeNumerically("==", 65),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Sword of Battle"),
					"Points": BeNumerically("==", 60),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Enchanted Shield"),
					"Points": BeNumerically("==", 10),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Dawnstone"),
					"Points": BeNumerically("==", 35),
					"GameID": Equal(gm.ID),
				})),
			))

			beastmenItems, err := repo.Find(ctx, &repos.FindItemsOpts{
				GameID:     gm.ID,
				ArmyTypeID: &at2.ID,
			})
			Expect(err).To(BeNil())
			Expect(beastmenItems).To(HaveLen(4))
			Expect(beastmenItems).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Ogre Blade"),
					"Points": BeNumerically("==", 65),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Sword of Battle"),
					"Points": BeNumerically("==", 60),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Charmed Shield"),
					"Points": BeNumerically("==", 5),
					"GameID": Equal(gm.ID),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":   Equal("Dawnstone"),
					"Points": BeNumerically("==", 30), // This is a test to ensure uniqueness
					"GameID": Equal(gm.ID),
				})),
			))

			bitm := beastmenItems[0]
			existingItem, err := repo.Get(ctx, bitm.ID, bitm.GameID)
			Expect(err).To(BeNil())
			Expect(existingItem).To(BeEquivalentTo(bitm))
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

		It("should return items from a specific army type", func() {
			htgm, err := gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
			Expect(err).To(BeNil())

			htat, err := gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "Beastmen Brayherds", GameID: htgm.ID})
			Expect(err).To(BeNil())

			itms, err := repo.Find(ctx, &repos.FindItemsOpts{GameID: htgm.ID})
			Expect(err).To(BeNil())
			Expect(itms).To(HaveLen(48))
			Expect(itms).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Ogre Blade"),
					"Points":       BeNumerically("==", 65),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   BeNil(),
				})),
			))
			Expect(itms).NotTo(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Primeval Club"),
					"Points":       BeNumerically("==", 60),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   PointTo(Equal(htat.ID)),
				})),
			))
			Expect(itms[0]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":         Equal("Ogre Blade"),
				"Points":       BeNumerically("==", 65),
				"ItemTypeName": Equal("Magic Weapons"),
				"ArmyTypeID":   BeNil(),
			})))
			Expect(itms[len(itms)-1]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":         Equal("Earthing Rod"),
				"Points":       BeNumerically("==", 5),
				"ItemTypeName": Equal("Arcane Items"),
				"ArmyTypeID":   BeNil(),
			})))

			itms, err = repo.Find(ctx, &repos.FindItemsOpts{GameID: htgm.ID, ArmyTypeID: &htat.ID})
			Expect(err).To(BeNil())
			Expect(itms).To(HaveLen(81))
			Expect(itms).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Ogre Blade"),
					"Points":       BeNumerically("==", 65),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   BeNil(),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Primeval Club"),
					"Points":       BeNumerically("==", 60),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   PointTo(Equal(htat.ID)),
				})),
			))

			its, _, err := gr.ItemTypes().Find(ctx, &repos.FindItemTypeOpts{
				GameID: htgm.ID,
				Name:   []string{"Magic Weapons", "Magic Armour"},
			})
			Expect(err).To(BeNil())
			Expect(its).To(HaveLen(2))

			itms, err = repo.Find(ctx, &repos.FindItemsOpts{
				GameID: htgm.ID, ArmyTypeID: &htat.ID,
				ItemTypeIDs: []string{its[0].ID, its[1].ID},
			})
			Expect(err).To(BeNil())
			Expect(itms).To(HaveLen(28))
			Expect(itms).To(ContainElements(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Ogre Blade"),
					"Points":       BeNumerically("==", 65),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   BeNil(),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("Primeval Club"),
					"Points":       BeNumerically("==", 60),
					"ItemTypeName": Equal("Magic Weapons"),
					"ArmyTypeID":   PointTo(Equal(htat.ID)),
				})),
			))
			Expect(itms[len(itms)-1]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":         Equal("Charmed Shield"),
				"Points":       BeNumerically("==", 5),
				"ItemTypeName": Equal("Magic Armour"),
				"ArmyTypeID":   BeNil(),
			})))
		})
	})
})
