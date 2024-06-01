## B.F.G.

* [![license MIT](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
* [![GoReportCard](https://goreportcard.com/badge/github.com/tristanmorgan/bfg)](https://goreportcard.com/report/github.com/tristanmorgan/bfg)
* [![Go](https://github.com/tristanmorgan/bfg/actions/workflows/go-test-build.yml/badge.svg)](https://github.com/tristanmorgan/bfg/actions/workflows/go-test-build.yml)

BFG is an optimised [Brainfuck](https://esolangs.org/wiki/Brainfuck) interpreter written in Go.

Uses signed ints for data (platform specific 32/64), memory wraps around at 65535, EOF returns -1.

Optionally uses 8 bit data cells.

Buffered output flushes on newline, 200 chars or input.

## Optimisations

Operates with a instruction parse then execute pattern.

 * loop start/end are calculated up front.
 * repeat ++++ or --- are replaced with a single addition/subtraction
 * [-] is converted to a blind set 0 command
 * addition/subtraction after zero does a blind set.
 * repeat >>> or <<< are replaced with a single pointer jump
 * [>>>] is merged into a skip instruction.
 * [>>+<<-] is merged into a move instruction.
 * [<+++++>-] is converted to a multiply instruction.
 * [<<+>+>-] is converted to a duplicate instruction.
 * and dead code removal.

for performance comparison see no_optimisation branch.

## Usage

    Usage:
      bfg [option] source.bf [input]
    
    Options:
      -e, --eight	eight bit execution
      -p, --print	pretty print parsed program
      -v, --version	display version

May use - as source to read program from STDIN and output is STDOUT

    #!/usr/bin/env bf
    +++++++++[>++++++++>++++++++++++>++++>++++++++++++>+++++++++++>+<<<<<<-]>---.>++
    .----.+++++.++++++++++.>----.<-----.>>----.+.<<-.>.<<---.>-.>>>--.<.+++++.<<<+++
    +.>+++.>>>++.<---.<+.>>>+.

