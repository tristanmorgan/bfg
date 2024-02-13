package parser

import (
	"fmt"
	"io"
	"math"
)

const dataSize int = math.MaxUint16 + 1

func Execute(program []Instruction, reader io.ByteReader) {
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
