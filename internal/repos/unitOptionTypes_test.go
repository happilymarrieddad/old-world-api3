package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:CompositionTypes", func() {
	var repo repos.CompositionTypesRepo
	var gm *types.Game

	BeforeEach(func() {
		repo = gr.CompositionTypes()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should successfully create a composition type", func() {
			at, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))
			Expect(at.ID).To(Equal(at2.ID))

			at3, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Many To", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
			Expect(at3.ID).NotTo(Equal(at.ID))
			Expect(at3.ID).NotTo(Equal(at2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			at, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Single", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at.ID).NotTo(HaveLen(0))

			at2, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Many To", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at2.ID).NotTo(HaveLen(0))

			at3, err := repo.FindOrCreate(ctx, types.CreateCompositionType{Name: "Many From", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(at3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of composition types", func() {
			ats, err := repo.Find(ctx, gm.ID, 10, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(3))

			ats, err = repo.Find(ctx, gm.ID, 1, 0)
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))
		})
	})
})
