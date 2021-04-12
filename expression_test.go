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
			t:       ExpressionTypeAnd,
		},
		expr{
			field:   "bar",
			operand: OpNotEqual,
			value:   "baz",
			t:       ExpressionTypeAnd,
		},
	}

	actual := And(exprs...)

	assert.Equal(expr{exprs: exprs}, actual)
}

func TestOr(t *testing.T) {
	assert := assert.New(t)

	exprs := []Expression{
		expr{
			field:   "foo",
			operand: OpEqual,
			value:   "bar",
			t:       ExpressionTypeOr,
		},
		expr{
			field:   "bar",
			operand: OpNotEqual,
			value:   "baz",
			t:       ExpressionTypeOr,
		},
	}

	actual := Or(exprs...)

	assert.Equal(expr{exprs: exprs}, actual)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	actual := Equal("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpEqual, value: "bar", t: ExpressionTypeAnd}, actual)
}

func TestNotEqual(t *testing.T) {
	assert := assert.New(t)

	actual := NotEqual("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpNotEqual, value: "bar", t: ExpressionTypeAnd}, actual)
}

func TestGt(t *testing.T) {
	assert := assert.New(t)

	actual := Gt("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpGt, value: "bar", t: ExpressionTypeAnd}, actual)
}

func TestLt(t *testing.T) {
	assert := assert.New(t)

	actual := Lt("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpLt, value: "bar", t: ExpressionTypeAnd}, actual)
}

func TestGte(t *testing.T) {
	assert := assert.New(t)

	actual := Gte("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpGte, value: "bar", t: ExpressionTypeAnd}, actual)
}

func TestLte(t *testing.T) {
	assert := assert.New(t)

	actual := Lte("foo", "bar")

	assert.Equal(expr{field: "foo", operand: OpLte, value: "bar", t: ExpressionTypeAnd}, actual)
}
