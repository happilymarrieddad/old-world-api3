package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	"github.com/happilymarrieddad/old-world/api3/seeder/ensurer"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pkg/errors"
)

func clearAllData(ctx context.Context, driver neo4j.DriverWithContext) {
	db.WriteData(ctx, driver, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
		MATCH (n) DETACH DELETE n
		`, map[string]any{})
		if err != nil {
			panic(err)
		}
		return nil, nil
	})
}

func main() {
	ctx := context.Background()

	driver, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	gr, err := repos.NewGlobalRepo(driver)
	if err != nil {
		panic(err)
	}

	bts, err := os.ReadFile("internal/repos/data/example.json")
	if err != nil {
		panic(err)
	}

	ad := make(ensurer.Games)
	if err := json.Unmarshal(bts, &ad); err != nil {
		panic(errors.WithMessage(err, "unable to marshal data"))
	}

	//clearAllData(ctx, driver)

	if err := ensurer.EnsureData(ctx, gr, ad); err != nil {
		panic(err)
	}
}
