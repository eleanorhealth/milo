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

type ExpressionType string

const (
	ExpressionTypeAnd ExpressionType = "AND"
	ExpressionTypeOr  ExpressionType = "OR"
)

type Expression interface {
	Field() interface{}
	Operand() Op
	Value() interface{}
	Type() ExpressionType
	Expressions() []Expression
}

type expr struct {
	field   interface{}
	operand Op
	value   interface{}
	t       ExpressionType
	exprs   []Expression
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

func (e expr) Type() ExpressionType {
	return e.t
}

func (e expr) Expressions() []Expression {
	return e.exprs
}

func And(exprs ...Expression) Expression {
	return expr{exprs: exprs}
}

func Or(exprs ...Expression) Expression {
	for i, e := range exprs {
		if expr, ok := e.(expr); ok {
			expr.t = ExpressionTypeOr
			exprs[i] = expr
		}
	}

	return expr{exprs: exprs}
}

func Equal(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpEqual,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}

func NotEqual(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpNotEqual,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}

func Gt(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpGt,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}

func Lt(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpLt,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}

func Gte(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpGte,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}

func Lte(field interface{}, value interface{}) Expression {
	return expr{
		field:   field,
		operand: OpLte,
		value:   value,
		t:       ExpressionTypeAnd,
	}
}
