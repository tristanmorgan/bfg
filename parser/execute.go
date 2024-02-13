package parser

import (
	"bufio"
	"io"
	"math"
)

const dataSize int = math.MaxUint16 + 1

// Execute a compiled program
func Execute(program []Instruction, reader io.ByteReader, writer *bufio.Writer) []int16 {
	data := make([]int16, dataSize)
	var dataPtr, writeCount int = 0, 0
	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case opAddDp:
			dataPtr += program[pc].operand
			dataPtr = dataPtr & 0xFFFF
		case opAddVal:
			data[dataPtr] += int16(program[pc].operand)
		case opOut:
			writer.WriteByte(byte(data[dataPtr] & 0xff))
			writeCount++
			if data[dataPtr] == 10 || writeCount > 200 {
				writer.Flush()
				writeCount = 0
			}
		case opIn:
			if writeCount > 0 {
				writer.Flush()
				writeCount = 0
			}
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
	return data
}
