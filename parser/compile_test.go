package parser

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	table := []struct {
		name    string
		code    string
		program []Instruction
	}{
		{
			"small_prog",
			">>>>>[-]zero+++++[->>+<<]move>>",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opAddDp, 5},
				Instruction{opZero, 0},
				Instruction{opAddVal, 5},
				Instruction{opMove, 2},
				Instruction{opAddDp, 2},
			},
		},
		{
			"op_dp",
			">>>>>>><<<<<<>",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opAddDp, 2},
			},
		},
		{
			"op_val",
			"----++----++",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opAddVal, -4},
			},
		},
		{
			"op_move",
			"[->>+<<][<<<+>>>-]",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opMove, 2},
				Instruction{opMove, -3},
			},
		},
		{
			"op_jmp_z_nz",
			"[->>+>+<<<]",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opJmpZ, 8},
				Instruction{opAddVal, -1},
				Instruction{opAddDp, 2},
				Instruction{opAddVal, 1},
				Instruction{opAddDp, 1},
				Instruction{opAddVal, 1},
				Instruction{opAddDp, -3},
				Instruction{opJmpNz, 1},
			},
		},
		{
			"op_nested",
			"[[[[[[[,]]]]]]][comment.]",
			[]Instruction{
				Instruction{opNoop, 0},
				Instruction{opJmpZ, 15},
				Instruction{opJmpZ, 14},
				Instruction{opJmpZ, 13},
				Instruction{opJmpZ, 12},
				Instruction{opJmpZ, 11},
				Instruction{opJmpZ, 10},
				Instruction{opJmpZ, 9},
				Instruction{opIn, 0},
				Instruction{opJmpNz, 7},
				Instruction{opJmpNz, 6},
				Instruction{opJmpNz, 5},
				Instruction{opJmpNz, 4},
				Instruction{opJmpNz, 3},
				Instruction{opJmpNz, 2},
				Instruction{opJmpNz, 1},
				Instruction{opJmpZ, 18},
				Instruction{opOut, 0},
				Instruction{opJmpNz, 16},
			},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			got, err := Compile(bufio.NewReader(strings.NewReader(v.code)))
			want := v.program

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			} else if err != nil {
				t.Errorf("Error thrown  %v", err)
			}
		})
	}
}

func TestCompileError(t *testing.T) {
	table := []struct {
		name string
		code string
		err  error
	}{
		{
			"too_many_open",
			"[[[",
			errors.New("compilation error: unexpected EOF"),
		},
		{
			"too_many_close",
			"]]]",
			errors.New("compilation error: unbalanced braces"),
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			_, got := Compile(bufio.NewReader(strings.NewReader(v.code)))
			want := v.err

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}
