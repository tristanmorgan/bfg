# Development of Optimisations

## Initial project goals

I started this project as a learning exercise in Go and to make a portable BF interpreter I could more easily compile for the systems I use. I looked on GitHub for existing implementations and some came close to what I needed but in the end I took the option to DIY.

## Initial Parser

The initial parser takes a stream of characters and converts them to an array of instructions. One early choice was to combine addition and subtraction instructions into a single instruction with a value that could be positive or negative. Similarly, the data-pointer instructions could use a positive or negative offset. I knew I would expand on that later. The other early decision was to store the target program counter with jump instructions. This would save searching through instruction for a balanced jump when enacted.

## Simple Zero

A common code snippet used in BF is setting the current cell to zero, `[-]`. There are variation of this, one looping by subtracting one from the cell until zero and the other adding. Both potentially wrapping the integer value past the maximum/minimum value to eventually come back to zero. I decided both are functionally equivalent so I treat them both as "just set the value to zero". 

## Combining Instructions

As alluded earlier, the next optimisation was to combine repeated instructions and this also had the side effect of complementary instructions cancelling each other out. repeated pointer increments or decrements could them be treated as one "add 20 to the pointer" or a value increment followed by a decrement would result in no change to the data value.

## Complex Instructions

The next step for finding areas of improvement was to take existing BF code and analyse it for common patterns. Out of that was a new instruction internally called "move", it takes the data from one cell and moves it by an offset. This has two variations where a decrement instruction was at the start of the loop and alternatively at the end (`[->>+<<]` vs `[>>+<<-]`). Again, functionally the result is the same. Later another common patter was discovered was what I called a skip instruction where the code looks through a series of values to stop on the first zeroed one `[>]`.

## Performance Analysis

running benchmarks on different version of the code I discovered areas of improvement around the I/O buffering and realised the zero instruction could be extended to be a blind set instruction. It didn't need to look at the existing data value before setting it and this allowed combining increment and decrement instructions that followed a know zero point (such as exiting a loop) could also be treated as a "just set the value" stye of instruction.

## One more instruction

The last common pattern was identified as a fixed multiplication instruction. Where similar to the move instruction, for each decrement of the source cell another was incremented multiple times by a fixed amount `[->>++++<<]`. this too had an alternative where the decrement was at either the start or end `[>>++++<<-]`. This instruction required a second operand so I leveraged a following "noop" instruction to store it. 

## Dead code removal

Probably without a real-world performance improvement is dropping dead code from the output of the instruction parsing stage. Often a loop block is used to store comments in BF knowing that any characters within will never be interpreted during runtime. Identifying these comment blocks isn't always trivial but if a block starts directly after a cell being zeroed out then they could be skipped. To validate all the optimisations I kept a branch called no_optimisation updated to test a performance comparison and the correct operation of the optimised version.