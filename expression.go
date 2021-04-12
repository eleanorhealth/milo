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

type ExpressionListType string

const (
	ExpressionTypeAnd ExpressionListType = "AND"
	ExpressionTypeOr  ExpressionListType = "OR"
)

type Expression interface {
	Field() interface{}
	Operand() Op
	Value() interface{}
}

type ExpressionList interface {
	Expression
	Type() ExpressionListType
	Expressions() []Expression
}

type expr struct {
	field   interface{}
	operand Op
	value   interface{}
}

func (e expr) Field() interface{} {
	return e.field
}

func (e expr) Operand() Op {
	return e.operand
}

func (e expr) Value() interface{} {
	return e.value
}

type exprList struct {
	expr
	t     ExpressionListType
	exprs []Expression
}

func (e exprList) Type() ExpressionListType {
	return e.t
}

func (e exprList) Expressions() []Expression {
	return e.exprs
}

func And(exprs ...Expression) ExpressionList {
	return exprList{t: ExpressionTypeAnd, exprs: exprs}
}

func Or(exprs ...Expression) ExpressionList {
	return exprList{t: ExpressionTypeOr, exprs: exprs}
}

func Equal(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpEqual,
		value:   value,
	}
}

func NotEqual(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpNotEqual,
		value:   value,
	}
}

func Gt(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpGt,
		value:   value,
	}
}

func Lt(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpLt,
		value:   value,
	}
}

func Gte(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpGte,
		value:   value,
	}
}

func Lte(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpLte,
		value:   value,
	}
}
