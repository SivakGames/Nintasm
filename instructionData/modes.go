package instructionData

type InstructionModes int

const (
	None InstructionModes = iota + 0
	A
	ABS
	ABS_X
	ABS_Y
	IMM
	IMPL
	IND
	IND_X
	IND_Y
	REL
	ZP
	ZP_X
	ZP_Y
)

type ModeInfo struct {
	abbrev      string
	description string
}

var InstructionModeDetails = map[InstructionModes]ModeInfo{
	None:  {abbrev: "", description: ""},
	A:     {abbrev: "A", description: "Accumulator"},
	ABS:   {abbrev: "ABS", description: "Absolute"},
	ABS_X: {abbrev: "ABS_X", description: "Absolute,X"},
	ABS_Y: {abbrev: "ABS_Y", description: "Absolute,Y"},
	IMM:   {abbrev: "ABS_Y", description: "Immediate"},
	IMPL:  {abbrev: "IMPL", description: "Implied"},
	IND:   {abbrev: "IND", description: "Indirect"},
	IND_X: {abbrev: "IND_X", description: "Indirect,X"},
	IND_Y: {abbrev: "IND_Y", description: "Indirect,Y"},
	REL:   {abbrev: "REL", description: "Releative"},
	ZP:    {abbrev: "ZP", description: "Zero Page"},
	ZP_X:  {abbrev: "ZP_X", description: "Zero Page,X"},
	ZP_Y:  {abbrev: "ZP_Y", description: "Zero Page,Y"},
}

var impliedModes = []InstructionModes{IMPL}
var branchModes = []InstructionModes{REL, IMM}
var accumModModes = []InstructionModes{IMM, ZP, ZP_X, ABS, ABS_X, ABS_Y, IND_X, IND_Y}
var shiftModes = []InstructionModes{A, ZP, ZP_X, ABS, ABS_X}
var compareXYModes = []InstructionModes{IMM, ZP, ABS}
var incDecModes = []InstructionModes{ZP, ZP_X, ABS, ABS_X}
var bitModes = []InstructionModes{ZP, ABS}
var jmpModes = []InstructionModes{ABS, IND}
var jsrModes = []InstructionModes{ABS}
var ldxModes = []InstructionModes{IMM, ZP, ZP_Y, ABS, ABS_Y}
var ldyModes = []InstructionModes{IMM, ZP, ZP_X, ABS, ABS_X}
var staModes = []InstructionModes{ZP, ZP_X, ABS, ABS_X, ABS_Y, IND_X, IND_Y}
var stxModes = []InstructionModes{ZP, ZP_Y, ABS}
var styModes = []InstructionModes{ZP, ZP_X, ABS}

// ======================================
type individuallyAllowed struct {
	Base  uint8
	Modes []InstructionModes
}

// ======================================
// Commands that don't take an operand
var impliedInstructions = map[string]individuallyAllowed{
	"BRK": {Base: 0x00, Modes: impliedModes},
	"CLC": {Base: 0x18, Modes: impliedModes},
	"CLD": {Base: 0xd8, Modes: impliedModes},
	"CLI": {Base: 0x58, Modes: impliedModes},
	"CLV": {Base: 0xb8, Modes: impliedModes},
	"DEX": {Base: 0xca, Modes: impliedModes},
	"DEY": {Base: 0x88, Modes: impliedModes},
	"INX": {Base: 0xe8, Modes: impliedModes},
	"INY": {Base: 0xc8, Modes: impliedModes},
	"NOP": {Base: 0xea, Modes: impliedModes},
	"PHA": {Base: 0x48, Modes: impliedModes},
	"PHP": {Base: 0x08, Modes: impliedModes},
	"PLA": {Base: 0x68, Modes: impliedModes},
	"PLP": {Base: 0x28, Modes: impliedModes},
	"RTI": {Base: 0x40, Modes: impliedModes},
	"RTS": {Base: 0x60, Modes: impliedModes},
	"SEC": {Base: 0x38, Modes: impliedModes},
	"SED": {Base: 0xf8, Modes: impliedModes},
	"SEI": {Base: 0x78, Modes: impliedModes},
	"TAX": {Base: 0xaa, Modes: impliedModes},
	"TAY": {Base: 0xa8, Modes: impliedModes},
	"TSX": {Base: 0xba, Modes: impliedModes},
	"TXA": {Base: 0x8a, Modes: impliedModes},
	"TXS": {Base: 0x9a, Modes: impliedModes},
	"TYA": {Base: 0x98, Modes: impliedModes},
}

// ======================================
// Branch instructions
var branchInstructions = map[string]individuallyAllowed{
	"BPL": {Base: 0x10, Modes: branchModes},
	"BMI": {Base: 0x30, Modes: branchModes},
	"BVC": {Base: 0x50, Modes: branchModes},
	"BVS": {Base: 0x70, Modes: branchModes},
	"BCC": {Base: 0x90, Modes: branchModes},
	"BCS": {Base: 0xb0, Modes: branchModes},
	"BNE": {Base: 0xd0, Modes: branchModes},
	"BEQ": {Base: 0xf0, Modes: branchModes},
}

// ======================================
// CPX and CPY
var compareXYInstructions = map[string]individuallyAllowed{
	"CPX": {Base: 0xe4, Modes: compareXYModes},
	"CPY": {Base: 0xc4, Modes: compareXYModes},
}

// ======================================
// INC and DEC
var incDecInstructions = map[string]individuallyAllowed{
	"DEC": {Base: 0xc6, Modes: incDecModes},
	"INC": {Base: 0xe6, Modes: incDecModes},
}

// ======================================
// Shift instructions
var shiftInstructions = map[string]individuallyAllowed{
	"ASL": {Base: 0x06, Modes: shiftModes},
	"LSR": {Base: 0x46, Modes: shiftModes},
	"ROL": {Base: 0x26, Modes: shiftModes},
	"ROR": {Base: 0x66, Modes: shiftModes},
}

// ======================================
var modifyAccumulatorInstructions = map[string]individuallyAllowed{
	"ADC": {Base: 0x65, Modes: accumModModes},
	"AND": {Base: 0x25, Modes: accumModModes},
	"CMP": {Base: 0xc5, Modes: accumModModes},
	"EOR": {Base: 0x45, Modes: accumModModes},
	"LDA": {Base: 0xa5, Modes: accumModModes},
	"ORA": {Base: 0x05, Modes: accumModModes},
	"SBC": {Base: 0xe5, Modes: accumModModes},
}

// ======================================
var miscInstructions = map[string]individuallyAllowed{
	"BIT": {Base: 0x24, Modes: bitModes},
	"JMP": {Base: 0x4c, Modes: jmpModes},
	"JSR": {Base: 0x20, Modes: jsrModes},
	"LDX": {Base: 0xa6, Modes: ldxModes},
	"LDY": {Base: 0xa4, Modes: ldyModes},
	"STA": {Base: 0x85, Modes: staModes},
	"STX": {Base: 0x86, Modes: stxModes},
	"STY": {Base: 0x84, Modes: styModes},
}

var AllowedModes = map[string]individuallyAllowed{}

func init() {
	buildModes(&miscInstructions)
	buildModes(&modifyAccumulatorInstructions)
	buildModes(&shiftInstructions)
	buildModes(&impliedInstructions)
	buildModes(&branchInstructions)
	buildModes(&compareXYInstructions)
	buildModes(&incDecInstructions)
}

func buildModes(modeMap *map[string]individuallyAllowed) {
	for k, v := range *modeMap {
		AllowedModes[k] = v
	}
}
