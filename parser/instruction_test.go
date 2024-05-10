package parser

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

var opName = map[Opcode]string{
	opNoop:   "nop",
	opAddDp:  "ptr",
	opAddVal: "add",
	opOut:    "out",
	opIn:     "inp",
	opJmpZ:   "jmp",
	opJmpNz:  "jnz",
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

	for idx, val := range []byte(sourceCode) {
		t.Run(program[idx].String(), func(t *testing.T) {
			got := NewInstruction(val)
			want := program[idx]

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestSameOp(t *testing.T) {
	opsList := []Opcode{
		opNoop,
		opAddDp,
		opAddVal,
		opOut,
		opIn,
		opJmpZ,
		opJmpNz,
	}

	for row, rval := range opsList {
		t.Run(opName[rval], func(t *testing.T) {
			for col, cval := range opsList {
				rinst := Instruction{rval, col}
				cinst := Instruction{cval, row}
				want := row == col

				if rinst.SameOp(cinst) != want {
					t.Errorf("testing %v vs %v want %v", rinst, cinst, want)
				}
			}
		})
	}

}
