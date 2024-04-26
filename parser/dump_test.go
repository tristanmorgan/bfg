package parser

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestDump(t *testing.T) {
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 5},
		{opSetVal, 0},
		{opAddVal, 5},
		{opMove, 2},
		{opJmpZ, 7},
		{opIn, 1},
		{opJmpNz, 5},
		{opAddDp, 2},
	}
	var buf bytes.Buffer
	outputBuf := bufio.NewWriter(&buf)
	Dump(program, outputBuf)
	got := buf.String()
	want := "0  nop: 0\n1  ptr: 5\n2  set: 0\n3  add: 5\n4  mov: 2\n5  jmp: 7\n6 \t inp: 1\n7  jnz: 5\n8  ptr: 2\n"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
