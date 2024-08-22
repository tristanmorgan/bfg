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
			">>>>>++[-]zero+++++>+++++[[->>+<<]move]>>[[>>-<<-]movn]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 5},
				{opSetVal, 5},
				{opAddDp, 1},
				{opAddVal, 5},
				{opMove, 2},
				{opAddDp, 2},
				{opMovN, 2},
			},
		},
		{
			"op_mul",
			"+[<++++++>-]>[+>>+++<<]",
			[]Instruction{
				{opNoop, 0},
				{opSetVal, 1},
				{opMulVal, -1}, // dest value pointer
				{opNoop, 6},    // multiplication factor
				{opAddDp, 1},
				{opMulVal, 2}, // dest value pointer
				{opNoop, -3},  // multiplication factor
			},
		},
		{
			"op_dp",
			">>>>>>><<<<<<>",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 2},
			},
		},
		{
			"op_val",
			"----++->---++",
			[]Instruction{
				{opNoop, 0},
				{opSetVal, -3},
				{opAddDp, 1},
				{opAddVal, -1},
			},
		},
		{
			"op_skip",
			">[>>>>>]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opSkip, 5},
			},
		},
		{
			"op_move",
			">[->>+<<]>[<<<+>>>-]++[[>-<-]]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opMove, 2},
				{opAddDp, 1},
				{opMove, -3},
				{opSetVal, 2},
				{opMovN, 1},
			},
		},
		{
			"op_vec",
			">[->>+++>+++<<<]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opVec, 2},
				{opNoop, 3},
				{opNoop, 3},
			},
		},
		{
			"op_dup",
			">[->>+>+<<<]>[>>+<<-<<+>>]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opDupVal, 2},
				{opNoop, 3},
				{opAddDp, 1},
				{opDupVal, 2},
				{opNoop, -2},
			},
		},
		{
			"op_jmp_z_nz",
			">[->>,>+<<<]>[<<+>>--]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opJmpZ, 9},
				{opAddVal, -1},
				{opAddDp, 2},
				{opIn, 1},
				{opAddDp, 1},
				{opAddVal, 1},
				{opAddDp, -3},
				{opJmpNz, 2},
				{opAddDp, 1},
				{opJmpZ, 16},
				{opAddDp, -2},
				{opAddVal, 1},
				{opAddDp, 2},
				{opAddVal, -2},
				{opJmpNz, 11},
			},
		},
		{
			"op_nested",
			">[[[[[[[>]]]]]]][comment.]+++[[-]]",
			[]Instruction{
				{opNoop, 0},
				{opAddDp, 1},
				{opSkip, 1},
				{opSetVal, 0},
			},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			got, err := Tokenise(bufio.NewReader(strings.NewReader(v.code)))
			want := v.program

			if err != nil {
				t.Errorf("Error thrown  %v", err)
			} else if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
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
