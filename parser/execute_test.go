package parser

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
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
		{opAddDp, -1},
		{opDupVal, -1},
		{opNoop, -2},
		{opAddDp, -1},
		{opOut, 1},
		{opAddDp, -2},
		{opSetVal, 32},
		{opOut, 1},
		{opIn, 1},
		{opMovN, 2},
		{opAddDp, 1},
		{opVec, 2},
		{opNoop, 2},
		{opNoop, 3},
		{opNoop, 4},
	}
	startdata := make([]int, 65536)
	outputBuf := bufio.NewWriter(&bufferWriter{})
	inputBuf := bufio.NewReader(strings.NewReader("no input."))
	data := Execute(startdata, program, inputBuf, outputBuf)[:10]
	want := []int{0, 0, 0, 0, 0, -100, 20, 20, 20, 0}

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
-]
,+[-,+]`)))
	startdata := make([]int, 65536)
	outputBuf := bufio.NewWriter(&bufferWriter{})
	inputBuf := bufio.NewReader(iotest.ErrReader(io.EOF))
	data := Execute(startdata, program, inputBuf, outputBuf)[:10]
	want := []int{0, 0, 0, 170, 0, 0, 0, 0, 0, 0}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v want %v", data, want)
	}
}

type bufferWriter struct {
	buf []byte
}

func (writer *bufferWriter) Write(p []byte) (int, error) {
	writer.buf = append(writer.buf, p...)
	return len(p), nil
}

func BenchmarkExecute(b *testing.B) {
	sourceFile, err := os.Open("../sample/pi-digits.bf")
	if err != nil {
		b.Errorf("error opening program: err:")
	}
	buff := bufio.NewReader(sourceFile)
	program, _ := Tokenise(buff)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inputBuf := bufio.NewReader(strings.NewReader("80\n"))
		outputBuf := bufio.NewWriter(&bufferWriter{})
		startdata := make([]int, 65536)
		Execute(startdata, program, inputBuf, outputBuf)
	}
}
