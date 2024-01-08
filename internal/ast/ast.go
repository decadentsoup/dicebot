package ast

// We use a non-crypto rand since dice bots are a terrible option for cryptography anyway.
import "math/rand"

type Formula struct {
	Equations []Equation
}

type Equation struct {
	Name string
	Term Term
}

type TermKind int

const (
	Int TermKind = iota
	Dice
	Add
	Subtract
)

type Term interface{ Solve() int }

type MultiplyTerm struct{ Left, Right Term }

func (mulTerm MultiplyTerm) Solve() int {
	return mulTerm.Left.Solve() * mulTerm.Right.Solve()
}

type DivideTerm struct{ Left, Right Term }

func (divTerm DivideTerm) Solve() int {
	return divTerm.Left.Solve() / divTerm.Right.Solve()
}

type AddTerm struct{ Left, Right Term }

func (addTerm AddTerm) Solve() int {
	return addTerm.Left.Solve() + addTerm.Right.Solve()
}

type SubtractTerm struct{ Left, Right Term }

func (subTerm SubtractTerm) Solve() int {
	return subTerm.Left.Solve() - subTerm.Right.Solve()
}

type DiceTerm struct{ Count, Faces int }

func (diceTerm DiceTerm) Solve() int {
	total := 0

	for index := 0; index < diceTerm.Count; index++ {
		dieResult := rand.Intn(diceTerm.Faces) + 1 //nolint:gosec
		total += dieResult
	}

	return total
}

type IntTerm struct{ Value int }

func (intTerm IntTerm) Solve() int {
	return intTerm.Value
}
