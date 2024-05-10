package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tristanmorgan/bfg/parser"
)

var (
	// Version is the main version number that is being run at the moment.
	Version = "0.0.1"

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "-dev" (in development), -"beta", "-rc1", etc.
	VersionPrerelease = "-dev"
)

func main() {
	version := flag.Bool("version", false, "display version")
	eight := flag.Bool("eight", false, "eight bit execution")
	dump := flag.Bool("dump", false, "dump parsed program")
	print := flag.Bool("print", false, "pretty print parsed program")

	flag.Usage = func() {
		fmt.Printf("Usage:\n  %s [option] source.bf [input]\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *version {
		fmt.Printf("Version: v%s%s\n", Version, VersionPrerelease)
		fmt.Println("https://github.com/tristanmorgan/bfg")
		os.Exit(0)
	}

	var sourceBuf, inputBuf io.ByteReader
	source := flag.Arg(0)
	input := flag.Arg(1)
	sourceBuf, err := inputReader(source, false)
	if err != nil {
		fmt.Println("error opening program: err:", err)
		os.Exit(1)
	}
	inputBuf, err = inputReader(input, true)
	if err != nil {
		fmt.Println("error opening input: err:", err)
		os.Exit(1)
	}
	outputBuf := bufio.NewWriter(os.Stdout)
	defer outputBuf.Flush()

	program, err := parser.Tokenise(sourceBuf)
	if err != nil {
		fmt.Println("error compiling program: err:", err)
		os.Exit(1)
	}
	if *dump {
		parser.Dump(program, outputBuf)
	} else if *print {
		parser.Print(program, outputBuf)
	} else if *eight {
		data := make([]byte, parser.DataSize)
		parser.Execute(data, program, inputBuf, outputBuf)
	} else {
		data := make([]int, parser.DataSize)
		parser.Execute(data, program, inputBuf, outputBuf)
	}
}

func inputReader(pathStr string, useDefault bool) (buff io.ByteReader, err error) {
	if useDefault && pathStr == "" {
		pathStr = "-"
	}
	if pathStr == "" {
		return nil, errors.New("no path provided")
	} else if pathStr == "-" {
		buff = bufio.NewReader(os.Stdin)
	} else {
		sourceFile, err := os.Open(pathStr)
		if err != nil {
			return nil, err
		}
		buff = bufio.NewReader(sourceFile)
	}
	return
}
