package lexer

import (
	"errors"
	"strconv"
	"unicode"
)

type Lexer struct {
	source      string
	currentPos  int
	beginOfLine int
	currentRow  int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
	}
}

func (lexer *Lexer) isEmpty() bool {
	return lexer.currentPos >= len(lexer.source)
}

func (lexer *Lexer) chopCharacter() {
	if !lexer.isEmpty() {
		currentSym := lexer.source[lexer.currentPos]
		lexer.currentPos++
		if currentSym == '\n' {
			lexer.beginOfLine = lexer.currentPos
			lexer.currentRow++
		}
	}
}

func (lexer *Lexer) trimLeft() {
	for !lexer.isEmpty() && unicode.IsSpace(rune(lexer.source[lexer.currentPos])) {
		lexer.chopCharacter()
	}
}

func (lexer *Lexer) dropLine() {
	for !lexer.isEmpty() && lexer.source[lexer.currentPos] != '\n' {
		lexer.chopCharacter()
	}
	if !lexer.isEmpty() {
		lexer.chopCharacter()
	}
}

func (lexer *Lexer) currentPosition() Position {
	return Position{lexer.currentRow, lexer.currentPos - lexer.beginOfLine}
}

func (lexer *Lexer) NextToken() (Token, error) {
	lexer.trimLeft()
	for !lexer.isEmpty() && lexer.source[lexer.currentPos] == '#' {
		lexer.dropLine()
		lexer.trimLeft()
	}
	curPosBegin := lexer.currentPosition()

	if lexer.isEmpty() {
		return Token{EOF, "", curPosBegin}, nil
	}

	symBegin := lexer.source[lexer.currentPos]

	if unicode.IsLetter(rune(lexer.source[lexer.currentPos])) {
		i := lexer.currentPos
		for !lexer.isEmpty() && unicode.IsLetter(rune(lexer.source[lexer.currentPos])) {
			lexer.chopCharacter()
		}
		return Token{NAME, lexer.source[i:lexer.currentPos], curPosBegin}, nil
	}

	literal_tokens := map[uint8]TokenType{
		'(': LPAREN,
		')': RPAREN,
		'{': LBRACE,
		'}': RBRACE,
		',': COMMA,
		';': SEMICOLON,
	}
	literal_token, ok := literal_tokens[lexer.source[lexer.currentPos]]
	if ok {
		lexer.chopCharacter()
		return Token{literal_token, string(symBegin), curPosBegin}, nil
	}

	if symBegin == '"' {
		lexer.chopCharacter()
		start := lexer.currentPos
		for !lexer.isEmpty() && lexer.source[lexer.currentPos] != '"' {
			lexer.chopCharacter()
		}

		if !lexer.isEmpty() {
			text := lexer.source[start:lexer.currentPos]
			lexer.chopCharacter()
			return Token{STRING, text, curPosBegin}, nil
		}
	}
	if unicode.IsDigit(rune(symBegin)) {
		start := lexer.currentPos
		for !lexer.isEmpty() && unicode.IsDigit(rune(lexer.source[lexer.currentPos])) {
			lexer.chopCharacter()
		}
		_, err := strconv.Atoi(lexer.source[start:lexer.currentPos])
		if err != nil {
			return Token{}, err
		}
		return Token{NUMBER, lexer.source[start:lexer.currentPos],
			curPosBegin}, nil
	}

	return Token{}, errors.New("could not recognize token pattern")
}
