package instructionData

import enumInstructionModes "misc/nintasm/constants/enums/instructionModes"

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
type instructionOpcodesAndSupportedModes struct {
	BaseOpcode     uint8
	SupportedModes *[]instModes
	ModeOpcodes    map[instModes]uint8
}

func newInstOpcodeSet(baseOpcode uint8, supportedModes *[]instModes) instructionOpcodesAndSupportedModes {
	return instructionOpcodesAndSupportedModes{
		BaseOpcode:     baseOpcode,
		SupportedModes: supportedModes,
		ModeOpcodes:    make(map[instModes]uint8),
	}
}

// ======================================
// Commands that don't take an operand
var impliedInstructions = map[string]instructionOpcodesAndSupportedModes{
	"BRK": newInstOpcodeSet(0x00, &impliedModes),
	"CLC": newInstOpcodeSet(0x18, &impliedModes),
	"CLD": newInstOpcodeSet(0xd8, &impliedModes),
	"CLI": newInstOpcodeSet(0x58, &impliedModes),
	"CLV": newInstOpcodeSet(0xb8, &impliedModes),
	"DEX": newInstOpcodeSet(0xca, &impliedModes),
	"DEY": newInstOpcodeSet(0x88, &impliedModes),
	"INX": newInstOpcodeSet(0xe8, &impliedModes),
	"INY": newInstOpcodeSet(0xc8, &impliedModes),
	"NOP": newInstOpcodeSet(0xea, &impliedModes),
	"PHA": newInstOpcodeSet(0x48, &impliedModes),
	"PHP": newInstOpcodeSet(0x08, &impliedModes),
	"PLA": newInstOpcodeSet(0x68, &impliedModes),
	"PLP": newInstOpcodeSet(0x28, &impliedModes),
	"RTI": newInstOpcodeSet(0x40, &impliedModes),
	"RTS": newInstOpcodeSet(0x60, &impliedModes),
	"SEC": newInstOpcodeSet(0x38, &impliedModes),
	"SED": newInstOpcodeSet(0xf8, &impliedModes),
	"SEI": newInstOpcodeSet(0x78, &impliedModes),
	"TAX": newInstOpcodeSet(0xaa, &impliedModes),
	"TAY": newInstOpcodeSet(0xa8, &impliedModes),
	"TSX": newInstOpcodeSet(0xba, &impliedModes),
	"TXA": newInstOpcodeSet(0x8a, &impliedModes),
	"TXS": newInstOpcodeSet(0x9a, &impliedModes),
	"TYA": newInstOpcodeSet(0x98, &impliedModes),
}

// ======================================
// Branch instructions
var branchInstructions = map[string]instructionOpcodesAndSupportedModes{
	"BPL": newInstOpcodeSet(0x10, &branchModes),
	"BMI": newInstOpcodeSet(0x30, &branchModes),
	"BVC": newInstOpcodeSet(0x50, &branchModes),
	"BVS": newInstOpcodeSet(0x70, &branchModes),
	"BCC": newInstOpcodeSet(0x90, &branchModes),
	"BCS": newInstOpcodeSet(0xb0, &branchModes),
	"BNE": newInstOpcodeSet(0xd0, &branchModes),
	"BEQ": newInstOpcodeSet(0xf0, &branchModes),
}

// ======================================
// CPX and CPY
var compareXYInstructions = map[string]instructionOpcodesAndSupportedModes{
	"CPX": newInstOpcodeSet(0xe4, &compareXYModes),
	"CPY": newInstOpcodeSet(0xc4, &compareXYModes),
}

// ======================================
// INC and DEC
var incDecInstructions = map[string]instructionOpcodesAndSupportedModes{
	"DEC": newInstOpcodeSet(0xc6, &incDecModes),
	"INC": newInstOpcodeSet(0xe6, &incDecModes),
}

// ======================================
// Shift instructions
var shiftInstructions = map[string]instructionOpcodesAndSupportedModes{
	"ASL": newInstOpcodeSet(0x06, &shiftModes),
	"LSR": newInstOpcodeSet(0x46, &shiftModes),
	"ROL": newInstOpcodeSet(0x26, &shiftModes),
	"ROR": newInstOpcodeSet(0x66, &shiftModes),
}

// ======================================
var modifyAccumulatorInstructions = map[string]instructionOpcodesAndSupportedModes{
	"ADC": newInstOpcodeSet(0x65, &accumModModes),
	"AND": newInstOpcodeSet(0x25, &accumModModes),
	"CMP": newInstOpcodeSet(0xc5, &accumModModes),
	"EOR": newInstOpcodeSet(0x45, &accumModModes),
	"LDA": newInstOpcodeSet(0xa5, &accumModModes),
	"ORA": newInstOpcodeSet(0x05, &accumModModes),
	"SBC": newInstOpcodeSet(0xe5, &accumModModes),
}

// ======================================
var miscInstructions = map[string]instructionOpcodesAndSupportedModes{
	"BIT": newInstOpcodeSet(0x24, &bitModes),
	"JMP": newInstOpcodeSet(0x4c, &jmpModes),
	"JSR": newInstOpcodeSet(0x20, &jsrModes),
	"LDX": newInstOpcodeSet(0xa6, &ldxModes),
	"LDY": newInstOpcodeSet(0xa4, &ldyModes),
	"STA": newInstOpcodeSet(0x85, &staModes),
	"STX": newInstOpcodeSet(0x86, &stxModes),
	"STY": newInstOpcodeSet(0x84, &styModes),
}
