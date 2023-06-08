package lexer

type Token struct {
	Type     TokenType
	Text     string
	Position Position
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	NEWLINE
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	SEMICOLON
	NUMBER
	STRING

	NAME

	COMMA
	RETURN
)

var tokenNames = map[TokenType]string{
	ILLEGAL: "<illegal>",
	EOF:     "EOF",
	NEWLINE: "<newline>",

	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	SEMICOLON: ";",
	COMMA:     ",",

	RETURN: "return",

	NUMBER: "number",
	STRING: "string",
	NAME:   "name",
}

// String returns the string name of this token.
func (t TokenType) String() string {
	return tokenNames[t]
}
