package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tristanmorgan/bfg/parser"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s file\n", args[0])
		return
	}
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Println("Error reading source file")
		return
	}
	program, err := parser.Compile(bufio.NewReader(file))
	if err != nil {
		fmt.Println(err)
		return
	}
	input := bufio.NewReader(os.Stdin)
	if len(args) > 2 {
		file, err = os.Open(args[2])
		if err != nil {
			fmt.Println("Error reading input file")
			return
		}
		input = bufio.NewReader(file)
	}
	parser.Execute(program, input)
}
