package domain

import "github.com/eleanorhealth/milo/examples/store/entityid"

type Customer struct {
	ID entityid.ID

	NameFirst string
	NameLast  string

	Addresses []*Address
}

type Address struct {
	ID entityid.ID

	Street string
	City   string
	State  string
	Zip    string
}
