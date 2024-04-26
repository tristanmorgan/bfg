package parser

import (
	"errors"
	"io"
)

// Tokenise sourcecode into an array of operators
func Tokenise(input io.ByteReader) (program []Instruction, err error) {
	var pc int = 0
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
		instruction := NewInstruction(chr)
		program = append(program, instruction)
		switch instruction.operator {
		case opNoop:
			program = program[:pc]
			pc--
		case opAddDp:
			if program[pc-1].SameOp(instruction) {
				program[pc-1].operand += instruction.operand
				program = program[:pc]
				pc--
			}
		case opAddVal:
			if program[pc-1].SameOp(instruction) ||
				program[pc-1].operator == opSetVal {
				program[pc-1].operand += instruction.operand
				program = program[:pc]
				pc--
			} else if program[pc-1].SameOp(NewInstruction(']')) {
				operand := instruction.operand
				program = program[:pc]
				program = append(program, Instruction{opSetVal, operand})
			}
		case opJmpZ:
			jmpStack = append(jmpStack, pc)
		case opJmpNz:
			if len(jmpStack) == 0 {
				return nil, errors.New("tokenisation error: unbalanced braces")
			}
			jmpPc := jmpStack[len(jmpStack)-1]
			jmpStack = jmpStack[:len(jmpStack)-1]
			program[pc].operand = jmpPc
			program[jmpPc].operand = pc
			if pc-jmpPc == 2 && program[pc-1].SameOp(NewInstruction('+')) {
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opSetVal, 0})
			} else if pc-jmpPc == 2 && program[pc-1].SameOp(NewInstruction('>')) {
				offset := program[pc-1].operand
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opSkip, offset})
			} else if pc-jmpPc == 5 && // [<<<+>>>-] or [-<<<+>>>]
				program[pc-4].Complement(program[pc-2]) &&
				program[pc-3].Complement(program[pc-1]) {
				offset := program[pc-4].operand
				if program[pc-3].SameOp(NewInstruction('>')) {
					offset = program[pc-3].operand
				}
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opMove, offset})
			} else if pc-jmpPc == 5 && // [<<++++>>-]
				program[pc-4].Complement(program[pc-2]) &&
				program[pc-3].SameOp(NewInstruction('+')) &&
				program[pc-1].SameOp(NewInstruction('-')) &&
				program[pc-1].operand == -1 {
				offset := program[pc-4].operand
				factor := program[pc-3].operand
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opMulVal, offset})
				pc++
				program = append(program, Instruction{opNoop, factor})
			}
		}
		pc++
	}
	if len(jmpStack) != 0 {
		return nil, errors.New("tokenisation error: unexpected EOF")
	}
	return
}
