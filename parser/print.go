package parser

import (
	"bufio"
	"fmt"
	"iter"
	"strings"
)

func repeatDirection(neg, pos string, vect int) string {
	if vect > 0 {
		return strings.Repeat(pos, vect)
	}
	return strings.Repeat(neg, -vect)
}

func instPrints(program []Instruction) iter.Seq[string] {
	return func(yield func(string) bool) {

		depth := 0
		indent := ""
		for pc, inst := range program {
			indent = strings.Repeat("\t", depth)
			str := ""
			switch inst.operator {
			case opAddDp:
				str = repeatDirection("<", ">", inst.operand)
			case opAddVal:
				str = repeatDirection("-", "+", inst.operand)
			case opSetVal:
				prefix := "[-]"
				if program[pc-1].IsZeroOp() && inst.operand != 0 {
					prefix = ""
				}
				str = prefix + repeatDirection("-", "+", inst.operand)
			case opOut:
				str = "."
			case opIn:
				str = ","
			case opJmpZ:
				str = "["
				depth++
			case opJmpNz:
				str = "]"
				depth--
				indent = strings.Repeat("\t", depth)
			case opMovN:
				str = "[-" + repeatDirection("<", ">", inst.operand) + "-" + repeatDirection(">", "<", inst.operand) + "]"
			case opSkip:
				str = "[" + repeatDirection("<", ">", inst.operand) + "]"
			case opDupVal, opMove:
				str = "[-" + repeatDirection("<", ">", inst.operand) + "+"
				for pc+1 < len(program) && program[pc+1].operator == opNoop {
					str = str + repeatDirection("<", ">", program[pc+1].operand-inst.operand) + "+"
					pc++
					inst = program[pc]
				}
				str = str + repeatDirection(">", "<", program[pc].operand) + "]"
			case opVec, opMulVal:
				multiplier := repeatDirection("-", "+", program[pc+1].operand)
				str = "[-" + repeatDirection("<", ">", inst.operand) + multiplier
				for pc+2 < len(program) && program[pc+2].operator == opNoop {
					str = str + repeatDirection("<", ">", program[pc+2].operand-inst.operand) + multiplier
					pc++
					inst = program[pc+1]
				}
				str = str + repeatDirection(">", "<", inst.operand) + "]"
			case opNoop:
				continue
			default:
				return
			}
			if !yield(indent + str) {
				return
			}
		}
	}
}

// Print pretty prints out the parsed program.
func Print(program []Instruction, writer *bufio.Writer) {
	for printout := range instPrints(program) {
		fmt.Fprintln(writer, printout)
	}
	writer.Flush()
}
