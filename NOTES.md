# Approach
 1. Write a disassembler
   My thought here is that I can have all my data structures created and can move from machine code to assembly code.
   * Get Def. of opcodes
   * Get Def. of header
   * Get Def. of routines
   * Get Def. of game file

# Useful commands
* `paste (xxd -b -c 16 data/anchor.z8|psub) (xxd data/anchor.z8|psub) | head`
