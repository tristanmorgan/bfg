package parser

import (
	"errors"
	"io"
)

// Compile sourcecode into a program
func Compile(input io.ByteReader) (program []Instruction, err error) {
	var pc, jmpPc int = 0, 0
	jmpStack := make([]int, 0)
	program = append(program, Instruction{opNoop, 0})
	pc++
	for {
		chr, err := input.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New("compilation read error")
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
				return nil, errors.New("compilation error: unbalanced braces")
			}
			jmpPc = jmpStack[len(jmpStack)-1]
			jmpStack = jmpStack[:len(jmpStack)-1]
			program[pc].operand = jmpPc
			program[jmpPc].operand = pc
			if pc-jmpPc == 2 && program[pc-1].operator == opAddVal {
				pc--
				pc--
				program = program[:pc]
				program = append(program, Instruction{opZero, 0})
			}
			if pc-jmpPc == 5 &&
				program[pc-4].operator == opAddVal &&
				program[pc-3].operator == opAddDp &&
				program[pc-2].operator == opAddVal &&
				program[pc-1].operator == opAddDp &&
				program[pc-4].operand == -1 &&
				program[pc-2].operand == 1 &&
				program[pc-3].operand+program[pc-1].operand == 0 {
				offset := program[pc-3].operand
				pc -= 5
				program = program[:pc]
				program = append(program, Instruction{opMove, offset})
			}
			if pc-jmpPc == 5 &&
				program[pc-4].operator == opAddDp &&
				program[pc-3].operator == opAddVal &&
				program[pc-2].operator == opAddDp &&
				program[pc-1].operator == opAddVal &&
				program[pc-3].operand == 1 &&
				program[pc-1].operand == -1 &&
				program[pc-4].operand+program[pc-2].operand == 0 {
				offset := program[pc-4].operand
				pc -= 5
				program = program[:pc]
				program = append(program, Instruction{opMove, offset})
			}
		}
		pc++
	}
	if len(jmpStack) != 0 {
		return nil, errors.New("compilation error: unexpected EOF")
	}
	return
}
