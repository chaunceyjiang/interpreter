package ast

import (
	"Interpreter/token"
	"bytes"
)

type BlockStatement struct {
	Token      token.Token //  { Token
	Statements []Statement //  {} 代码块中的全部语句
}

func (bs *BlockStatement) statementNode() {

}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
