package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tristanmorgan/bfg/parser"
)

func main() {
	version := flag.Bool("version", false, "display version")

	flag.Usage = func() {
		fmt.Printf("Usage:\n  %s [option] source.bf [input]\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *version {
		fmt.Println("Version 0.0.1")
		os.Exit(0)
	}

	var sourceBuf, inputBuf io.ByteReader
	source := flag.Arg(0)
	input := flag.Arg(1)
	if source == "" {
		fmt.Println("please provide a source file")
		os.Exit(1)
	} else {
		sourceFile, err := os.Open(source)
		if err != nil {
			fmt.Println("error opening source file: err:", err)
			os.Exit(1)
		}
		sourceBuf = bufio.NewReader(sourceFile)
	}
	if input == "" {
		inputBuf = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(input)
		if err != nil {
			fmt.Println("error opening input file: err:", err)
			os.Exit(1)
		}
		inputBuf = bufio.NewReader(file)
	}

	program, err := parser.Compile(sourceBuf)
	if err != nil {
		fmt.Println("error compiling program: err:", err)
		os.Exit(1)
	}
	parser.Execute(program, inputBuf)
}
