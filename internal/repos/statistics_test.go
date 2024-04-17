package repos_test

import (
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:Statistics", func() {
	var repo repos.StatisticsRepo
	var gm *types.Game

	BeforeEach(func() {
		repo = gr.Statistics()

		clearAllData()

		var err error
		gm, err = gr.Games().FindOrCreate(ctx, "Old World 1st Edition")
		Expect(err).To(BeNil())
	})

	Context("FindOrCreate", func() {
		It("should successfully create a statistic", func() {
			stat, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Movement", Display: "M", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat.ID).NotTo(HaveLen(0))

			stat2, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Movement", Display: "M", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat2.ID).NotTo(HaveLen(0))
			Expect(stat.ID).To(Equal(stat2.ID))

			stat3, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Weapon Skill", Display: "WS", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat3.ID).NotTo(HaveLen(0))
			Expect(stat3.ID).NotTo(Equal(stat.ID))
			Expect(stat3.ID).NotTo(Equal(stat2.ID))
		})
	})

	Context("Find", func() {
		BeforeEach(func() {
			stat, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Movement", Display: "M", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat.ID).NotTo(HaveLen(0))

			stat2, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Weapon Skill", Display: "WS", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat2.ID).NotTo(HaveLen(0))

			stat3, err := repo.FindOrCreate(ctx, types.CreateStatistic{Name: "Strength", Display: "S", GameID: gm.ID})
			Expect(err).To(BeNil())
			Expect(stat3.ID).NotTo(HaveLen(0))
		})

		It("should return the list of games", func() {
			ats, _, err := repo.Find(ctx, &repos.FindStatisticsOpts{GameIDs: []string{gm.ID}, Limit: 10})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(3))

			ats, _, err = repo.Find(ctx, &repos.FindStatisticsOpts{GameIDs: []string{gm.ID}, Limit: 1})
			Expect(err).To(BeNil())
			Expect(ats).To(HaveLen(1))
		})
	})
})
