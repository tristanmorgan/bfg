## B.F.G.

[![Go](https://github.com/tristanmorgan/bfg/actions/workflows/go-test-build.yml/badge.svg)](https://github.com/tristanmorgan/bfg/actions/workflows/go-test-build.yml)

BFG is an optimised [Brainfuck](https://esolangs.org/wiki/Brainfuck) interpreter written in Go.

Uses signed ints for data (platform specific 32/64), memory wraps around at 65535, EOF returns -1.

## Optimisations

[-] is replaced with a blind set 0 command

repeat ++++ or --- are replaced with a single addition/subtraction

addition/subtraction after zero does a blind set.

repeat >>> or <<< are replaced with a single pointer jump

[>>+<<-] and [->>+<<] merged into a move instruction.

buffered output prints on newline, 200 chars or input.
