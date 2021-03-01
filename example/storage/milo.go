package storage

import (
	"reflect"

	"github.com/fterrag/milo/example/domain"
)

// MiloEntityModelMap is used by milo to map domain entities to storage models.
var MiloEntityModelMap = map[reflect.Type]reflect.Type{
	reflect.TypeOf(&domain.CarePlan{}): reflect.TypeOf(&carePlan{}),
}
