package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

// Instruction structure for intermediate program
type Instruction struct {
	operator int
	operand  int
}

const (
	opNoop = iota
	opAddDp
	opAddVal
	opOut
	opIn
	opJmpZ
	opJmpNz
	opZero
	opMove
)

const dataSize int = math.MaxUint16 + 1

func compile(input io.ByteReader) (program []Instruction, err error) {
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
		switch chr {
		case '>':
			program = append(program, Instruction{opAddDp, 1})
			if program[pc-1].operator == opAddDp {
				program = program[:pc]
				pc--
				program[pc].operand++
			}
		case '<':
			program = append(program, Instruction{opAddDp, -1})
			if program[pc-1].operator == opAddDp {
				program = program[:pc]
				pc--
				program[pc].operand--
			}
		case '+':
			program = append(program, Instruction{opAddVal, 1})
			if program[pc-1].operator == opAddVal {
				program = program[:pc]
				pc--
				program[pc].operand++
			}
		case '-':
			program = append(program, Instruction{opAddVal, -1})
			if program[pc-1].operator == opAddVal {
				program = program[:pc]
				pc--
				program[pc].operand--
			}
		case '.':
			program = append(program, Instruction{opOut, 0})
		case ',':
			program = append(program, Instruction{opIn, 0})
		case '[':
			program = append(program, Instruction{opJmpZ, 0})
			jmpStack = append(jmpStack, pc)
		case ']':
			if len(jmpStack) == 0 {
				return nil, errors.New("compilation error: unbalanced braces")
			}
			jmpPc = jmpStack[len(jmpStack)-1]
			jmpStack = jmpStack[:len(jmpStack)-1]
			program = append(program, Instruction{opJmpNz, jmpPc})
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
		default:
			pc--
		}
		pc++
	}
	if len(jmpStack) != 0 {
		return nil, errors.New("compilation error: unexpected EOF")
	}
	return
}

func execute(program []Instruction, reader io.ByteReader) {
	data := make([]int16, dataSize)
	var dataPtr int = 0
	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case opAddDp:
			dataPtr += program[pc].operand
			dataPtr = dataPtr & 0xFFFF
		case opAddVal:
			data[dataPtr] += int16(program[pc].operand)
		case opOut:
			fmt.Printf("%c", data[dataPtr]&0xff)
		case opIn:
			readVal, err := reader.ReadByte()
			data[dataPtr] = int16(readVal)
			if err == io.EOF {
				data[dataPtr] = int16(-1)
			} else if err != nil {
				break
			}
		case opJmpZ:
			if data[dataPtr] == 0 {
				pc = program[pc].operand
			}
		case opJmpNz:
			if data[dataPtr] != 0 {
				pc = program[pc].operand
			}
		case opZero:
			data[dataPtr] = 0
		case opMove:
			destPtr := (program[pc].operand + dataPtr) & 0xFFFF
			data[destPtr] += data[dataPtr]
			data[dataPtr] = 0
		case opNoop:
			continue
		default:
			panic("Unknown operator.")
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s file\n", args[0])
		return
	}
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println("Error reading source file")
		return
	}
	program, err := compile(bufio.NewReader(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	input := bufio.NewReader(os.Stdin)
	if len(args) > 2 {
		file, err = os.Open(args[2])
		if err != nil {
			fmt.Println("Error reading input file")
			return
		}
		input = bufio.NewReader(file)
	}
	execute(program, input)
}
