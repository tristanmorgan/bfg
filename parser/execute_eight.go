package parser

import (
	"bufio"
	"io"
)

// ExecuteEight a compiled program
func ExecuteEight(program []Instruction, reader io.ByteReader, writer *bufio.Writer) []byte {
	data := make([]byte, dataSize)
	var dataPtr, operand, writeCount int = 0, 0, 0
	for pc := 0; pc < len(program); pc++ {
		operand = program[pc].operand
		switch program[pc].operator {
		case opAddDp:
			dataPtr = (operand + dataPtr) & dataMask
		case opAddVal:
			data[dataPtr] += byte(operand)
		case opSetVal:
			data[dataPtr] = byte(operand)
		case opOut:
			writer.WriteByte(data[dataPtr])
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
			data[dataPtr] = readVal
			if err == io.EOF {
				data[dataPtr] = byte(255)
			} else if err != nil {
				break
			}
		case opSkip:
			for data[dataPtr] != 0 {
				dataPtr = (operand + dataPtr) & dataMask
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
