## B.F.G.
BFG is an optimised [Brainfuck](https://esolangs.org/wiki/Brainfuck) interpreter written in Go.

Uses signed 16bit ints for data, memory wraps around at 65535, EOF returns -1.

## Optimisations

[-] is replaced with a ZERO command

repeat ++++ or --- are replaced with a single addition/subtraction

repeat >>> or <<< are replaced with a single pointer jump

[>>+<<-] and [->>+<<] merged into a move instruction.

buffered output prints on newline, 200 chars or input.
