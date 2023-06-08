package parser

import (
	"C_compiler/lexer"
	"errors"
)

func fetchExpectedToken(lex *lexer.Lexer, tokenTypes ...lexer.TokenType) (lexer.Token, error) {
	token, err := lex.NextToken()
	if err != nil {
		return lexer.Token{}, err
	}

	isGood := false
	for _, tokenType := range tokenTypes {
		if token.Type == tokenType {
			isGood = true
			//return lexer.Token{}, errors.New("expected token type: " + tokenType.String() + ", got: " + token.Type.String())
		}
	}

	if !isGood {
		return lexer.Token{}, errors.New("expected token not found")
	}

	return token, nil
}

func parseType(lex *lexer.Lexer) (ReturnType, error) {
	returnType, err := fetchExpectedToken(lex, lexer.NAME)
	if err != nil {
		return VOID, err
	}

	if returnType.Text != "int" {
		return VOID, errors.New("unimplemented error type")
	}

	return INT, err
}

func parseArgList(lex *lexer.Lexer) (args []string, err error) {
	_, err = fetchExpectedToken(lex, lexer.LPAREN)
	if err != nil {
		return args, err
	}
	expr, err := fetchExpectedToken(lex, lexer.STRING, lexer.NUMBER, lexer.RPAREN)
	if err != nil {
		return args, err
	}
	if expr.Type == lexer.RPAREN {
		return args, nil
	}
	args = append(args, expr.Text)
	for {
		expr, err = fetchExpectedToken(lex, lexer.RPAREN, lexer.COMMA)
		if err != nil {
			return args, err
		}
		if expr.Type == lexer.RPAREN {
			break
		}
		expr, err = fetchExpectedToken(lex, lexer.STRING, lexer.NUMBER)
		if err != nil {
			return args, err
		}
		args = append(args, expr.Text)
	}
	return args, nil
}

func parseBlock(lex *lexer.Lexer) (block []any, err error) {
	_, err = fetchExpectedToken(lex, lexer.LBRACE)
	if err != nil {
		return block, err
	}

	for {
		name, err := fetchExpectedToken(lex, lexer.NAME, lexer.RBRACE)
		if err != nil {
			return block, err
		}
		if name.Type == lexer.RBRACE {
			break
		}
		if name.Text == "return" {
			expr, err := fetchExpectedToken(lex, lexer.NUMBER, lexer.STRING)
			if err != nil {
				return block, err
			}
			block = append(block, ReturnStatement{expr.Text})
		} else {
			args, err := parseArgList(lex)
			if err != nil {
				return block, err
			}
			block = append(block, FuncCallStatement{name.Text, args})

		}
		_, err = fetchExpectedToken(lex, lexer.SEMICOLON)
		if err != nil {
			return block, err
		}
	}

	return block, nil
}

func ParseFunction(lex *lexer.Lexer) (Function, error) {
	_, err := parseType(lex)
	if err != nil {
		return Function{}, err
	}

	name, err := fetchExpectedToken(lex, lexer.NAME)
	if err != nil {
		return Function{}, err
	}
	_, err = fetchExpectedToken(lex, lexer.LPAREN)
	if err != nil {
		return Function{}, err
	}
	_, err = fetchExpectedToken(lex, lexer.RPAREN)
	if err != nil {
		return Function{}, err
	}

	body, err := parseBlock(lex)
	if err != nil {
		return Function{}, err
	}
	return Function{name.Text, body}, nil
}
