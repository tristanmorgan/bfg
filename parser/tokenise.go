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
			} else if program[pc-1].IsZeroOp() {
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
			switch {
			case program[jmpPc-1].IsZeroOp():
				pc = jmpPc
				program = program[:pc]
				pc--
			case pc-jmpPc == 2 &&
				(program[pc-1].SameOp(NewInstruction('+')) ||
					program[pc-1].operator == opSetVal):
				pc = jmpPc
				if program[jmpPc-1].SameOp(NewInstruction('+')) {
					pc--
				}
				program = program[:pc]
				program = append(program, Instruction{opSetVal, 0})
			case pc-jmpPc == 2 &&
				(program[pc-1].SameOp(NewInstruction('>')) ||
					program[pc-1].operator == opSkip):
				offset := program[pc-1].operand
				pc = jmpPc
				program = program[:pc]
				program = append(program, Instruction{opSkip, offset})
			case pc-jmpPc == 5: // looking for opMulVal and opMove
				pointers, factors, ok := evalFactors(program[jmpPc+1 : pc])

				if ok && factors[0] == 1 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opMove, pointers[0]})
				} else if ok && factors[0] != 0 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opMulVal, pointers[0]})
					pc++
					program = append(program, Instruction{opNoop, factors[0]})
				}
			case pc-jmpPc == 7: //looking for opDupVal
				pointers, factors, ok := evalFactors(program[jmpPc+1 : pc])
				if ok && factors[0] == 1 && factors[1] == 1 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opDupVal, pointers[0]})
					pc++
					program = append(program, Instruction{opNoop, pointers[1]})
				}
			}
		}
		pc++
	}
	if len(jmpStack) != 0 {
		return nil, errors.New("tokenisation error: unexpected EOF")
	}
	return
}

func evalFactors(program []Instruction) (pointers, factors []int, ok bool) {
	ok = false
	collect := 0
	pointers = make([]int, 0)
	factors = make([]int, 0)
	for _, inst := range program {
		if inst.SameOp(NewInstruction('>')) {
			collect += inst.operand
		} else if collect == 0 && inst.Complement(NewInstruction('+')) {
			ok = true
		} else if inst.SameOp(NewInstruction('-')) {
			pointers = append(pointers, collect)
			factors = append(factors, inst.operand)
		} else {
			return pointers, factors, false
		}
	}
	if collect != 0 || !ok {
		return pointers, factors, false
	}
	return
}
