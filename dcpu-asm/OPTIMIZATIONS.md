## Optimizing processors

This is a listing of all pre- and post-processors which have a code 
optimization function, along with a description of what they do.

### -shorthand

This finds instances of "IFE B, A" and "IFN B, A" and checks if the first
operand is a short-form numeric literal. If so, the operands are swapped.

The DCPU spec states that short-form numbers can not be encoded in
the first operand, since its maximum value is not large enough to
hold all allowed literals. In this case, the assembler would have
to set B to 0x1f (next word) and store the value in a new word.

This needlessly increases the size of the program by one word.
For the IFE and IFN instructions, we can prevent this from 
happening by simply swapping the operands around.

For other instructions this is problematic, since swapping them out
changes the semantics of the operation. In those cases, we simply
allow the assembler to generate the extra word.

