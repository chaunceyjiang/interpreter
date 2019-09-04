package ast

import (
	"Interpreter/token"
	"bytes"
	"strings"
)

type CallExpression struct {
	Token     token.Token  // ( Token
	Function  Expression   // 变量或函数 Identifier or FunctionLiteral
	Arguments []Expression //
}

func (ce *CallExpression) expressionNode()      {}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
