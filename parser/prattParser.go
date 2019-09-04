package parser

import (
	"Interpreter/ast"
	"Interpreter/token"
)

const (
	// 优先级
	_           int = iota
	LOWEST          //                       1
	EQUALS          // == !=                 2
	LESSGREATER     // > or < <= >=          3
	SUM             // + -                   4
	PRODUCT         // * /                   5
	PREFIX          // -x or !x              6
	CALL            // f(x)                  7

)

type (
	// 都返回一个表达式
	prefixParseFn func() ast.Expression                          // 一元操作符
	infixParseFn  func(ast.Expression) ast.Expression // 二元操作符

)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN: CALL,     // 函数调用 ，add()
}
