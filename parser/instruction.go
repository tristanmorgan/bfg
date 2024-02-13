package parser

// Instruction structure for intermediate program
type Instruction struct {
	operator Opcode
	operand  int
}

type Opcode int

const (
	opNoop Opcode = iota
	opAddDp
	opAddVal
	opOut
	opIn
	opJmpZ
	opJmpNz
	opZero
	opMove
)
