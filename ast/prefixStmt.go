package ast

import (
	"Interpreter/token"
	"bytes"
)

// PrefixExpression 前缀表达式
type PrefixExpression struct {
	Token    token.Token
	Operator string     // 前缀操作符
	Right    Expression // 右边的表达式
}

func (pe *PrefixExpression) expressionNode() {

}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()

}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
