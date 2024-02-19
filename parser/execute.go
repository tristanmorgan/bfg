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
	var dataPtr, writeCount int = 0, 0
	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case opAddDp:
			dataPtr = (program[pc].operand + dataPtr) & dataMask
		case opAddVal:
			data[dataPtr] += program[pc].operand
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
