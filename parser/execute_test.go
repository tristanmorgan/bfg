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
		{opNoop, 0},
		{opAddDp, 5},
		{opSetVal, 0},
		{opAddVal, 5},
		{opMove, 2},
		{opAddDp, 2},
		{opMulVal, -1},
		{opNoop, 2},
	}
	startdata := make([]int, 65536)
	outputBuf := bufio.NewWriter(os.Stdout)
	inputBuf := bufio.NewReader(strings.NewReader("no input."))
	data := Execute(startdata, program, inputBuf, outputBuf)[:10]
	want := []int{0, 0, 0, 0, 0, 0, 10, 0, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}

func TestTokeniseExecuteSmall(t *testing.T) {
	program, _ := Tokenise(bufio.NewReader(strings.NewReader(
		`>>>>>>>>>>
<+< <+< <+< <+ [>>]
>++++++++
[
  <
  <[>+<-]> Add the next bit
  [<++>-]  double the result
  >[<+>-]< and move the bit counter
-]`)))
	startdata := make([]int, 65536)
	outputBuf := bufio.NewWriter(os.Stdout)
	inputBuf := bufio.NewReader(strings.NewReader("no input."))
	data := Execute(startdata, program, inputBuf, outputBuf)[:10]
	want := []int{0, 0, 0, 170, 0, 0, 0, 0, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}
