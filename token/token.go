package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string // 存储字面值，代表代码里面的值 因此都是字符串
}

const (
	//  two special types
	ILLEGAL = "ILLEGAL" // 非法的
	EOF     = "EOF"

	// 标示符 和 字面量
	// Identifiers + literals
	IDENT = "IDENT" // add ,foo,bar ,x, y ....
	INT   = "INT"   // 12342

	//操作符Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	// 双字符操作符

	EQ          = "=="
	PLUS_EQ     = "+="
	MINUS_EQ    = "-="
	NOT_EQ      = "!="
	ASTERISK_EQ = "*="
	SLASH_EQ    = "/="
	LE          = "<="
	GE          = ">="

	// 分隔符 Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"



)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"false":  FALSE,
	"true":   TRUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		// 如果是关键字，则返回 关键字
		return tok
	}
	// 不是关键字 则返回变量 IDENT
	return IDENT
}

var twoCharToken = map[string]TokenType{
	"==": EQ,
	"+=": PLUS_EQ,
	"-=": MINUS_EQ,
	"!=": NOT_EQ,
	"*=": ASTERISK_EQ,
	"/=": SLASH_EQ,
	"<=": LE,
	">=": GE,
}

func LookupTwoCharToken(ch byte) TokenType {
	literal := string(ch) + "="
	return twoCharToken[literal]
}
