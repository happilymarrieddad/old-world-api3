package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:ItemTypes", func() {
	var repo repos.ItemTypesRepo
	var gm *types.Game

	BeforeEach(func() {
		repo = gr.ItemTypes()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should successfully create a item type", func() {
			at, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Magic Weapons", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Magic Weapons", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))
			Expect(at.ID).To(Equal(at2.ID))

			at3, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Melee Weapons", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
			Expect(at3.ID).NotTo(Equal(at.ID))
			Expect(at3.ID).NotTo(Equal(at2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			at, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Magic Weapons", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Melee Weapons", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))

			at3, err := repo.FindOrCreate(ctx, types.CreateItemType{Name: "Magic Standards", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of item types", func() {
			ats, _, err := repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Limit: 10})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(3))

			ats, _, err = repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Limit: 1})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))

			ats, _, err = repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Name: []string{"Magic Weapons"}})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))

			ats, _, err = repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Name: []string{"Magic Standards"}})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))

			ats, _, err = repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Name: []string{"Magic Weapons", "Magic Standards"}})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(2))

			ats, _, err = repo.Find(ctx, &repos.FindItemTypeOpts{GameID: gm.ID, Name: []string{
				"Magic Weapons", "Magic Standards", "garbage"}})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(2))
		})
	})
})
