package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"meganruggiero.com/dicebot/internal/ast"
)

func TestSolve(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 8, ast.MultiplyTerm{Left: ast.IntTerm{4}, Right: ast.IntTerm{2}}.Solve())
	assert.Equal(t, 2, ast.DivideTerm{Left: ast.IntTerm{42}, Right: ast.IntTerm{21}}.Solve())
	assert.Equal(t, 42, ast.AddTerm{Left: ast.IntTerm{40}, Right: ast.IntTerm{2}}.Solve())
	assert.Equal(t, -2, ast.SubtractTerm{Left: ast.IntTerm{2}, Right: ast.IntTerm{4}}.Solve())
	// Skip ast.DiceTerm so we don't have to deal with changes to the randomizer.
	assert.Equal(t, 42, ast.IntTerm{Value: 42}.Solve())
}
