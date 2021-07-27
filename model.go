package milo

import "context"

type Model interface {
	FromEntity(interface{}) error
	ToEntity() (interface{}, error)
}

type Hook interface {
	BeforeSave(ctx context.Context, store Storer, entity interface{}) error
	BeforeDelete(ctx context.Context, store Storer, entity interface{}) error
}
