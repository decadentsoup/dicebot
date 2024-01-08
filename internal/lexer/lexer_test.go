package lexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"meganruggiero.com/dicebot/internal/lexer"
	"meganruggiero.com/dicebot/internal/token"
)

func TestLexer(t *testing.T) {
	t.Parallel()

	lexer := lexer.New("2d4 + d20 - 4 D8?,\nkeyword = Dice_1234")

	expectations := []token.Token{
		token.New(1, 1, token.Int, "2"),
		token.New(1, 2, token.D, "d4"),
		token.New(1, 5, token.Add, "+"),
		token.New(1, 7, token.D, "d20"),
		token.New(1, 11, token.Subtract, "-"),
		token.New(1, 13, token.Int, "4"),
		token.New(1, 15, token.D, "D8"),
		token.New(1, 17, token.Unrecognized, "?"),
		token.New(2, 1, token.Word, "keyword"),
		token.New(2, 9, token.Equal, "="),
		token.New(2, 11, token.Word, "Dice_1234"),
		token.New(2, 19, token.EOF, ""),
	}

	for index, expectation := range expectations {
		assert.Equal(t, expectation, lexer.Read(), "token %v should match expectation", index)
	}
}
