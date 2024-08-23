package parser

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTokenise(t *testing.T) {
	table := []struct {
		name    string
		code    string
		program []Instruction
	}{
		{
			"small_prog",
			">>>>>[-]+++++[->>+<<]>>",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opAddDp, 1},
				{opAddDp, 1},
				{opAddDp, 1},
				{opAddDp, 1},
				{opJmpZ, 0},
				{opAddVal, -1},
				{opJmpNz, 0},
				{opAddVal, 1},
				{opAddVal, 1},
				{opAddVal, 1},
				{opAddVal, 1},
				{opAddVal, 1},
				{opJmpZ, 0},
				{opAddVal, -1},
				{opAddDp, 1},
				{opAddDp, 1},
				{opAddVal, 1},
				{opAddDp, -1},
				{opAddDp, -1},
				{opJmpNz, 0},
				{opAddDp, 1},
				{opAddDp, 1},
			},
		},
		{
			"op_nested",
			">[[[[[[[,]]]]]]][.]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opJmpZ, 0},
				{opIn, 1},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpNz, 0},
				{opJmpZ, 0},
				{opOut, 1},
				{opJmpNz, 0},
			},
		},
		{
			"op_noops",
			"no_code",
			[]Instruction{
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
				{opNoop, 0},
			},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			got, err := Tokenise(bufio.NewReader(strings.NewReader(v.code)))
			want := v.program

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			} else if err != nil {
				t.Errorf("Error thrown  %v", err)
			}
		})
	}
}

func TestTokeniseError(t *testing.T) {
	table := []struct {
		name string
		code string
		err  error
	}{
		{
			"too_many_open",
			"[[[",
			errors.New("tokenisation error: unexpected EOF"),
		},
		{
			"too_many_close",
			"]]]",
			errors.New("tokenisation error: unbalanced braces"),
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			_, got := Tokenise(bufio.NewReader(strings.NewReader(v.code)))
			want := v.err

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func BenchmarkTokenise(b *testing.B) {
	sourceFile, err := os.Open("../sample/pi-digits.bf")
	if err != nil {
		b.Errorf("error opening program: err:")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buff := bufio.NewReader(sourceFile)
		Tokenise(buff)
	}
}
