package storage

import (
	"reflect"

	"github.com/eleanorhealth/milo"
	"github.com/eleanorhealth/milo/example/domain"
)

// MiloEntityModelMap is used by Milo to map domain entities to storage models.
var MiloEntityModelMap = milo.EntityModelMap{
	reflect.TypeOf(&domain.Customer{}): milo.ModelConfig{
		Model:          reflect.TypeOf(&customer{}),
		FieldColumnMap: milo.FieldColumnMap{},
	},
}
