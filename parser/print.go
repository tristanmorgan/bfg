package parser

import (
	"bufio"
	"fmt"
	"iter"
	"strings"
)

func repeatDirection(neg, pos string, vect int) string {
	if vect > 0 {
		return pos
	}
	return neg
}

func instPrints(program []Instruction) iter.Seq[string] {
	return func(yield func(string) bool) {

		depth := 0
		indent := ""
		for _, inst := range program {
			indent = strings.Repeat("\t", depth)
			str := ""
			switch inst.operator {
			case opAddDp:
				str = repeatDirection("<", ">", inst.operand)
			case opAddVal:
				str = repeatDirection("-", "+", inst.operand)
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
