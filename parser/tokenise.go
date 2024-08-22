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

			redOk, reduceInst := reduceInst(program[pc-1])
			switch {
			case program[jmpPc-1].IsZeroOp():
				pc = jmpPc
				program = program[:pc]
				pc--
			case pc-jmpPc == 2 && redOk:
				pc = jmpPc
				if (program[jmpPc-1].operator == opAddVal ||
					program[jmpPc-1].operator == opSetVal) &&
					reduceInst.operator == opSetVal {
					pc--
				}
				program = program[:pc]
				program = append(program, reduceInst)
			case pc-jmpPc == 5: // looking for opMulVal and opMove
				pointers, factors, ok := evalFactors(program[jmpPc+1 : pc])

				if ok && factors[0] == 1 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opMove, pointers[0]})
				} else if ok && factors[0] == -1 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opMovN, pointers[0]})
				} else if ok && factors[0] != 0 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opMulVal, pointers[0]})
					pc++
					program = append(program, Instruction{opNoop, factors[0]})
				}
			case pc-jmpPc == 7 || pc-jmpPc == 8: //looking for opDupVal
				pointers, factors, ok := evalFactors(program[jmpPc+1 : pc])
				if ok && factors[0] == 1 && factors[1] == 1 {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opDupVal, pointers[0]})
					pc++
					program = append(program, Instruction{opNoop, pointers[1]})
				} else if ok && factors[0] == factors[1] {
					pc = jmpPc
					program = program[:pc]
					program = append(program, Instruction{opVec, pointers[0]})
					pc++
					program = append(program, Instruction{opNoop, pointers[1]})
					pc++
					program = append(program, Instruction{opNoop, factors[1]})
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
	invert := false
	collect := 0
	pointers = make([]int, 0)
	factors = make([]int, 0)
	for _, inst := range program {
		if inst.SameOp(NewInstruction('>')) {
			collect += inst.operand
		} else if collect == 0 && inst.SameOp(NewInstruction('-')) {
			ok = inst.operand == 1 || inst.operand == -1
			invert = inst.operand == 1
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
	if invert {
		for i := range len(factors) {
			factors[i] = 0 - factors[i]
		}
	}
	return
}

func reduceInst(inst Instruction) (bool, Instruction) {
	switch inst.operator {
	case opAddDp:
		return true, Instruction{opSkip, inst.operand}
	case opSkip:
		return true, inst
	case opAddVal:
		if inst.operand == 1 || inst.operand == -1 {
			return true, Instruction{opSetVal, 0}
		}
		return false, Instruction{opNoop, 0}
	case opSetVal:
		if inst.operand == 0 {
			return true, inst
		}
		return false, Instruction{opNoop, 0}
	case opMove:
		return true, inst
	case opMovN:
		return true, inst
	default:
		return false, Instruction{opNoop, 0}
	}
}
