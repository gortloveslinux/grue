package main

import (
	"fmt"
	"os"
	"path"
)

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

/* HEADER
----+---+----+----+----+--------------------------------------------------------------------------------------+
Hex |V  |Dyn |Int |Rst |Contents
----+---+----+----+----+--------------------------------------------------------------------------------------+
0   |1  |    |    |    |Version number (1 to 6)
----+---+----+----+----+--------------------------------------------------------------------------------------+
1   |3  |    |    |    |Flags 1 (in Versions 1 to 3):
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |    |    |Bit 1: Status line type: 0=score/turns, 1=hours:mins
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |    |    |2: Story file split across two discs?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |*   |*   |4: Status line not available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |*   |*   |5: Screen-splitting available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |*   |*   |6: Is a variable-pitch font the default?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |4  |    |    |    |Flags 1 (from Version 4):
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |*   |*   |Bit 0: Colours available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |*   |*   |1: Picture displaying available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |4  |    |*   |*   |2: Boldface available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |4  |    |*   |*   |3: Italic available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |4  |    |*   |*   |4: Fixed-space style available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |*   |*   |5: Sound effects available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |4  |    |*   |*   |7: Timed keyboard input available?
----+---+----+----+----+--------------------------------------------------------------------------------------+
4   |1  |    |    |    |Base of high memory (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
6   |1  |    |    |    |Initial value of program counter (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |    |    |Packed address of initial "main" routine
----+---+----+----+----+--------------------------------------------------------------------------------------+
8   |1  |    |    |    |Location of dictionary (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
A   |1  |    |    |    |Location of object table (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
C   |1  |    |    |    |Location of global variables table (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
E   |1  |    |    |    |Base of static memory (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
10  |1  |    |    |    |Flags 2:
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |1  | *  |*   |*   |Bit 0: Set when transcripting is on
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |3  | *  |    |*   |1: Game sets to force printing in fixed-pitch font
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  | *  |*   |    |2: Int sets to request screen redraw: game clears when it complies with this.
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |*   |*   |3: If set, game wants to use pictures
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |*   |*   |4: If set, game wants to use the UNDO opcodes
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |*   |*   |5: If set, game wants to use a mouse
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |    |    |6: If set, game wants to use colours
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |5  |    |*   |*   |7: If set, game wants to use sound effects
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |*   |*   |8: If set, game wants to use menus
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |   |    |    |    |(For bits 3,4,5,7 and 8, Int clears again if it cannot provide the requested effect.)
----+---+----+----+----+--------------------------------------------------------------------------------------+
18  |2  |    |    |    |Location of abbreviations table (byte address)
----+---+----+----+----+--------------------------------------------------------------------------------------+
1A  |3+ |    |    |    |Length of file (see note)
----+---+----+----+----+--------------------------------------------------------------------------------------+
1C  |3+ |    |    |    |Checksum of file
----+---+----+----+----+--------------------------------------------------------------------------------------+
1E  |4  |    |*   |*   |Interpreter number
----+---+----+----+----+--------------------------------------------------------------------------------------+
1F  |4  |    |*   |*   |Interpreter version
----+---+----+----+----+--------------------------------------------------------------------------------------+
Hex |V  |Dyn |Int |Rst |Contents
----+---+----+----+----+--------------------------------------------------------------------------------------+
20  |4  |    |*   |*   |Screen height (lines): 255 means "infinite"
----+---+----+----+----+--------------------------------------------------------------------------------------+
21  |4  |    |*   |*   |Screen width (characters)
----+---+----+----+----+--------------------------------------------------------------------------------------+
22  |5  |    |*   |*   |Screen width in units
----+---+----+----+----+--------------------------------------------------------------------------------------+
24  |5  |    |*   |*   |Screen height in units
----+---+----+----+----+--------------------------------------------------------------------------------------+
26  |5  |    |*   |*   |Font width in units (defined as width of a '0')
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |*   |*   |Font height in units
----+---+----+----+----+--------------------------------------------------------------------------------------+
27  |5  |    |*   |*   |Font height in units
----+---+----+----+----+--------------------------------------------------------------------------------------+
    |6  |    |*   |*   |Font width in units (defined as width of a '0')
----+---+----+----+----+--------------------------------------------------------------------------------------+
28  |6  |    |    |    |Routines offset (divided by 8)
----+---+----+----+----+--------------------------------------------------------------------------------------+
2A  |6  |    |    |    |Static strings offset (divided by 8)
----+---+----+----+----+--------------------------------------------------------------------------------------+
2C  |5  |    |*   |*   |Default background colour
----+---+----+----+----+--------------------------------------------------------------------------------------+
2D  |5  |    |*   |*   |Default foreground colour
----+---+----+----+----+--------------------------------------------------------------------------------------+
2E  |5  |    |    |    |Address of terminating characters table (bytes)
----+---+----+----+----+--------------------------------------------------------------------------------------+
30  |6  |    |*   |    |Total width in pixels of text sent to output stream 3
----+---+----+----+----+--------------------------------------------------------------------------------------+
32  |1  |    |*   |*   |Standard revision number
----+---+----+----+----+--------------------------------------------------------------------------------------+
34  |5  |    |    |    |Alphabet table address (bytes), or 0 for default
----+---+----+----+----+--------------------------------------------------------------------------------------+
36  |5  |    |    |    |Header extension table address (bytes)
----+---+----+----+----+--------------------------------------------------------------------------------------+
*/

type Header struct {
	data []byte
}

type Game struct {
	header *Header
	data   []byte
	size   int
}

func newGame(fileName string) (*Game, error) {
	gf, err := os.Open(path.Clean(fileName))
	if err != nil {
		return nil, fmt.Errorf("Couldn't load game ", err)
	}

	g := &Game{}
	g.size, err = gf.Read(g.data)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load game ", err)
	}

	g.header, err = newHeader(g.data)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load game ", err)
	}

	return g, nil
}

func newHeader(d []byte) (*Header, error) {
	h := &Header{data: d[0:296]}
	return h, nil
}

func (h *Header) getVersion() byte {
	return h.data[0:1][0]
}
