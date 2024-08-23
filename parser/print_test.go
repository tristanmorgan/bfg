package parser

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var testprogram = ` >>>>>
 [-]+++++
 [->>+<<]
 >
 [
	 ,
	 [<<<]
	 ++
	 [->>++<<]
	 >>
	 [->+>+<<]
	 >
	 [->+++>+++<<]
 ]
 >>
 .
 [->>-<<]
`

func TestPrint(t *testing.T) {
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 5},
		{opSetVal, 5},
		{opMove, 2},
		{opAddDp, 1},
		{opJmpZ, 7},
		{opIn, 1},
		{opSkip, -3},
		{opSetVal, 2},
		{opMulVal, 2},
		{opNoop, 2},
		{opAddDp, 2},
		{opDupVal, 1},
		{opNoop, 2},
		{opAddDp, 1},
		{opVec, 1},
		{opNoop, 3},
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
	want := testprogram

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestParseThenPrint(t *testing.T) {
	want := testprogram

	program, _ := Tokenise(bufio.NewReader(strings.NewReader(testprogram)))

	var buf bytes.Buffer
	outputBuf := bufio.NewWriter(&buf)
	Print(program, outputBuf)
	got := buf.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
