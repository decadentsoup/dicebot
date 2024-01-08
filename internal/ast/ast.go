package ast

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

type Term any

type ExponentiateTerm struct{ Left, Right Term }

type MultiplyTerm struct{ Left, Right Term }

type DivideTerm struct{ Left, Right Term }

type AddTerm struct{ Left, Right Term }

type SubtractTerm struct{ Left, Right Term }

type DiceTerm struct{ Count, Faces int }

type IntTerm struct{ Value int }
