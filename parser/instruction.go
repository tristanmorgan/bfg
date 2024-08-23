// Package parser provides the parsing and execution of BF code.
package parser

import (
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
	opIn
	opJmpNz
	opJmpZ
	opOut
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
			if err == io.EOF {
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
