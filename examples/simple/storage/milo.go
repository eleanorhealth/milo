package storage

import (
	"reflect"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/examples/simple/domain"
)

// MiloEntityModelMap is used by Milo to map domain entities to storage models.
var MiloEntityModelMap = milo.EntityModelMap{
	reflect.TypeOf(&domain.Customer{}): reflect.TypeOf(&customer{}),
}
