package parser

import (
	"errors"
	"io"
)

// Tokenise sourcecode into an array of operators
func Tokenise(input io.ByteReader) (program []Instruction, err error) {
	var jmpStack = 0
	program = append(program, Instruction{opNoop, 0})
	for instruction := range Instructions(input) {
		program = append(program, instruction)
		switch instruction.operator {
		case opJmpZ:
			jmpStack++
		case opJmpNz:
			if jmpStack == 0 {
				return nil, errors.New("tokenisation error: unbalanced braces")
			}
			jmpStack--
		}
	}
	if jmpStack != 0 {
		return nil, errors.New("tokenisation error: unexpected EOF")
	}
	return
}
