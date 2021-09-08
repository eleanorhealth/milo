package entityid

import "github.com/google/uuid"

var DefaultGenerator = NewUUIDGenerator()

type ID string

func (i ID) String() string {
	return string(i)
}

type Generator interface {
	Generate() ID
}

type UUIDGenerator struct {
}

var _ Generator = (*UUIDGenerator)(nil)

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (u *UUIDGenerator) Generate() ID {
	return ID(uuid.New().String())
}
