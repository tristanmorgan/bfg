package parser

import (
	"bufio"
	"io"
	"math"
)

// DataSize is the size of the input data
const DataSize int = math.MaxUint16 + 1

// DataMask is the mask to match the size
const DataMask int = math.MaxUint16

// Number interface defines different types that can be used in execution
type Number interface {
	byte | int | int8 | int16 | int32 | int64
}

// Execute a compiled program
func Execute[T Number](data []T, program []Instruction, reader io.ByteReader, writer *bufio.Writer) []T {
	var dataPtr, operand, writeCount int = 0, 0, 0
	for pc := 0; pc < len(program); pc++ {
		operand = program[pc].operand
		switch program[pc].operator {
		case opAddDp:
			dataPtr = (operand + dataPtr) & DataMask
		case opAddVal:
			data[dataPtr] += T(operand)
		case opSetVal:
			data[dataPtr] = T(operand)
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
			data[dataPtr] = T(readVal)
			if err == io.EOF {
				data[dataPtr] = T(0)
				data[dataPtr]--
			} else if err != nil {
				break
			}
		case opSkip:
			for data[dataPtr] != 0 {
				dataPtr = (operand + dataPtr) & DataMask
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
			destPtr := (operand + dataPtr) & DataMask
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
