package milo

type Model interface {
	FromEntity(interface{}) error
	ToEntity() (interface{}, error)
}
