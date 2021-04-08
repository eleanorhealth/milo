package storage

import (
	"reflect"

	"github.com/eleanorhealth/milo/example/domain"
)

// MiloEntityModelMap is used by Milo to map domain entities to storage models.
var MiloEntityModelMap = map[reflect.Type]reflect.Type{
	reflect.TypeOf(&domain.Customer{}): reflect.TypeOf(&customer{}),
}
