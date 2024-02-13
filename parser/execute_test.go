package parser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestExecuteSmall(t *testing.T) {
	//func Execute(program []Instruction, reader io.ByteReader) []int16
	program := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opAddDp, 5},
		Instruction{opZero, 0},
		Instruction{opAddVal, 5},
		Instruction{opMove, 2},
		Instruction{opAddDp, 2},
	}
	data := Execute(program, bufio.NewReader(strings.NewReader("no input.")))[:10]
	want := []int16{0, 0, 0, 0, 0, 0, 0, 5, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}
