package parser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestCompileSmall(t *testing.T) {
	program, err := Compile(bufio.NewReader(strings.NewReader(">>>>>[-]+++++[->>+<<]>>")))
	want := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opAddDp, 5},
		Instruction{opZero, 0},
		Instruction{opAddVal, 5},
		Instruction{opMove, 2},
		Instruction{opAddDp, 2},
	}

	if !reflect.DeepEqual(program, want) {
		t.Errorf("got %v want %v", program, want)
	} else if err != nil {
		t.Errorf("Error thrown  %v", err)
	}
}

func TestCompileOpDp(t *testing.T) {
	program, err := Compile(bufio.NewReader(strings.NewReader(">>>>>>><<<<<<>")))
	want := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opAddDp, 2},
	}

	if !reflect.DeepEqual(program, want) {
		t.Errorf("got %v want %v", program, want)
	} else if err != nil {
		t.Errorf("Error thrown  %v", err)
	}
}

func TestCompileOpVal(t *testing.T) {
	program, err := Compile(bufio.NewReader(strings.NewReader("----++----++")))
	want := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opAddVal, -4},
	}

	if !reflect.DeepEqual(program, want) {
		t.Errorf("got %v want %v", program, want)
	} else if err != nil {
		t.Errorf("Error thrown  %v", err)
	}
}

func TestCompileOpMove(t *testing.T) {
	program, err := Compile(bufio.NewReader(strings.NewReader("[->>+<<][<<<+>>>-]")))
	want := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opMove, 2},
		Instruction{opMove, -3},
	}

	if !reflect.DeepEqual(program, want) {
		t.Errorf("got %v want %v", program, want)
	} else if err != nil {
		t.Errorf("Error thrown  %v", err)
	}
}

func TestCompileOpJmpZNZ(t *testing.T) {
	program, err := Compile(bufio.NewReader(strings.NewReader("[->>+>+<<<]")))
	want := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opJmpZ, 8},
		Instruction{opAddVal, -1},
		Instruction{opAddDp, 2},
		Instruction{opAddVal, 1},
		Instruction{opAddDp, 1},
		Instruction{opAddVal, 1},
		Instruction{opAddDp, -3},
		Instruction{opJmpNz, 1},
	}

	if !reflect.DeepEqual(program, want) {
		t.Errorf("got %v want %v", program, want)
	} else if err != nil {
		t.Errorf("Error thrown  %v", err)
	}
}
