package repos_test

import (
	"github.com/google/uuid"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("repo:Users", func() {
	var repo repos.UsersRepo

	BeforeEach(func() {
		repo = gr.Users()

		clearAllData()
	})

	Context("Create", func() {
		It("should successfully create a user", func() {
			nu, err := repo.Create(ctx, types.CreateUser{
				FirstName:       "Nick",
				LastName:        "Kotenberg",
				Email:           "nick@mail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			})
			Expect(err).To(BeNil())
			Expect(nu.ID).NotTo(HaveLen(0))

			// _, err = repo.Create(ctx, types.CreateUser{
			// 	FirstName:       "Nick",
			// 	LastName:        "Kotenberg",
			// 	Email:           "nick@mail.com",
			// 	Password:        "1234",
			// 	PasswordConfirm: "1234",
			// })
			// Expect(err).NotTo(BeNil())
			// fmt.Println(err)
		})
	})

	Context("Get(id|email)", func() {
		var id string
		BeforeEach(func() {
			nu, err := repo.Create(ctx, types.CreateUser{
				FirstName:       "Nick",
				LastName:        "Kotenberg",
				Email:           "nick@mail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			})
			Expect(err).To(BeNil())
			id = nu.ID
		})

		It("should return an err when not found", func() {
			eu, err := repo.GetByID(ctx, uuid.New().String())
			Expect(types.IsNotFoundError(err)).To(BeTrue())
			Expect(eu).To(BeNil())
		})

		It("should be found by ID", func() {
			eu, err := repo.GetByID(ctx, id)
			Expect(err).To(BeNil())
			Expect(eu.FirstName).To(Equal("Nick"))
		})

		It("should be found by Email", func() {
			eu, err := repo.GetByEmail(ctx, "nick@mail.com")
			Expect(err).To(BeNil())
			Expect(eu.ID).To(Equal(id))
		})
	})

	Context("Delete", func() {
		var id string
		BeforeEach(func() {
			nu, err := repo.Create(ctx, types.CreateUser{
				FirstName:       "Nick",
				LastName:        "Kotenberg",
				Email:           "nick@mail.com",
				Password:        "1234",
				PasswordConfirm: "1234",
			})
			Expect(err).To(BeNil())
			id = nu.ID
		})

		It("should delete the existing user", func() {
			Expect(repo.Delete(ctx, id)).To(Succeed())

			eu, err := repo.GetByID(ctx, id)
			Expect(types.IsNotFoundError(err)).To(BeTrue())
			Expect(eu).To(BeNil())
		})
	})
})
