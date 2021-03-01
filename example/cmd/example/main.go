package main

import (
	"log"
	"time"

	"github.com/fterrag/milo"
	"github.com/fterrag/milo/example/domain"
	"github.com/fterrag/milo/example/entityid"
	"github.com/fterrag/milo/example/storage"
	"github.com/go-pg/pg/v10"
)

func main() {
	// See docker-compose.yml
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:8200",
		User:     "postgres",
		Password: "password",
		Database: "milo",
	})
	defer db.Close()

	err := storage.CreateSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	store := milo.NewStore(db, storage.MiloEntityModelMap)

	carePlan := &domain.CarePlan{
		ID: entityid.DefaultGenerator.Generate(),

		Goals: []*domain.Goal{
			{
				ID:    entityid.DefaultGenerator.Generate(),
				Title: "Exercise",
				Body:  "Exercise at least 3 times a week.",
			},
		},

		YearOfCare: time.Now(),
	}

	err = store.Save(carePlan)
	if err != nil {
		log.Fatal(err)
	}
}
