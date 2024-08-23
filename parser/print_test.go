package parser

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"
)

var testprogram = `>
>
>
>
+
+
+
+
[
	,
]
>
>
`

func TestPrint(t *testing.T) {
	program := []Instruction{
		{opNoop, 0},
		{opAddDp, 1},
		{opAddDp, 1},
		{opAddDp, 1},
		{opAddDp, 1},
		{opAddVal, 1},
		{opAddVal, 1},
		{opAddVal, 1},
		{opAddVal, 1},
		{opJmpZ, 0},
		{opIn, 1},
		{opJmpNz, 0},
		{opAddDp, 1},
		{opAddDp, 1},
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
