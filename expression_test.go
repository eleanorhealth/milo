package milo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert := assert.New(t)

	exprs := []Expression{
		expr{
			field:   "foo",
			operand: OpEqual,
			value:   "bar",
		},
		expr{
			field:   "bar",
			operand: OpNotEqual,
			value:   "baz",
		},
	}

	actual := And(exprs...)

	assert.Equal(exprList{t: ExpressionTypeAnd, exprs: exprs}, actual)
}

func TestOr(t *testing.T) {
	assert := assert.New(t)

	exprs := []Expression{
		expr{
			field:   "foo",
			operand: OpEqual,
			value:   "bar",
		},
		expr{
			field:   "bar",
			operand: OpNotEqual,
			value:   "baz",
		},
	}

	actual := Or(exprs...)

	assert.Equal(exprList{t: ExpressionTypeOr, exprs: exprs}, actual)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	actual := Equal("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpEqual, value: "bar"}, actual)
}

func TestNotEqual(t *testing.T) {
	assert := assert.New(t)

	actual := NotEqual("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpNotEqual, value: "bar"}, actual)
}

func TestGt(t *testing.T) {
	assert := assert.New(t)

	actual := Gt("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpGt, value: "bar"}, actual)
}

func TestLt(t *testing.T) {
	assert := assert.New(t)

	actual := Lt("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpLt, value: "bar"}, actual)
}

func TestGte(t *testing.T) {
	assert := assert.New(t)

	actual := Gte("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpGte, value: "bar"}, actual)
}

func TestLte(t *testing.T) {
	assert := assert.New(t)

	actual := Lte("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpLte, value: "bar"}, actual)
}
