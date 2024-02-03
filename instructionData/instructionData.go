package instructionData

import enumInstructionModes "misc/nintasm/enums/instructionModes"

type instModes = enumInstructionModes.Def

type modeDetailsKeys struct {
	abbrev      string
	description string
}

var InstructionModeDetails = map[enumInstructionModes.Def]modeDetailsKeys{
	enumInstructionModes.None:  {abbrev: "", description: ""},
	enumInstructionModes.A:     {abbrev: "A", description: "Accumulator"},
	enumInstructionModes.ABS:   {abbrev: "ABS", description: "Absolute"},
	enumInstructionModes.ABS_X: {abbrev: "ABS_X", description: "Absolute,X"},
	enumInstructionModes.ABS_Y: {abbrev: "ABS_Y", description: "Absolute,Y"},
	enumInstructionModes.IMM:   {abbrev: "IMM", description: "Immediate"},
	enumInstructionModes.IMPL:  {abbrev: "IMPL", description: "Implied"},
	enumInstructionModes.IND:   {abbrev: "IND", description: "Indirect"},
	enumInstructionModes.IND_X: {abbrev: "IND_X", description: "Indirect,X"},
	enumInstructionModes.IND_Y: {abbrev: "IND_Y", description: "Indirect,Y"},
	enumInstructionModes.REL:   {abbrev: "REL", description: "Releative"},
	enumInstructionModes.ZP:    {abbrev: "ZP", description: "Zero Page"},
	enumInstructionModes.ZP_X:  {abbrev: "ZP_X", description: "Zero Page,X"},
	enumInstructionModes.ZP_Y:  {abbrev: "ZP_Y", description: "Zero Page,Y"},
}

type xyConvert struct {
	X enumInstructionModes.Def
	Y enumInstructionModes.Def
}

// The X/Y indexing modes of supported modes
var ModesWithXYIndexes = map[enumInstructionModes.Def]xyConvert{
	enumInstructionModes.ABS: {X: enumInstructionModes.ABS_X, Y: enumInstructionModes.ABS_Y},
	enumInstructionModes.IND: {X: enumInstructionModes.IND_X, Y: enumInstructionModes.IND_Y},
	enumInstructionModes.ZP:  {X: enumInstructionModes.ZP_X, Y: enumInstructionModes.ZP_Y},
}

// The zero page equivalent of an absolute mode
var ABStoZP = map[instModes]instModes{
	enumInstructionModes.ABS:   enumInstructionModes.ZP,
	enumInstructionModes.ABS_X: enumInstructionModes.ZP_X,
	enumInstructionModes.ABS_Y: enumInstructionModes.ZP_Y,
}

var impliedModes = []instModes{enumInstructionModes.IMPL}
var branchModes = []instModes{enumInstructionModes.REL, enumInstructionModes.IMM}
var accumModModes = []instModes{enumInstructionModes.IMM, enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS, enumInstructionModes.ABS_X, enumInstructionModes.ABS_Y, enumInstructionModes.IND_X, enumInstructionModes.IND_Y}
var shiftModes = []instModes{enumInstructionModes.A, enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS, enumInstructionModes.ABS_X}
var compareXYModes = []instModes{enumInstructionModes.IMM, enumInstructionModes.ZP, enumInstructionModes.ABS}
var incDecModes = []instModes{enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS, enumInstructionModes.ABS_X}
var bitModes = []instModes{enumInstructionModes.ZP, enumInstructionModes.ABS}
var jmpModes = []instModes{enumInstructionModes.ABS, enumInstructionModes.IND}
var jsrModes = []instModes{enumInstructionModes.ABS}
var ldxModes = []instModes{enumInstructionModes.IMM, enumInstructionModes.ZP, enumInstructionModes.ZP_Y, enumInstructionModes.ABS, enumInstructionModes.ABS_Y}
var ldyModes = []instModes{enumInstructionModes.IMM, enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS, enumInstructionModes.ABS_X}
var staModes = []instModes{enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS, enumInstructionModes.ABS_X, enumInstructionModes.ABS_Y, enumInstructionModes.IND_X, enumInstructionModes.IND_Y}
var stxModes = []instModes{enumInstructionModes.ZP, enumInstructionModes.ZP_Y, enumInstructionModes.ABS}
var styModes = []instModes{enumInstructionModes.ZP, enumInstructionModes.ZP_X, enumInstructionModes.ABS}

// ======================================
type instructionOpcodeAndSupportedModes struct {
	BaseOpcode uint8
	Modes      *[]instModes
}

