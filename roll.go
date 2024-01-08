package main

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"meganruggiero.com/dicebot/internal/parser"
)

func roll(input string) string {
	var output strings.Builder

	fmt.Fprintf(&output, "**Rolling**: %v", discordEscapeMarkdown(input))

	formula, err := parser.Parse(input)
	if err != nil {
		fmt.Fprintf(&output, "\n**Syntax Error**: %v", discordEscapeMarkdown(err.Error()))

		return output.String()
	}

	for index, equation := range formula.Equations {
		name := equation.Name
		if name == "" {
			name = humanize.Ordinal(index + 1)
		}

		fmt.Fprintf(&output, "\n**%v**: %v", discordEscapeMarkdown(name), equation.Term.Solve())
	}

	return output.String()
}
