package parser

import (
	"errors"
	"io"
)

// Tokenise sourcecode into an array of operators
func Tokenise(input io.ByteReader) (program []Instruction, err error) {
	var pc, jmpPc int = 0, 0
	jmpStack := make([]int, 0)
	program = append(program, Instruction{opNoop, 0})
	pc++
	for {
		chr, err := input.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New("tokenisation read error")
		}
		program = append(program, NewInstruction(chr))
		switch program[pc].operator {
		case opNoop:
			program = program[:pc]
			pc--
		case opAddDp:
			if program[pc-1].SameOp(program[pc]) {
				program[pc-1].operand += program[pc].operand
				program = program[:pc]
				pc--
			}
		case opAddVal:
			if program[pc-1].SameOp(program[pc]) {
				program[pc-1].operand += program[pc].operand
				program = program[:pc]
				pc--
			}
		case opJmpZ:
			jmpStack = append(jmpStack, pc)
		case opJmpNz:
			if len(jmpStack) == 0 {
				return nil, errors.New("tokenisation error: unbalanced braces")
			}
			jmpPc = jmpStack[len(jmpStack)-1]
			jmpStack = jmpStack[:len(jmpStack)-1]
			program[pc].operand = jmpPc
			program[jmpPc].operand = pc
			if pc-jmpPc == 2 && program[pc-1].operator == opAddVal {
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opZero, 0})
			}
			if pc-jmpPc == 5 &&
				program[pc-4].Complement(program[pc-2]) &&
				program[pc-3].Complement(program[pc-1]) &&
				!program[pc-2].SameOp(program[pc-1]) {
				offset := program[pc-4].operand
				if program[pc-3].operator == opAddDp {
					offset = program[pc-3].operand
				}
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opMove, offset})
			}
		}
		pc++
	}
	if len(jmpStack) != 0 {
		return nil, errors.New("tokenisation error: unexpected EOF")
	}
	return
}
