package parser

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestPrint(t *testing.T) {
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 5},
		{opSetVal, 0},
		{opAddVal, 5},
		{opMove, 2},
		{opJmpZ, 7},
		{opIn, 1},
		{opSkip, -3},
		{opSetVal, 2},
		{opMulVal, 2},
		{opNoop, 2},
		{opDupVal, 1},
		{opNoop, 2},
		{opJmpNz, 5},
		{opAddDp, 2},
		{opOut, 1},
		{opMovN, 2},
	}
	var buf bytes.Buffer
	outputBuf := bufio.NewWriter(&buf)
	Print(program, outputBuf)
	got := buf.String()
	want := " \n >>>>>\n [-]\n +++++\n [->>+<<]\n [\n\t ,\n\t [<<<]\n\t ++\n\t [->>++<<]\n\t [->+>+<<]\n ]\n >>\n .\n [->>-<<]\n"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
