package storage

import (
	"time"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/example/domain"
	"github.com/eleanorhealth/milo/example/entityid"
)

type carePlan struct {
	ID string `pg:"id"`

	Goals      []*goal   `pg:"rel:has-many"`
	YearOfCare time.Time `pg:"year_of_care"`
}

var _ milo.Model = (*carePlan)(nil)

type goal struct {
	ID         string `pg:"id"`
	CarePlanID string `pg:"care_plan_id"`

	Title string `pg:"title"`
	Body  string `pg:"body"`
}

func (c *carePlan) FromEntity(e interface{}) error {
	entity := e.(*domain.CarePlan)

	c.ID = entity.ID.String()

	for _, g := range entity.Goals {
		c.Goals = append(c.Goals, &goal{
			ID:         g.ID.String(),
			CarePlanID: c.ID,

			Title: g.Title,
			Body:  g.Body,
		})
	}

	c.YearOfCare = entity.YearOfCare

	return nil
}

func (c *carePlan) ToEntity() (interface{}, error) {
	entity := &domain.CarePlan{}

	entity.ID = entityid.ID(c.ID)

	for _, g := range c.Goals {
		entity.Goals = append(entity.Goals, &domain.Goal{
			ID:    entityid.ID(g.ID),
			Title: g.Title,
			Body:  g.Body,
		})
	}

	entity.YearOfCare = c.YearOfCare

	return entity, nil
}
