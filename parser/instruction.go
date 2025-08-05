// Package parser provides the parsing and execution of BF code.
package parser

import (
	"errors"
	"io"
	"iter"
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
	opDupVal
	opIn
	opJmpNz
	opJmpZ
	opMovN
	opMove
	opMulVal
	opOut
	opSetVal
	opSkip
	opVec
)

var instMap = map[byte]Instruction{
	'+': {opAddVal, 1},
	',': {opIn, 1},
	'-': {opAddVal, -1},
	'.': {opOut, 1},
	'<': {opAddDp, -1},
	'>': {opAddDp, 1},
	'[': {opJmpZ, 0},
	']': {opJmpNz, 0},
}

// Instructions returns an iterator of the instructions.
func Instructions(input io.ByteReader) iter.Seq[Instruction] {
	return func(yield func(Instruction) bool) {

		for {
			chr, err := input.ReadByte()
			if errors.Is(err, io.EOF) {
				break
			}
			if !yield(NewInstruction(chr)) {
				return
			}
		}
	}
}

// NewInstruction created from a sourcecode byte
func NewInstruction(chr byte) Instruction {
	return instMap[chr]
}

// SameOp compares Instructions operator but not operand
func (inst Instruction) SameOp(instTwo Instruction) bool {
	return inst.operator == instTwo.operator
}

// IsZeroOp returns true for ops that have left the pointer on a zero
func (inst Instruction) IsZeroOp() bool {
	return inst.operator != opAddDp &&
		inst.operator != opAddVal &&
		inst.operator != opJmpZ &&
		inst.operator != opOut &&
		inst.operator != opIn &&
		(inst.operator != opSetVal || inst.operand == 0)
}
