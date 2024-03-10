package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:ArmyTypes", func() {
	var repo repos.ArmyTypesRepo
	var gm *types.Game

	BeforeEach(func() {
		repo = gr.ArmyTypes()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should successfully create a army type", func() {
			at, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Bretonnia", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Bretonnia", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))
			Expect(at.ID).To(Equal(at2.ID))

			at3, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Empire of Man", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
			Expect(at3.ID).NotTo(Equal(at.ID))
			Expect(at3.ID).NotTo(Equal(at2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			at, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Bretonnia", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Empire of Man", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))

			at3, err := repo.FindOrCreate(ctx, types.CreateArmyType{Name: "Warriors of Chaos", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of army types", func() {
			// We are testing the <>
			allAt, err := gr.ArmyTypes().FindOrCreate(ctx, types.CreateArmyType{Name: "All Armies", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(allAt).NotTo(BeNil())

			ats, _, err := repo.Find(ctx, gm.ID, 10, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(3))

			ats, _, err = repo.Find(ctx, gm.ID, 1, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))
		})
	})
})
