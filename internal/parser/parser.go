package parser

import (
	"meganruggiero.com/dicebot/internal/ast"
	"meganruggiero.com/dicebot/internal/lexer"
	"meganruggiero.com/dicebot/internal/token"
)

func Parse(input string) (*ast.Formula, error) {
	parser := parser{
		lexer:        lexer.New(input),
		currentToken: token.New(0, 0, token.Unrecognized, ""),
	}

	parser.readToken()

	formula, err := parser.parseFormula()
	if err != nil {
		return nil, err
	}

	return formula, nil
}

type parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
}

func (parser *parser) expected(expected ...string) *expectation {
	return &expectation{expected: expected, received: parser.currentToken}
}

func (parser *parser) readToken() {
	nextToken := parser.lexer.Read()
	parser.currentToken = nextToken
}

func (parser *parser) parseFormula() (*ast.Formula, *expectation) {
	equations := make([]ast.Equation, 0)

	for parser.currentToken.Kind != token.EOF {
		equation, err := parser.parseEquation()
		if err != nil {
			return nil, err
		}

		equations = append(equations, *equation)
	}

	return &ast.Formula{Equations: equations}, nil
}

func (parser *parser) parseEquation() (*ast.Equation, *expectation) {
	name, err := parser.parseOptionalEquationName()
	if err != nil {
		return nil, err
	}

	term, err := parser.parseETerm()
	if err != nil {
		return nil, err
	}

	return &ast.Equation{Name: name, Term: term}, nil
}

func (parser *parser) parseOptionalEquationName() (string, *expectation) {
	if parser.currentToken.Kind != token.ID {
		return "", nil
	}

	name := parser.currentToken.String

	parser.readToken()

	if parser.currentToken.Kind != token.Equal {
		return "", parser.expected(`"="`)
	}

	parser.readToken()

	return name, nil
}

func (parser *parser) parseETerm() (ast.Term, *expectation) {
	left, err := parser.parseMDTerm()
	if err != nil {
		return nil, err
	}

	if parser.currentToken.Kind == token.Exponentiate {
		parser.readToken()

		right, err := parser.parseETerm()
		if err != nil {
			return nil, err
		}

		return ast.ExponentiateTerm{Left: left, Right: right}, nil
	}

	return left, nil
}

func (parser *parser) parseMDTerm() (ast.Term, *expectation) {
	left, err := parser.parseASTerm()
	if err != nil {
		return nil, err
	}

	for {
		switch parser.currentToken.Kind { //nolint:exhaustive
		case token.Multiply:
			parser.readToken()

			right, err := parser.parseASTerm()
			if err != nil {
				return nil, err
			}

			left = ast.MultiplyTerm{Left: left, Right: right}
		case token.Divide:
			parser.readToken()

			right, err := parser.parseASTerm()
			if err != nil {
				return nil, err
			}

			left = ast.DivideTerm{Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

func (parser *parser) parseASTerm() (ast.Term, *expectation) {
	left, err := parser.parseBottomTerm()
	if err != nil {
		return nil, err
	}

	for {
		switch parser.currentToken.Kind { //nolint:exhaustive
		case token.Add:
			parser.readToken()

			right, err := parser.parseBottomTerm()
			if err != nil {
				return nil, err
			}

			left = ast.AddTerm{Left: left, Right: right}
		case token.Subtract:
			parser.readToken()

			right, err := parser.parseBottomTerm()
			if err != nil {
				return nil, err
			}

			left = ast.SubtractTerm{Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

//nolint:cyclop,funlen
func (parser *parser) parseBottomTerm() (ast.Term, *expectation) {
	switch parser.currentToken.Kind { //nolint:exhaustive
	case token.D:
		parser.readToken()

		if parser.currentToken.Kind != token.Int {
			return nil, parser.expected("integer")
		}

		faces := parser.currentToken.Int()
		parser.readToken()

		return ast.DiceTerm{Count: 1, Faces: faces}, nil
	case token.Int:
		intOrCount := parser.currentToken.Int()

		parser.readToken()

		if parser.currentToken.Kind == token.D {
			parser.readToken()

			if parser.currentToken.Kind != token.Int {
				return nil, parser.expected("integer")
			}

			faces := parser.currentToken.Int()
			parser.readToken()

			return ast.DiceTerm{Count: intOrCount, Faces: faces}, nil
		}

		return ast.IntTerm{Value: intOrCount}, nil
	case token.Add:
		parser.readToken()

		if parser.currentToken.Kind != token.Int {
			return nil, parser.expected("integer")
		}

		parser.readToken()

		return ast.IntTerm{Value: +parser.currentToken.Int()}, nil
	case token.Subtract:
		parser.readToken()

		if parser.currentToken.Kind != token.Int {
			return nil, parser.expected("integer")
		}

		parser.readToken()

		return ast.IntTerm{Value: -parser.currentToken.Int()}, nil
	case token.LeftParentheses:
		parser.readToken()

		term, err := parser.parseETerm()
		if err != nil {
			return nil, err
		}

		if parser.currentToken.Kind != token.RightParentheses {
			return nil, parser.expected(`")"`)
		}

		parser.readToken()

		return term, nil
	default:
		return nil, parser.expected("integer", "dice term", `"("`)
	}
}
