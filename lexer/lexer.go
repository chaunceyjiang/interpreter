package lexer

import "Interpreter/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)  当前字符的位置
	readPosition int  // current reading position in input (after current char) 当前字符的位置的后一个位置
	ch           byte // current char under examination   当前正在处理的字符

}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 初始化 position
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	// 处理当前字符
	switch l.ch {
	case '=':
		tok = l.makeTwoCharToken(token.ASSIGN, l.ch)
	case '!':
		tok = l.makeTwoCharToken(token.BANG, l.ch)
	case '+':
		tok = l.makeTwoCharToken(token.PLUS, l.ch)
	case '-':
		tok = l.makeTwoCharToken(token.MINUS, l.ch)
	case '*':
		tok = l.makeTwoCharToken(token.ASTERISK, l.ch)
	case '/':
		tok = l.makeTwoCharToken(token.SLASH, l.ch)
	case '<':
		tok = l.makeTwoCharToken(token.LT, l.ch)
	case '>':
		tok = l.makeTwoCharToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			// 如果当前是字母，那么获取 该字母后面的连续字母 ()
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			// 如果当前位置是 数字，那么将后面的连续数值取处理
			Literal,isFloat:= l.readNumber()
			tok.Literal =Literal
			if isFloat {
				tok.Type = token.FlOAT
			}else {
				tok.Type = token.INT
			}


			return tok

		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar() // 继续处理下一个字符
	return tok
}
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		// 一直读取到 "  结束
		l.readChar()
		if l.ch == '"' {
			break
		}
	}
	return l.input[position:l.position]
}

// readChar 获取当前字符
func (l *Lexer) readChar() {
	// readPosition always points to the "next" character in the input.

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar 提前查看下一个要解析的字符
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) makeTwoCharToken(orgTok token.TokenType, ch byte) token.Token {
	var tok token.TokenType
	var literal string
	var peekChar = l.peekChar()
	if peekChar == '=' {
		// 这里处理了 下一个字符，所以这个应该让指针前进一个字符
		tok = token.LookupTwoCharToken(ch)
		literal = string(ch) + string(peekChar)
		l.readChar()
	} else {
		tok = orgTok
		literal = string(ch)
	}
	return token.Token{Type: tok, Literal: literal}
}
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position] // 获取连续的字母

}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	var isFloat = false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar()
		isFloat = true
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position],isFloat
}

// isLetter 判断是不是字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit 判断是不是数字
func isDigit(ch byte) bool {

	return '0' <= ch && ch <= '9'
}
