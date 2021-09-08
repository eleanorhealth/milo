package storage

import (
	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/examples/simple/domain"
	"github.com/eleanorhealth/milo/examples/simple/entityid"
)

type customer struct {
	ID string `pg:"id"`

	NameFirst string `pg:"name_first"`
	NameLast  string `pg:"name_last"`
}

var _ milo.Model = (*customer)(nil)

func (c *customer) FromEntity(e interface{}) error {
	entity := e.(*domain.Customer)

	c.ID = entity.ID.String()

	c.NameFirst = entity.NameFirst
	c.NameLast = entity.NameLast

	return nil
}

func (c *customer) ToEntity() (interface{}, error) {
	entity := &domain.Customer{}

	entity.ID = entityid.ID(c.ID)

	entity.NameFirst = c.NameFirst
	entity.NameLast = c.NameLast

	return entity, nil
}
