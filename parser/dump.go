package parser

import (
	"bufio"
	"fmt"
	"strings"
)

// Dump prints out the parsed program.
func Dump(program []Instruction, writer *bufio.Writer) {
	depth := 0
	startLoop := NewInstruction('[')
	endLoop := NewInstruction(']')
	for idx, inst := range program {
		if inst.SameOp(endLoop) {
			depth--
		}
		fmt.Fprintln(writer, idx, strings.Repeat("\t", depth), inst.String())
		if inst.SameOp(startLoop) {
			depth++
		}
	}
	writer.Flush()
}
