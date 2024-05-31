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

func instPrint(inst, lastInst Instruction) string {
	switch inst.operator {
	case opAddDp:
		return repeatDirection("<", ">", inst.operand)
	case opAddVal:
		return repeatDirection("-", "+", inst.operand)
	case opSetVal:
		prefix := "[-]"
		if lastInst.IsZeroOp() && inst.operand != 0 {
			prefix = ""
		}
		return prefix + repeatDirection("-", "+", inst.operand)
	case opOut:
		return "."
	case opIn:
		return ","
	case opJmpZ:
		return "["
	case opJmpNz:
		return "]"
	case opMove:
		return "[-" + repeatDirection("<", ">", inst.operand) + "+" + repeatDirection(">", "<", inst.operand) + "]"
	case opSkip:
		return "[" + repeatDirection("<", ">", inst.operand) + "]"
	case opMulVal:
		return ""
	case opDupVal:
		return ""
	case opNoop:
		if lastInst.operator == opMulVal {
			multiplier := repeatDirection("-", "+", inst.operand)
			return "[-" + repeatDirection("<", ">", lastInst.operand) + multiplier + repeatDirection(">", "<", lastInst.operand) + "]"
		} else if lastInst.operator == opDupVal {
			return "[-" + repeatDirection("<", ">", lastInst.operand) + "+" + repeatDirection("<", ">", inst.operand-lastInst.operand) + "+" + repeatDirection(">", "<", inst.operand) + "]"
		}
		return ""
	default:
		return ""
	}
}

// Print pretty prints out the parsed program.
func Print(program []Instruction, writer *bufio.Writer) {
	depth := 0
	startLoop := NewInstruction('[')
	endLoop := NewInstruction(']')
	lastInst := NewInstruction('!')
	for _, inst := range program {
		if inst.operator == opMulVal {
			lastInst = inst
			continue
		}
		if inst.SameOp(endLoop) {
			depth--
		}
		fmt.Fprintln(writer, strings.Repeat("\t", depth), instPrint(inst, lastInst))
		if inst.SameOp(startLoop) {
			depth++
		}
		lastInst = inst
	}
	writer.Flush()
}
