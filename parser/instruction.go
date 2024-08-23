package parser

// Instruction structure for intermediate program
type Instruction struct {
	operator Opcode
	operand  int
}

// Opcode type for instructions.
type Opcode int

const (
	opNoop Opcode = iota
	opAddDp
	opAddVal
	opOut
	opIn
	opJmpZ
	opJmpNz
)

var instMap = map[byte]Instruction{
	'>': {opAddDp, 1},
	'<': {opAddDp, -1},
	'+': {opAddVal, 1},
	'-': {opAddVal, -1},
	'.': {opOut, 1},
	',': {opIn, 1},
	'[': {opJmpZ, 0},
	']': {opJmpNz, 0},
}

// NewInstruction created from a sourcecode byte
func NewInstruction(chr byte) Instruction {
	return instMap[chr]
}

// SameOp compares Instructions operator but not operand
func (inst Instruction) SameOp(instTwo Instruction) bool {
	return inst.operator == instTwo.operator
}
