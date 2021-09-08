package domain

import "github.com/eleanorhealth/milo/examples/simple/entityid"

type Customer struct {
	ID entityid.ID

	NameFirst string
	NameLast  string
}
