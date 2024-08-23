package parser

import (
	"bufio"
	"fmt"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func repeatDirection(neg, pos string, vect int) string {
	if vect > 0 {
		return strings.Repeat(pos, vect)
	}
	return strings.Repeat(neg, abs(vect))
}

func instPrint(inst Instruction) string {
	switch inst.operator {
	case opAddDp:
		return repeatDirection("<", ">", inst.operand)
	case opAddVal:
		return repeatDirection("-", "+", inst.operand)
	case opOut:
		return "."
	case opIn:
		return ","
	case opJmpZ:
		return "["
	case opJmpNz:
		return "]"
	default:
		return ""
	}
}

// Print pretty prints out the parsed program.
func Print(program []Instruction, writer *bufio.Writer) {
	depth := 0
	startLoop := NewInstruction('[')
	endLoop := NewInstruction(']')
	for _, inst := range program {
		if inst.SameOp(endLoop) {
			depth--
		}
		printout := instPrint(inst)
		if printout != "" {
			fmt.Fprintln(writer, strings.Repeat("\t", depth), printout)
		}
		if inst.SameOp(startLoop) {
			depth++
		}
	}
	writer.Flush()
}
