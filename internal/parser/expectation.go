package parser

import (
	"fmt"
	"strings"

	"meganruggiero.com/dicebot/internal/token"
)

type expectation struct {
	expected []string
	received token.Token
}

func (expectation *expectation) Error() string {
	return fmt.Sprintf("line %v column %v: expected %v, got %v",
		expectation.received.Line,
		expectation.received.Column,
		strings.Join(expectation.expected, " or "),
		expectation.received.Quote())
}
