package milo

type Model interface {
	FromEntity(interface{}) error
	ToEntity() (interface{}, error)
}

type Hook interface {
	BeforeSave(store Storer, entity interface{}) error
}
