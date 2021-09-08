package milo

type Op string

const (
	OpIsNull    Op = "IS NULL"
	OpIsNotNull Op = "IS NOT NULL"
	OpEqual     Op = "="
	OpNotEqual  Op = "!="
	OpGt        Op = ">"
	OpLt        Op = "<"
	OpGte       Op = ">="
	OpLte       Op = "<="
)

type expressionType int

const (
	expressionTypeAnd expressionType = iota
	expressionTypeOr
)

type Expression struct {
	column interface{}
	op     Op
	value  interface{}
	t      expressionType
	exprs  []Expression
}

func (e Expression) Column() interface{} {
	return e.column
}

func (e Expression) Op() interface{} {
	return e.op
}

func (e Expression) Value() interface{} {
	return e.value
}

func And(exprs ...Expression) Expression {
	for i, expr := range exprs {
		expr.t = expressionTypeAnd
		exprs[i] = expr
	}

	return Expression{exprs: exprs}
}

func Or(exprs ...Expression) Expression {
	for i, expr := range exprs {
		expr.t = expressionTypeOr
		exprs[i] = expr
	}

	return Expression{exprs: exprs}
}

func Equal(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpEqual,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func NotEqual(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpNotEqual,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func Gt(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpGt,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func Lt(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpLt,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func Gte(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpGte,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func Lte(column interface{}, value interface{}) Expression {
	return Expression{
		column: column,
		op:     OpLte,
		value:  value,
		t:      expressionTypeAnd,
	}
}

func IsNull(column interface{}) Expression {
	return Expression{
		column: column,
		op:     OpIsNull,
		value:  nil,
		t:      expressionTypeAnd,
	}
}

func IsNotNull(column interface{}) Expression {
	return Expression{
		column: column,
		op:     OpIsNotNull,
		value:  nil,
		t:      expressionTypeAnd,
	}
}
