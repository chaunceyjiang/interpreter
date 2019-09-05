package parser

import (
	"Interpreter/ast"
	"Interpreter/lexer"
	"Interpreter/token"
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	errors         []error
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn // 根绝不同的 Token 来解析不同的前缀表达式
	infixParseFns  map[token.TokenType]infixParseFn  // 根绝不同的 Token 来解析不同的中缀表达式
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []error{}}

	p.nextToken() // 因为要让  curToken 和 peekToken 都初始化，所以这里调用了两次
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // 初始化 一元 解析函数map
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // 注册 变量表达式 的解析器
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         // 注册 int表达式  的解析器
	p.registerPrefix(token.BANG, p.parsePrefixExpression)      // 注册   ! 的 解析器      ！true
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     // 注册  - 的 解析器   -2
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.FlOAT, p.parseFloatLiteral)
	p.infixParseFns = make(map[token.TokenType]infixParseFn) // 初始化 二元 解析函数map
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression) // ( 这个 在两个地方注册，这里代表 函数调用
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// 解析let 语法
		return p.parseLetStatement()
	case token.RETURN:
		// 解析return 语句
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	// 如果当前是 ; 标示语句结束 ，因此 调用 Token 指针 向前走一步，准备下次调用
	//for !p.curTokenIs(token.SEMICOLON) {
	//	p.nextToken()
	//}

	stmt.Value = p.parseExpression(LOWEST) // 解析 let 关键字后面的表达式
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST) // 解析return 后面的表达式
	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression 将每一个Token 解析成不同的表达式，核心处理就是 ParseFns map 中注册是 Token 解析方法
func (p *Parser) parseExpression(precedence int) ast.Expression {

	prefixFn := p.prefixParseFns[p.curToken.Type] // 根据不同的Token ，调用不同的处理方法 ，基本都是 INT IDENT 类型的处理函数

	if prefixFn == nil {
		// 增加错误信息，方便调试
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefixFn() // 因为这里 都是INT IDENT 类型的处理函数 ，所以传过来的 优先级，是前一个操作符的优先级precedence

	// p.peekPrecedence() 这里比较是 当前的INT IDENT，的Token的前面一个操作符的优先级 与下一个操作符优先级的比较
	// 进入for 循环，代表则下一个操作符的优先级更高，因此，应该先处理优先级更高的操作符
	// for 循环结束，代表着 下一个优先级比较低，代表着，高优先级已经处理结束，退回到 从左往右计算的 正常处理中， 又因为curToken 指针一直在忘前走，所以不会重复的解析表达式

	// 括号 表达式 正是应用了该原理， 先用 ( 降低 前面操作符的优先级，然后进入 for 循环处理 意味着 先处理了括号里面的表达式
	// 然后 再遇到后，用 ) 降低 括号里面表达式的优先级 ， 意味着 括号中的表达式 解析结束，递归 结束。 退回 从左往右计算的 正常处理中，
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { // 如果当前的优先级 小于 未来的优先级 则继续查找
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil { // 结束 表达式 查找
			return leftExp
		}

		p.nextToken() // 循环迭代

		leftExp = infixFn(leftExp)

	}

	return leftExp

}

// parseBoolean  实现了 prefixParseFns func() ast.Expression   解析布尔类型
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)} // 这里的Values的值 巧用了 curTokenIs 函数，判断当前的 TokenType 是 TRUE,还是FALSE
}

// parseIntegerLiteral  实现了 prefixParseFns func() ast.Expression   解析数值类型
func (p *Parser) parseIntegerLiteral() ast.Expression {
	ilt := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		// 不能解析该值，然后将错误保存传递出去
		p.errors = append(p.errors, errors.New(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)))
	}
	ilt.Value = value
	return ilt

}

func (p *Parser) parseFloatLiteral() ast.Expression {
	ilt := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal,64)
	if err != nil {
		// 不能解析该值，然后将错误保存传递出去
		p.errors = append(p.errors, errors.New(fmt.Sprintf("could not parse %q as float", p.curToken.Literal)))
	}
	ilt.Value = value
	return ilt

}

