package domain

import (
	"time"

	"github.com/eleanorhealth/milo/example/entityid"
)

type CarePlan struct {
	ID entityid.ID

	Goals []*Goal

	YearOfCare time.Time
}

type Goal struct {
	ID entityid.ID

	Title string
	Body  string
}
