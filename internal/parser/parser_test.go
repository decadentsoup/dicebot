package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"meganruggiero.com/dicebot/internal/ast"
	"meganruggiero.com/dicebot/internal/parser"
)

//nolint:funlen
func TestLexer(t *testing.T) {
	t.Parallel()

	formula, err := parser.Parse(`
		50, 5d8, d20-8
		pemdas = 5 + 2 * 8 / 4 ^ 7 ^ (8 - 1 D 100)
		unary operations
			aka signs = 1 + +1 - -1
	`)
	assert.NoError(t, err)
	assert.Equal(t, &ast.Formula{Equations: []ast.Equation{
		{Name: "", Term: ast.IntTerm{Value: 50}},
		{Name: "", Term: ast.DiceTerm{Count: 5, Faces: 8}},
		{
			Name: "", Term: ast.SubtractTerm{
				Left:  ast.DiceTerm{Count: 1, Faces: 20},
				Right: ast.IntTerm{Value: 8},
			},
		},
		{
			Name: "pemdas", Term: ast.ExponentiateTerm{
				Left: ast.DivideTerm{
					Left: ast.MultiplyTerm{
						Left: ast.AddTerm{
							Left:  ast.IntTerm{Value: 5},
							Right: ast.IntTerm{Value: 2},
						},
						Right: ast.IntTerm{Value: 8},
					},
					Right: ast.IntTerm{Value: 4},
				},
				Right: ast.ExponentiateTerm{
					Left: ast.IntTerm{Value: 7},
					Right: ast.SubtractTerm{
						Left:  ast.IntTerm{Value: 8},
						Right: ast.DiceTerm{Count: 1, Faces: 100},
					},
				},
			},
		},
		{
			Name: "unary operations aka signs", Term: ast.SubtractTerm{
				Left: ast.AddTerm{
					Left:  ast.IntTerm{Value: 1},
					Right: ast.IntTerm{Value: 0},
				},
				Right: ast.IntTerm{Value: 0},
			},
		},
	}}, formula)

	formula, err = parser.Parse("2d4 + d 20 - -1, keyword 4 D8")
	assert.EqualError(t, err, `line 1 column 26: expected "=", got "4"`)
	assert.Nil(t, formula)

	formula, err = parser.Parse("!")
	assert.EqualError(t, err, `line 1 column 1: expected integer or dice term or "(", got "!"`)
	assert.Nil(t, formula)

	formula, err = parser.Parse("1 ! 2")
	assert.EqualError(t, err, `line 1 column 3: expected integer or dice term or "(", got "!"`)
	assert.Nil(t, formula)

	formula, err = parser.Parse("(5d8")
	assert.EqualError(t, err, `line 1 column 4: expected ")", got end of input`)
	assert.Nil(t, formula)
}
