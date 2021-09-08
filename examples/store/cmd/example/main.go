package main

import (
	"context"
	"fmt"
	"log"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/examples/store/domain"
	"github.com/eleanorhealth/milo/examples/store/entityid"
	"github.com/eleanorhealth/milo/examples/store/storage"
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

	miloStore, err := milo.NewStore(db, storage.MiloEntityModelMap)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewStore(miloStore)

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

	err = store.Customers().Save(context.Background(), customer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully saved customer %s %s\n", customer.NameFirst, customer.NameLast)

	store.Transaction(context.Background(), func(txStore domain.Storer) error {
		foundCustomer, err := txStore.Customers().FindByIDForUpdate(customer.ID, false)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found customer %s %s for update\n", foundCustomer.NameFirst, foundCustomer.NameLast)

		foundCustomer.NameLast = "Doe"

		err = txStore.Customers().Save(context.Background(), foundCustomer)
		if err != nil {
			log.Fatal(err)
		}

		foundUpdatedCustomer, err := txStore.Customers().FindByIDForUpdate(foundCustomer.ID, false)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found updated customer %s %s\n", foundUpdatedCustomer.NameFirst, foundUpdatedCustomer.NameLast)

		return nil
	})
}