// parseGroupedExpression 小括号的表达式
func (p *Parser) parseGroupedExpression() ast.Expression {

	// 跳过 括号
	p.nextToken()

	// 降低当前的Token的优先级，然后调用 parseExpression 进行递归处理， 进入parseExpression 中的for 中 进行括号中的表达式解析。
	exp := p.parseExpression(LOWEST)

	// 递归结束 (parseExpression 中的for 循环 结束)  这里进行检查，如果让递归结束的不是 右括号 ) 则代表 括号表达式有问题，因此返回 nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// 让递归结束的是 ) 括号 同时 调用 nextToken ，跳过 ) 括号  `1 + (1 + 2 ) + 2` ,因为 ) 右边是另个 操作符
	// exp 整个括号中的表达式
	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}
func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	args = make([]ast.Expression, 0)
	if p.peekTokenIs(token.RPAREN) { // 无参数调用
		p.nextToken()
		return args
	}
	p.nextToken()

	args = append(args, p.parseExpression(LOWEST)) // 解析第一个参数

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.peekTokenIs(token.RPAREN) {
			break
		}
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIdentifier  实现了 prefixParseFns func() ast.Expression
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parsePrefixExpression  实现了 prefixParseFns func() ast.Expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	// 将 Token(-)  Token(212) 组合成 Token(-212)
	prefixEp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	p.nextToken() // 因为这个是prefix表达式，所以要获取下一个 Token ，

	prefixEp.Right = p.parseExpression(PREFIX) // 将一下个 Token 解析 成表达式  这里是一个递归

	return prefixEp
}

// parseInfixExpression 实现了 infixParseFn func(ast.Expression) ast.Expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left} // 左边的表达式

	precedence := p.curPrecedence() // 获取当前Token 的优先级 ，以便后续处理

	p.nextToken() // 取下一个 Token

	expression.Right = p.parseExpression(precedence) // 右边的表达式

	return expression

}

// parseIfExpression
func (p *Parser) parseIfExpression() ast.Expression {
	// if ( x < y ) { x } else { y }
	// 能进这个函数，表示当前的 Token 是 If 关键字
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) { // 判断下一个 Token是不是 （ ，如果是 则前进一个 Token
		return nil
	}

	p.nextToken() // 下一个 是 ( ，因此指针前进一个Token
	// 解析 if 括号()中的 的条件表达式
	expression.Condition = p.parseExpression(LOWEST) // 解析括号中的表达式

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// ) {  这两个符号是连在一起的
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// 解析 {} 中的代码块
	expression.Consequence = p.parseBlockStatement()

	// 解析 else 中的表达式
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// fn(x, y) { x + y; }
	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn

}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var params []*ast.Identifier
	params = make([]*ast.Identifier, 0)
	if p.peekTokenIs(token.RPAREN) { // 无参数列表
		p.nextToken()
		return params
	}
	// fn(x, y) { x + y; }
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, ident)
	for p.peekTokenIs(token.COMMA) { // 参数列表，都是以为 `,` 为分分隔符

		p.nextToken()

		if p.peekTokenIs(token.RPAREN) { // 参数列表 (x,y,)
			break
		}
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, ident)
	}
	if !p.expectPeek(token.RPAREN) { // 函数参数列表解析结束 ，同时 跳过 )
		return nil
	}
	return params
}

// parseBlockStatement 这里返回的是  {} 代码块中的全部语句
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // 跳过 { 左大括号

	for !p.curTokenIs(token.RBRACE) { // 直到 } 右大括号结束
		stmt := p.parseStatement() // 进入大括号后，解析大括号中的所有语句
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// curTokenIs 检查当前Token 的类型
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t

}

// peekTokenIs 检查下一个Token 的类型
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t

}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else { // 如果下一个不是期望的 TokenType 则说明语法有问题，则加入错误中列表中
		p.peekError(t)
		return false
	}

}

// noPrefixParseFnError  当 该类型 没有解析函数时，增加错误信息
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors, errors.New(fmt.Sprintf("no prefix parse function for %s found", t)))
}
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	err := errors.New(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	p.errors = append(p.errors, err)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
