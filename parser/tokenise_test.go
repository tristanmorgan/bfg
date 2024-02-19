package parser

import (
	"bufio"
	"errors"
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
				Instruction{opNoop, 0},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
				Instruction{opJmpZ, 0},
				Instruction{opAddVal, -1},
				Instruction{opJmpNz, 0},
				Instruction{opAddVal, 1},
				Instruction{opAddVal, 1},
				Instruction{opAddVal, 1},
				Instruction{opAddVal, 1},
				Instruction{opAddVal, 1},
				Instruction{opJmpZ, 0},
				Instruction{opAddVal, -1},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
				Instruction{opAddVal, 1},
				Instruction{opAddDp, -1},
				Instruction{opAddDp, -1},
				Instruction{opJmpNz, 0},
				Instruction{opAddDp, 1},
				Instruction{opAddDp, 1},
			},
		},
		{
			"op_nested",
			"[[[[[[[,]]]]]]][.]",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opJmpZ, 0},
				Instruction{opIn, 1},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpNz, 0},
				Instruction{opJmpZ, 0},
				Instruction{opOut, 1},
				Instruction{opJmpNz, 0},
			},
		},
		{
			"op_noops",
			"no_code",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
				Instruction{opNoop, 0},
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
