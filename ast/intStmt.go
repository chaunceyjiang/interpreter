package ast

import "Interpreter/token"

// IntegerLiteral 也是一个表达式，因此也要实现Expression
type IntegerLiteral struct {
	Token token.Token // 存储 类型
	Value int64       // 具体的值
}

func (il *IntegerLiteral) expressionNode() {

}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}


