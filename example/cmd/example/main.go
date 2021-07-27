package main

import (
	"context"
	"fmt"
	"log"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/example/domain"
	"github.com/eleanorhealth/milo/example/entityid"
	"github.com/eleanorhealth/milo/example/storage"
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

	customer := &domain.Customer{
		ID: entityid.DefaultGenerator.Generate(),

		NameFirst: "John",
		NameLast:  "Smith",

		Addresses: []*domain.Address{
			{
				ID: entityid.DefaultGenerator.Generate(),

				Street: "1 City Hall Square #500",
				City:   "Boston",
				State:  "MA",
				Zip:    "02201",
			},
		},
	}

	err = store.Save(context.Background(), customer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully saved customer %s %s\n", customer.NameFirst, customer.NameLast)
}
