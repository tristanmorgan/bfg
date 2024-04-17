package parser

import (
	"fmt"
	"reflect"
	"testing"
)

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
		opSetVal,
		opOut,
		opIn,
		opJmpZ,
		opJmpNz,
		opMove,
		opSkip,
	}

	for row, rval := range opsList {
		t.Run(fmt.Sprintf("%d", row), func(t *testing.T) {
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

func TestComplement(t *testing.T) {
	opsList := []Opcode{
		opAddDp,
		opAddVal,
	}

	for row, rval := range opsList {
		for col, cval := range opsList {
			t.Run(fmt.Sprintf("%d-%d", row, col), func(t *testing.T) {
				for operand := range [6]int{} {
					rinst := Instruction{rval, operand}
					cinst := Instruction{cval, operand - 6}
					want := row == col && operand == 3

					if rinst.Complement(cinst) != want {
						t.Errorf("testing %v vs %v want %v", rinst, cinst, want)
					}
				}
			})
		}
	}

}
