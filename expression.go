package milo

type Op string

const (
	OpEqual    Op = "="
	OpNotEqual Op = "!="
	OpGt       Op = ">"
	OpLt       Op = "<"
	OpGte      Op = ">="
	OpLte      Op = "<="
)

type expressionType int

const (
	expressionTypeAnd expressionType = iota
	expressionTypeOr
)

type Expression struct {
	field interface{}
	op    Op
	value interface{}
	t     expressionType
	exprs []Expression
}

func (e Expression) Field() interface{} {
	return e.field
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

func Equal(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpEqual,
		value: value,
		t:     expressionTypeAnd,
	}
}

func NotEqual(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpNotEqual,
		value: value,
		t:     expressionTypeAnd,
	}
}

func Gt(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpGt,
		value: value,
		t:     expressionTypeAnd,
	}
}

func Lt(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpLt,
		value: value,
		t:     expressionTypeAnd,
	}
}

func Gte(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpGte,
		value: value,
		t:     expressionTypeAnd,
	}
}

func Lte(field interface{}, value interface{}) Expression {
	return Expression{
		field: field,
		op:    OpLte,
		value: value,
		t:     expressionTypeAnd,
	}
}
