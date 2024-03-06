package main

import (
	"log"

	"github.com/happilymarrieddad/old-world/api3/internal/api"
	"github.com/happilymarrieddad/old-world/api3/internal/db"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
)

func main() {
	driver, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	gr, err := repos.NewGlobalRepo(driver)
	if err != nil {
		log.Fatal(err)
	}

	api.Run(gr)
}
