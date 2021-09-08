package domain

import (
	"context"

	"github.com/eleanorhealth/milo/example/entityid"
)

type Storer interface {
	Transaction(context.Context, func(Storer) error) error

	Customers() CustomerStorer
}

type CustomerStorer interface {
	FindAll() ([]*Customer, error)

	FindByID(id entityid.ID) (*Customer, error)
	FindByIDForUpdate(id entityid.ID, skipLocked bool) (*Customer, error)

	Save(context.Context, *Customer) error
	Delete(context.Context, *Customer) error
}
