package main

import (
	"context"
	"fmt"
	"log"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/examples/simple/domain"
	"github.com/eleanorhealth/milo/examples/simple/entityid"
	"github.com/eleanorhealth/milo/examples/simple/storage"
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

	store, err := milo.NewStore(db, storage.MiloEntityModelMap)
	if err != nil {
		log.Fatal(err)
	}

	customer := &domain.Customer{
		ID:        entityid.DefaultGenerator.Generate(),
		NameFirst: "Jane",
		NameLast:  "Doe",
	}

	err = store.Save(context.Background(), customer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully saved customer %s %s\n", customer.NameFirst, customer.NameLast)

}
