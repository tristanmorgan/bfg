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
	opDupVal: "dup",
	opIn:     "inp",
	opJmpNz:  "jnz",
	opJmpZ:   "jmp",
	opMovN:   "nmv",
	opMove:   "mov",
	opMulVal: "mul",
	opOut:    "out",
	opSetVal: "set",
	opSkip:   "skp",
	opVec:    "vec",
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

func TestIsZeroOp(t *testing.T) {
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 1},
		{opAddVal, 1},
		{opSetVal, 0},
		{opSetVal, 1},
		{opOut, 1},
		{opIn, 1},
		{opJmpZ, 0},
		{opJmpNz, 0},
		{opMove, 1},
		{opSkip, 1},
		{opMovN, 1},
		{opMulVal, 1},
		{opDupVal, 1},
		{opVec, 1},
	}
	want := []bool{
		true,
		false,
		false,
		true,
		false,
		false,
		false,
		false,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
	}

	for idx, val := range program {
		if want[idx] != val.IsZeroOp() {
			t.Errorf("testing %v got %v want %v", val, val.IsZeroOp(), want[idx])
		}
	}
}
func TestSameOp(t *testing.T) {
	opsList := []Opcode{
		opNoop,
		opAddDp,
		opAddVal,
		opSetVal,
		opOut,
		opIn,
		opJmpZ,
		opJmpNz,
		opMove,
		opMovN,
		opSkip,
		opMulVal,
		opDupVal,
		opVec,
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
