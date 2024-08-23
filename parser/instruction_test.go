package parser

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

var opName = map[Opcode]string{
	opNoop:   "nop",
	opAddDp:  "ptr",
	opAddVal: "add",
	opIn:     "inp",
	opJmpNz:  "jnz",
	opJmpZ:   "jmp",
	opOut:    "out",
}

func (inst Instruction) String() string {
	return fmt.Sprintf("%s:%v", opName[inst.operator], inst.operand)
}

func TestInstructionSize(t *testing.T) {
	sizes := []uintptr{2, 4, 8, 16, 32, 64, 128}
	ok := false
	got := unsafe.Sizeof(NewInstruction('-'))
	for _, s := range sizes {
		if s == got {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("size not in array: got %v", got)
	}
}

func TestNewInstruction(t *testing.T) {
	sourceCode := "g><+-.,[]"
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 1},
		{opAddDp, -1},
		{opAddVal, 1},
		{opAddVal, -1},
		{opOut, 1},
		{opIn, 1},
		{opJmpZ, 0},
		{opJmpNz, 0},
	}

	idx := 0
	for got := range Instructions(bufio.NewReader(strings.NewReader(sourceCode))) {
		t.Run(program[idx].String(), func(t *testing.T) {
			want := program[idx]

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			}
		})
		idx++
	}
}
