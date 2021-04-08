package milo

import "reflect"

type EntityModelMap map[reflect.Type]reflect.Type

type Model interface {
	FromEntity(interface{}) error
	ToEntity() (interface{}, error)
}

type Field string
