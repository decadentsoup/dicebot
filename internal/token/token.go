package token

import "fmt"

type Kind int

const (
	Unrecognized Kind = iota
	RuneError
	EOF
	Equal
	LeftParentheses
	RightParentheses
	Exponentiate
	Multiply
	Divide
	Add
	Subtract
	D
	Int
	ID
)

type Token struct {
	Line   int
	Column int
	Kind   Kind
	String string
}

func New(line, column int, kind Kind, str string) Token {
	return Token{Line: line, Column: column, Kind: kind, String: str}
}

func (token Token) Int() int {
	value := 0

	// We can iterate on byte instead of rune because integers are currently only ASCII.
	for _, currentByte := range []byte(token.String) {
		if '0' <= currentByte && currentByte <= '9' {
			value = value*10 + (int(currentByte) - '0') //nolint:gomnd
		}
	}

	return value
}

func (token Token) Quote() string {
	switch token.Kind { //nolint:exhaustive
	case RuneError:
		return "invalid utf-8 sequence"
	case EOF:
		return "end of input"
	default:
		return fmt.Sprintf("%q", token.String)
	}
}
