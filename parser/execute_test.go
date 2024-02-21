package parser

import (
	"bufio"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestExecuteSmall(t *testing.T) {
	program := []Instruction{
		Instruction{opNoop, 0},
		Instruction{opAddDp, 5},
		Instruction{opZero, 0},
		Instruction{opAddVal, 5},
		Instruction{opMove, 2},
		Instruction{opAddDp, 2},
	}
	outputBuf := bufio.NewWriter(os.Stdout)
	inputBuf := bufio.NewReader(strings.NewReader("no input."))
	data := Execute(program, inputBuf, outputBuf)[:10]
	want := []int16{0, 0, 0, 0, 0, 0, 0, 5, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}

func TestTokeniseExecuteSmall(t *testing.T) {
	program, _ := Tokenise(bufio.NewReader(strings.NewReader(
		`+> >+> >+> >+> >
>++++++++
[
  <
  <[>+<-]> Add the next bit
  [<++>-]  double the result
  >[<+>-]< and move the bit counter
-]`)))
	outputBuf := bufio.NewWriter(os.Stdout)
	inputBuf := bufio.NewReader(strings.NewReader("no input."))
	data := Execute(program, inputBuf, outputBuf)[:10]
	want := []int16{170, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}
