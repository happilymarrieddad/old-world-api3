package repos_test

import (
	"context"
	"testing"

	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var gr repos.GlobalRepo
var driver neo4j.DriverWithContext
var ctx context.Context

var _ = BeforeSuite(func() {
	ctx = context.Background()

	var err error
	driver, err = db.NewDB()
	Expect(err).To(BeNil())

	gr, err = repos.NewGlobalRepo(driver)
	Expect(err).To(BeNil())
	Expect(gr).NotTo(BeNil())
})

func clearAllData() {
	db.WriteData(ctx, driver, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
		MATCH (n) DETACH DELETE n
		`, map[string]any{})
		Expect(err).To(BeNil())
		return nil, nil
	})
}

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
}
