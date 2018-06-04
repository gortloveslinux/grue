package main

/* OPCODES
A single Z-machine instruction consists of the following sections (and in the order shown):
  Opcode               1 or 2 bytes
  (Types of operands)  1 or 2 bytes: 4 or 8 2-bit fields
  Operands             Between 0 and 8 of these: each 1 or 2 bytes
  (Store variable)     1 byte
  (Branch offset)      1 or 2 bytes
  (Text to print)      An encoded string (of unlimited length)
*/

/* ROUTINES
5.1
A routine is required to begin at an address in memory which can be represented by a packed address (for instance, in Version 5 it must occur at a byte address which is divisible by 4).
5.2
A routine begins with one byte indicating the number of local variables it has (between 0 and 15 inclusive).
5.2.1
In Versions 1 to 4, that number of 2-byte words follows, giving initial values for these local variables. In Versions 5 and later, the initial values are all zero.
5.3
Execution of instructions begins from the byte after this header information. There is no formal 'end-marker' for a routine (it is simply assumed that execution eventually results in a return taking place).
5.4
In Version 6, there is a "main" routine (whose packed address is stored in the word at $06 in the header) called when the game starts up. It is illegal to return from this routine.
5.5
In all other Versions, the word at $06 contains the byte address of the first instruction to execute. The Z-machine starts in an environment with no local variables from which, again, a return is illegal.
------
Remarks
Note that it is permissible for a routine to be in dynamic memory. Marnix Klooster suggests this might be used for compiling code at run time!
In Versions 3 and 4, Inform always stores 0 as the initial values for local variables.
*/

type ZMachine struct {
	header ZHeader
	ip     uint32
	buf    []uint8
	logger func(string, ...interface{})
}

func getUint16(buf []byte, offset uint32) uint16 {
	return (uint16(buf[offset]) << 8) | (uint16)(buf[offset+1])
}

func getUint32(buf []byte, offset uint32) uint32 {
	return (uint32(buf[offset]) << 24) | (uint32(buf[offset+1]) << 16) | (uint32(buf[offset+2]) << 8) | uint32(buf[offset+3])
}

func (zm *ZMachine) Init(buf []byte, logger func(string, ...interface{})) error {
	var header ZHeader
	header.load(buf)

	zm.buf = buf
	zm.header = header
	zm.ip = uint32(header.ip)
	zm.logger = logger

	return nil
}

func (zm *ZMachine) peekByte() uint8 {
	return zm.buf[zm.ip]
}

func (zm *ZMachine) readByte() uint8 {
	zm.ip++
	return zm.buf[zm.ip-1]
}

func (zm *ZMachine) InterpretInstruction() {
	op := zm.peekByte()
	zm.logger("ip: %x, op: %x\n", zm.ip, op)

	if op == 0xbe && zm.header.version >= 5 {
		zm.logger("Extended codes not implemented yet")
		// return error?
	} else if op>>6 == 0x3 {
		// var operation form
		zm.InterpretVarInstruction()
	} else if op>>6 == 0x2 {
		// short operation form
		zm.logger("short form instructions not implemented yet")
	}
	// long operation form
	zm.InterpretLongInstruction()
}

func (zm *ZMachine) InterpretLongInstruction() {

}

func (zm *ZMachine) InterpretVarInstruction() {
	//opcode := zm.readByte()
	//opcount :=
}

func (zm *ZMachine) verify() {
	scale := uint32(0)

	switch zm.header.version {
	case 3, 4, 5:
		scale = 2
	case 6, 7, 8:
		scale = 8
	}
	end := zm.header.fileLength * scale

	var sum = uint32(0)
	for _, v := range zm.buf[0x40:end] {
		sum = (sum + uint32(v)) % uint32(0x10000)
	}
	// TODO more
}
