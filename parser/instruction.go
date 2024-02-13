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

func (inst Instruction) String() string {
	opName := [...]string{
		"nop",
		"ptr",
		"add",
		"out",
		"in",
		"jmpz",
		"jmpnz",
		"zero",
		"mov",
	}
	return opName[inst.operator]
}

func (inst Instruction) SameOp(instTwo Instruction) bool {
	return inst.operator == instTwo.operator
}
