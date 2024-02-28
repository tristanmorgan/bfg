package parser

import (
	"bufio"
	"io"
	"math"
)

const dataSize int = math.MaxUint16 + 1
const dataMask int = math.MaxUint16

// Execute a compiled program
func Execute(program []Instruction, reader io.ByteReader, writer *bufio.Writer) []int {
	data := make([]int, dataSize)
	var dataPtr, operand, writeCount int = 0, 0, 0
	for pc := 0; pc < len(program); pc++ {
		operand = program[pc].operand
		switch program[pc].operator {
		case opAddDp:
			dataPtr = (operand + dataPtr) & dataMask
		case opAddVal:
			data[dataPtr] += operand
		case opSetVal:
			data[dataPtr] = operand
		case opOut:
			writer.WriteByte(byte(data[dataPtr]))
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
			data[dataPtr] = int(readVal)
			if err == io.EOF {
				data[dataPtr] = int(-1)
			} else if err != nil {
				break
			}
		case opJmpZ:
			if data[dataPtr] == 0 {
				pc = operand
			}
		case opJmpNz:
			if data[dataPtr] != 0 {
				pc = operand
			}
		case opMove:
			destPtr := (operand + dataPtr) & dataMask
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
