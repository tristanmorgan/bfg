package parser

import (
	"errors"
	"io"
)

// Tokenise sourcecode into an array of operators
func Tokenise(input io.ByteReader) (program []Instruction, err error) {
	var pc, jmpStack int = 0, 0
	program = append(program, Instruction{opNoop, 0})
	pc++
	for {
		chr, err := input.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New("tokenisation read error")
		}
		program = append(program, NewInstruction(chr))
		switch program[pc].operator {
		case opJmpZ:
			jmpStack++
		case opJmpNz:
			if jmpStack == 0 {
				return nil, errors.New("tokenisation error: unbalanced braces")
			}
			jmpStack--
		}
		pc++
	}
	if jmpStack != 0 {
		return nil, errors.New("tokenisation error: unexpected EOF")
	}
	return
}
