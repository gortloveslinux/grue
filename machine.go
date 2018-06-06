package main

import "fmt"

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

const (
	operand_mask     = 0x3
	operand_omitted  = 0x3
	operand_variable = 0x2
	operand_small    = 0x1
	operand_large    = 0x0
)

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

func (zm *ZMachine) readShort() uint16 {
	zm.ip = zm.ip + 2
	return getUint16(zm.buf, zm.ip-2)
}

func (zm *ZMachine) InterpretInstruction() error {
	op := zm.peekByte()
	context := fmt.Sprintf("ip: %x, opcode: %x\n", zm.ip, op)

	if op == 0xbe && zm.header.version >= 5 {
		zm.logger("Extended form instructions not implemented yet")
		return nil
	} else if op>>6 == 0x3 {
		zm.logger("Interpreting VAR form instruction: %s", context)
		return zm.InterpretVarInstruction()
	} else if op>>6 == 0x2 {
		zm.logger("Short form instructions not implemented yet")
		return nil
	}
	zm.logger("Long form instructions not implemented yet")
	return nil
}

// In variable form, if bit 5 is 0 then the count is 2OP; if it is 1, then the count is VAR.
// The opcode number is given in the bottom 5 bits.
func (zm *ZMachine) InterpretVarInstruction() error {
	opcode := zm.readByte()
	opcount := opcode >> 5 & 0x1
	instruction := opcode & 0x1f
	zm.logger("instruction: %x, operand count: %x\n", instruction, opcount)

	if opcount == 0 {
		// 2OP instructions
		return nil
	}
	// VAR OP instruction
	switch instruction {
	case 0x0:
		return zm.zCall()
	default:
		return fmt.Errorf("Unknown instruction: %x, ip: %x", opcode, zm.ip)
	}
}

func (zm *ZMachine) zCall() error {
	args, _, err := zm.getOperands()
	if err != nil {
		return fmt.Errorf("error interpreting call instruction. ip: %x. %v", zm.ip, err)
	}
	addr, err := zm.getPackedAddress(args[0], true)
	if err != nil {
		return fmt.Errorf("error interpreting call instruction. ip: %x. %v", zm.ip, err)
	}

	zm.logger("Call %x", addr)
	return nil
}

// In variable or extended forms, a byte of 4 operand types is given next.
// This contains 4 2-bit fields: bits 6 and 7 are the first field, bits 0 and 1 the fourth.
// The values are operand types as above. Once one type has been given as 'omitted', all subsequent ones must be.
// Example: $$00101111 means large constant followed by variable (and no third or fourth opcode).
func (zm *ZMachine) getOperands() ([]uint16, uint16, error) {
	args := make([]uint16, 4)
	argc := uint16(0)
	opTypes := zm.readByte()
	for i := uint8(0); i < 4; i++ {
		optyp := (opTypes >> (6 - i)) & 0x3
		if optyp == operand_omitted {
			break
		}
		switch optyp {
		case operand_variable:
			// to implement
		case operand_small:
			// to implement
		case operand_large:
			args[i] = zm.readShort()
		default:
			return nil, 0, fmt.Errorf("unknown operand type: %x, operand number: %x, ip: %x", optyp, i, zm.ip-1)
		}
	}

	return args, argc, nil
}

// A packed address specifies where a routine or string begins in high memory.
// Given a packed address P, the formula to obtain the corresponding byte address B is:
//
//  2P           Versions 1, 2 and 3
//  4P           Versions 4 and 5
//  4P + 8R_O    Versions 6 and 7, for routine calls
//  4P + 8S_O    Versions 6 and 7, for print_paddr
//  8P           Version 8
//R_O and S_O are the routine and strings offsets (specified in the header as words at $28 and $2a, respectively).
func (zm *ZMachine) getPackedAddress(addr uint16, routine bool) (uint32, error) {
	switch zm.header.version {
	case 1, 2, 3:
		return uint32(addr) * 2, nil
	case 4, 5:
		return uint32(addr) * 4, nil
	case 6, 7:
		offset := zm.header.stringOffset
		if routine {
			offset = zm.header.routineOffset
		}
		return uint32(addr)*4 + 8*uint32(offset), nil
	case 8:
		return uint32(addr) * 8, nil
	default:
		return 0, fmt.Errorf("unsupported version for packed address. version: %x", zm.header.version)
	}
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
