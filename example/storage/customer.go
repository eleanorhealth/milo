package storage

import (
	"context"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/example/domain"
	"github.com/eleanorhealth/milo/example/entityid"
	"github.com/pkg/errors"
)

type customer struct {
	ID string `pg:"id"`

	NameFirst string `pg:"name_first"`
	NameLast  string `pg:"name_last"`

	Addresses []*address `pg:"rel:has-many"`
}

var _ milo.Model = (*customer)(nil)

type address struct {
	ID         string `pg:"id"`
	CustomerID string `pg:"customer_id"`

	Street string `pg:"street"`
	City   string `pg:"city"`
	State  string `pg:"state"`
	Zip    string `pg:"zip"`
}

func (c *customer) FromEntity(e interface{}) error {
	entity := e.(*domain.Customer)

	c.ID = entity.ID.String()

	c.NameFirst = entity.NameFirst
	c.NameLast = entity.NameLast

	for _, a := range entity.Addresses {
		c.Addresses = append(c.Addresses, &address{
			ID:         a.ID.String(),
			CustomerID: c.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return nil
}

func (c *customer) ToEntity() (interface{}, error) {
	entity := &domain.Customer{}

	entity.ID = entityid.ID(c.ID)

	entity.NameFirst = c.NameFirst
	entity.NameLast = c.NameLast

	for _, a := range c.Addresses {
		entity.Addresses = append(entity.Addresses, &domain.Address{
			ID: entityid.ID(a.ID),

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return entity, nil
}

type CustomerStore struct {
	miloStore milo.Storer
}

var _ domain.CustomerStorer = (*CustomerStore)(nil)

func NewCustomerStore(miloStore milo.Storer) *CustomerStore {
	return &CustomerStore{
		miloStore: miloStore,
	}
}

func (s *CustomerStore) FindAll() ([]*domain.Customer, error) {
	entities := []*domain.Customer{}
	err := s.miloStore.FindAll(&entities)
	if err != nil {
		return nil, errors.Wrap(err, "finding all entities")
	}

	return entities, nil
}

func (s *CustomerStore) FindByID(id entityid.ID) (*domain.Customer, error) {
	entity := &domain.Customer{}
	err := s.miloStore.FindByID(entity, id)
	if err != nil {
		return nil, errors.Wrap(err, "finding entity by ID")
	}

	return entity, nil
}

func (s *CustomerStore) FindByIDForUpdate(id entityid.ID, skipLocked bool) (*domain.Customer, error) {
	entity := &domain.Customer{}
	err := s.miloStore.FindByIDForUpdate(entity, id, skipLocked)
	if err != nil {
		return nil, errors.Wrap(err, "finding entity by ID for update")
	}

	return entity, nil
}

func (s *CustomerStore) Save(ctx context.Context, entity *domain.Customer) error {
	err := s.miloStore.Save(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "saving entity")
	}

	return nil
}

func (s *CustomerStore) Delete(ctx context.Context, entity *domain.Customer) error {
	err := s.miloStore.Delete(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "deleting entity")
	}

	return nil
}
