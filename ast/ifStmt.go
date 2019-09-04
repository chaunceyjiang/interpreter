package ast

import (
	"Interpreter/token"
	"bytes"
)

type IfExpression struct {
	Token token.Token  // if
	Condition Expression // 条件
	Consequence *BlockStatement  // if 为true 的结果
	Alternative *BlockStatement // else 可选
}


func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString(" if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil{
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}