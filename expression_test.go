package milo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert := assert.New(t)

	exprs := []Expression{
		{
			field: "foo",
			op:    OpEqual,
			value: "bar",
			t:     expressionTypeAnd,
		},
		{
			field: "bar",
			op:    OpNotEqual,
			value: "baz",
			t:     expressionTypeAnd,
		},
	}

	actual := And(exprs...)

	assert.Equal(Expression{exprs: exprs}, actual)
}

func TestOr(t *testing.T) {
	assert := assert.New(t)

	exprs := []Expression{
		{
			field: "foo",
			op:    OpEqual,
			value: "bar",
			t:     expressionTypeOr,
		},
		{
			field: "bar",
			op:    OpNotEqual,
			value: "baz",
			t:     expressionTypeOr,
		},
	}

	actual := Or(exprs...)

	assert.Equal(Expression{exprs: exprs}, actual)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	actual := Equal("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpEqual, value: "bar", t: expressionTypeAnd}, actual)
}

func TestNotEqual(t *testing.T) {
	assert := assert.New(t)

	actual := NotEqual("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpNotEqual, value: "bar", t: expressionTypeAnd}, actual)
}

func TestGt(t *testing.T) {
	assert := assert.New(t)

	actual := Gt("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpGt, value: "bar", t: expressionTypeAnd}, actual)
}

func TestLt(t *testing.T) {
	assert := assert.New(t)

	actual := Lt("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpLt, value: "bar", t: expressionTypeAnd}, actual)
}

func TestGte(t *testing.T) {
	assert := assert.New(t)

	actual := Gte("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpGte, value: "bar", t: expressionTypeAnd}, actual)
}

func TestLte(t *testing.T) {
	assert := assert.New(t)

	actual := Lte("foo", "bar")

	assert.Equal(Expression{field: "foo", op: OpLte, value: "bar", t: expressionTypeAnd}, actual)
}