// ======================================
// Commands that don't take an operand
var impliedInstructions = map[string]instructionOpcodeAndSupportedModes{
	"BRK": {BaseOpcode: 0x00, Modes: &impliedModes},
	"CLC": {BaseOpcode: 0x18, Modes: &impliedModes},
	"CLD": {BaseOpcode: 0xd8, Modes: &impliedModes},
	"CLI": {BaseOpcode: 0x58, Modes: &impliedModes},
	"CLV": {BaseOpcode: 0xb8, Modes: &impliedModes},
	"DEX": {BaseOpcode: 0xca, Modes: &impliedModes},
	"DEY": {BaseOpcode: 0x88, Modes: &impliedModes},
	"INX": {BaseOpcode: 0xe8, Modes: &impliedModes},
	"INY": {BaseOpcode: 0xc8, Modes: &impliedModes},
	"NOP": {BaseOpcode: 0xea, Modes: &impliedModes},
	"PHA": {BaseOpcode: 0x48, Modes: &impliedModes},
	"PHP": {BaseOpcode: 0x08, Modes: &impliedModes},
	"PLA": {BaseOpcode: 0x68, Modes: &impliedModes},
	"PLP": {BaseOpcode: 0x28, Modes: &impliedModes},
	"RTI": {BaseOpcode: 0x40, Modes: &impliedModes},
	"RTS": {BaseOpcode: 0x60, Modes: &impliedModes},
	"SEC": {BaseOpcode: 0x38, Modes: &impliedModes},
	"SED": {BaseOpcode: 0xf8, Modes: &impliedModes},
	"SEI": {BaseOpcode: 0x78, Modes: &impliedModes},
	"TAX": {BaseOpcode: 0xaa, Modes: &impliedModes},
	"TAY": {BaseOpcode: 0xa8, Modes: &impliedModes},
	"TSX": {BaseOpcode: 0xba, Modes: &impliedModes},
	"TXA": {BaseOpcode: 0x8a, Modes: &impliedModes},
	"TXS": {BaseOpcode: 0x9a, Modes: &impliedModes},
	"TYA": {BaseOpcode: 0x98, Modes: &impliedModes},
}

// ======================================
// Branch instructions
var branchInstructions = map[string]instructionOpcodeAndSupportedModes{
	"BPL": {BaseOpcode: 0x10, Modes: &branchModes},
	"BMI": {BaseOpcode: 0x30, Modes: &branchModes},
	"BVC": {BaseOpcode: 0x50, Modes: &branchModes},
	"BVS": {BaseOpcode: 0x70, Modes: &branchModes},
	"BCC": {BaseOpcode: 0x90, Modes: &branchModes},
	"BCS": {BaseOpcode: 0xb0, Modes: &branchModes},
	"BNE": {BaseOpcode: 0xd0, Modes: &branchModes},
	"BEQ": {BaseOpcode: 0xf0, Modes: &branchModes},
}

// ======================================
// CPX and CPY
var compareXYInstructions = map[string]instructionOpcodeAndSupportedModes{
	"CPX": {BaseOpcode: 0xe4, Modes: &compareXYModes},
	"CPY": {BaseOpcode: 0xc4, Modes: &compareXYModes},
}

// ======================================
// INC and DEC
var incDecInstructions = map[string]instructionOpcodeAndSupportedModes{
	"DEC": {BaseOpcode: 0xc6, Modes: &incDecModes},
	"INC": {BaseOpcode: 0xe6, Modes: &incDecModes},
}

// ======================================
// Shift instructions
var shiftInstructions = map[string]instructionOpcodeAndSupportedModes{
	"ASL": {BaseOpcode: 0x06, Modes: &shiftModes},
	"LSR": {BaseOpcode: 0x46, Modes: &shiftModes},
	"ROL": {BaseOpcode: 0x26, Modes: &shiftModes},
	"ROR": {BaseOpcode: 0x66, Modes: &shiftModes},
}

// ======================================
var modifyAccumulatorInstructions = map[string]instructionOpcodeAndSupportedModes{
	"ADC": {BaseOpcode: 0x65, Modes: &accumModModes},
	"AND": {BaseOpcode: 0x25, Modes: &accumModModes},
	"CMP": {BaseOpcode: 0xc5, Modes: &accumModModes},
	"EOR": {BaseOpcode: 0x45, Modes: &accumModModes},
	"LDA": {BaseOpcode: 0xa5, Modes: &accumModModes},
	"ORA": {BaseOpcode: 0x05, Modes: &accumModModes},
	"SBC": {BaseOpcode: 0xe5, Modes: &accumModModes},
}

// ======================================
var miscInstructions = map[string]instructionOpcodeAndSupportedModes{
	"BIT": {BaseOpcode: 0x24, Modes: &bitModes},
	"JMP": {BaseOpcode: 0x4c, Modes: &jmpModes},
	"JSR": {BaseOpcode: 0x20, Modes: &jsrModes},
	"LDX": {BaseOpcode: 0xa6, Modes: &ldxModes},
	"LDY": {BaseOpcode: 0xa4, Modes: &ldyModes},
	"STA": {BaseOpcode: 0x85, Modes: &staModes},
	"STX": {BaseOpcode: 0x86, Modes: &stxModes},
	"STY": {BaseOpcode: 0x84, Modes: &styModes},
}

var OpcodesAndSupportedModes = map[string]instructionOpcodeAndSupportedModes{}

func init() {
	buildModes(&miscInstructions)
	buildModes(&modifyAccumulatorInstructions)
	buildModes(&shiftInstructions)
	buildModes(&impliedInstructions)
	buildModes(&branchInstructions)
	buildModes(&compareXYInstructions)
	buildModes(&incDecInstructions)
}

func buildModes(modeMap *map[string]instructionOpcodeAndSupportedModes) {
	for k, v := range *modeMap {
		OpcodesAndSupportedModes[k] = v
	}
}
