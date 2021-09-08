package storage

import (
	"context"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/examples/store/domain"
)

type Store struct {
	miloStore milo.Storer

	customers *CustomerStore
}

var _ domain.Storer = (*Store)(nil)

func NewStore(miloStore milo.Storer) *Store {
	return &Store{
		miloStore: miloStore,

		customers: NewCustomerStore(miloStore),
	}
}

func (s *Store) Transaction(ctx context.Context, fn func(domain.Storer) error) error {
	return s.miloStore.Transaction(ctx, func(miloStore milo.Storer) error {
		return fn(NewStore(miloStore))
	})
}

func (s *Store) Customers() domain.CustomerStorer {
	return s.customers
}
