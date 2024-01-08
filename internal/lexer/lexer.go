package lexer

import (
	"regexp"
	"unicode"
	"unicode/utf8"

	"meganruggiero.com/dicebot/internal/token"
)

const eof = rune(-1)

// Matches "d", case-insensitive, with digits after it.
var regexpD = regexp.MustCompile(`\A[Dd]\d+\z`)

type Lexer struct {
	input           string
	offset          int
	line            int
	column          int
	currentRuneSize int
	currentRune     rune
}

func New(input string) *Lexer {
	lexer := Lexer{
		input:           input,
		offset:          0,
		line:            1,
		column:          0,
		currentRuneSize: 0,
		currentRune:     0,
	}

	lexer.readRune()

	return &lexer
}

//nolint:cyclop,funlen,wsl
func (lexer *Lexer) Read() token.Token {
	// Eat whitespace and do not include it in the token.
	for lexer.currentRune == ',' || unicode.IsSpace(lexer.currentRune) {
		lexer.readRune()
	}

	kind := token.Unrecognized
	start := lexer.offset
	line := lexer.line
	column := lexer.column

	switch lexer.currentRune {
	case utf8.RuneError:
		kind = token.RuneError
		lexer.readRune()
	case eof:
		kind = token.EOF
	case '=':
		kind = token.Equal
		lexer.readRune()
	case '(':
		kind = token.LeftParentheses
		lexer.readRune()
	case ')':
		kind = token.RightParentheses
		lexer.readRune()
	case '^':
		kind = token.Exponentiate
		lexer.readRune()
	case '*':
		kind = token.Multiply
		lexer.readRune()
	case '/':
		kind = token.Divide
		lexer.readRune()
	case '+': // Plus Sign
		kind = token.Add
		lexer.readRune()
	case '-':
		kind = token.Subtract
		lexer.readRune()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		kind = token.Int
		for '0' <= lexer.currentRune && lexer.currentRune <= '9' {
			lexer.readRune()
		}
	default:
		if lexer.currentRune == '_' || unicode.IsLetter(lexer.currentRune) {
			kind = token.Word
			for lexer.currentRune == '_' || unicode.IsLetter(lexer.currentRune) || unicode.IsNumber(lexer.currentRune) {
				lexer.readRune()
			}
		} else {
			// Keep kind set to token.Unrecognized.
			lexer.readRune()
		}
	}

	str := lexer.input[start:lexer.offset]

	if regexpD.MatchString(str) {
		kind = token.D
	}

	return token.New(line, column, kind, str)
}

func (lexer *Lexer) readRune() {
	const (
		FileSeparator           = rune(0x1C)
		GroupSeparator          = rune(0x1D)
		InformationSeparatorTwo = rune(0x1E)
		NextLine                = rune(0x85)
		LineSeparator           = rune(0x2028)
		ParagraphSeparator      = rune(0x2029)
	)

	lexer.offset += lexer.currentRuneSize
	nextRune, nextRuneSize := utf8.DecodeRuneInString(lexer.input[lexer.offset:])

	if nextRuneSize == 0 {
		nextRune = eof // Easier for switch statements.
	} else {
		switch nextRune {
		case '\n', '\r', FileSeparator, GroupSeparator, InformationSeparatorTwo, NextLine, LineSeparator, ParagraphSeparator:
			lexer.line++
			lexer.column = 0
		default:
			lexer.column++
		}
	}

	lexer.currentRuneSize = nextRuneSize
	lexer.currentRune = nextRune
}
