package ast

import "Interpreter/token"

type Identifier struct {
	Token token.Token
	Value string
}

// 实现了 expressionNode() 因此 也实现了 Expression
func (i *Identifier) expressionNode() {

}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string { return i.Value }
