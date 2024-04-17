package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:Games", func() {
	var repo repos.GamesRepo

	BeforeEach(func() {
		repo = gr.Games()

		clearAllData()
	})

	Context("FindOrCreate", func() {
		It("should successfully create a user", func() {
			gm, err := repo.FindOrCreate(ctx, "Some Game")
			Expect(err).To(BeNil())
			Expect(gm.ID).NotTo(HaveLen(0))

			gm2, err := repo.FindOrCreate(ctx, "Some Game")
			Expect(err).To(BeNil())
			Expect(gm2.ID).NotTo(HaveLen(0))
			Expect(gm.ID).To(Equal(gm2.ID))

			gm3, err := repo.FindOrCreate(ctx, "Some Game 2")
			Expect(err).To(BeNil())
			Expect(gm3.ID).NotTo(HaveLen(0))
			Expect(gm3.ID).NotTo(Equal(gm.ID))
			Expect(gm3.ID).NotTo(Equal(gm2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			gm, err := repo.FindOrCreate(ctx, "Some Game")
			Expect(err).To(BeNil())
			Expect(gm.ID).NotTo(HaveLen(0))

			gm2, err := repo.FindOrCreate(ctx, "Some Game 2")
			Expect(err).To(BeNil())
			Expect(gm2.ID).NotTo(HaveLen(0))

			gm3, err := repo.FindOrCreate(ctx, "Some Game 3")
			Expect(err).To(BeNil())
			Expect(gm3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of games", func() {
			gms, err := repo.Find(ctx, &repos.FindGameOpts{Limit: 10})
			Expect(err).To(BeNil())
			Expect(gms).To(HaveLen(3))

			gms, err = repo.Find(ctx, &repos.FindGameOpts{Limit: 1})
			Expect(err).To(BeNil())
			Expect(gms).To(HaveLen(1))
		})
	})
})
