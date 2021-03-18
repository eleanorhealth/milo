package storage

import (
	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/example/domain"
	"github.com/eleanorhealth/milo/example/entityid"
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
