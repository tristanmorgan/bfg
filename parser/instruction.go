package parser

import (
	"fmt"
)

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
	opZero
	opMove
)

// String representation of Instruction
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
	return fmt.Sprintf("%s: %v", opName[inst.operator], inst.operand)
}

// NewInstruction created from a sourcecode byte
func NewInstruction(chr byte) Instruction {
	switch chr {
	case '>':
		return Instruction{opAddDp, 1}
	case '<':
		return Instruction{opAddDp, -1}
	case '+':
		return Instruction{opAddVal, 1}
	case '-':
		return Instruction{opAddVal, -1}
	case '.':
		return Instruction{opOut, 1}
	case ',':
		return Instruction{opIn, 1}
	case '[':
		return Instruction{opJmpZ, 0}
	case ']':
		return Instruction{opJmpNz, 0}
	default:
		return Instruction{opNoop, 0}
	}
}

// SameOp compares Instructions operator but not operand
func (inst Instruction) SameOp(instTwo Instruction) bool {
	return inst.operator == instTwo.operator
}

// Complement compares Instructions operator and operand
func (inst Instruction) Complement(instTwo Instruction) bool {
	return inst.SameOp(instTwo) && inst.operand+instTwo.operand == 0
}
