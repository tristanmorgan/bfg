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
	var dataPtr, writeCount = 0, 0
	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case opAddDp:
			dataPtr = (program[pc].operand + dataPtr) & DataMask
		case opAddVal:
			data[dataPtr] += T(program[pc].operand)
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
		case opJmpZ:
			if data[dataPtr] == 0 {
				depth := 1
				for depth >= 1 {
					pc++
					if program[pc].operator == opJmpZ {
						depth++
					} else if program[pc].operator == opJmpNz {
						depth--
					}
				}
			}
		case opJmpNz:
			if data[dataPtr] != 0 {
				depth := 1
				for depth >= 1 {
					pc--
					if program[pc].operator == opJmpNz {
						depth++
					} else if program[pc].operator == opJmpZ {
						depth--
					}
				}
			}
		case opNoop:
			continue
		default:
			panic("Unknown operator.")
		}
	}
	return data
}
