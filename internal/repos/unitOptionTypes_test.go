package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:UnitOptionTypes", func() {
	var repo repos.UnitOptionTypesRepo
	var gm *types.Game

	BeforeEach(func() {
		repo = gr.UnitOptionTypes()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should successfully create a unit option type", func() {
			uot, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(uot.ID).NotTo(HaveLen(0))

			uot2, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(uot2.ID).NotTo(HaveLen(0))
			Expect(uot.ID).To(Equal(uot2.ID))

			uot3, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Many To", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(uot3.ID).NotTo(HaveLen(0))
			Expect(uot3.ID).NotTo(Equal(uot.ID))
			Expect(uot3.ID).NotTo(Equal(uot2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			at, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Many To", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))

			at3, err := repo.FindOrCreate(ctx, types.CreateUnitOptionType{Name: "Many Of", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of unit option types", func() {
			ats, _, err := repo.Find(ctx, gm.ID, 10, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(4)) // only 3 created above but the game ensures all are created

			ats, _, err = repo.Find(ctx, gm.ID, 1, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))
		})
	})
})
