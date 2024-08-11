package instructionData

import (
	enumInstructionModes "misc/nintasm/constants/enums/instructionModes"
)

type instModes = enumInstructionModes.Def

type modeDetailsKeys struct {
	Abbrev      string
	Description string
}

var InstructionModeEnumDetails = map[enumInstructionModes.Def]modeDetailsKeys{
	enumInstructionModes.None:  {Abbrev: "", Description: ""},
	enumInstructionModes.A:     {Abbrev: "A", Description: "Accumulator"},
	enumInstructionModes.ABS:   {Abbrev: "ABS", Description: "Absolute"},
	enumInstructionModes.ABS_X: {Abbrev: "ABS_X", Description: "Absolute,X"},
	enumInstructionModes.ABS_Y: {Abbrev: "ABS_Y", Description: "Absolute,Y"},
	enumInstructionModes.IMM:   {Abbrev: "IMM", Description: "Immediate"},
	enumInstructionModes.IMPL:  {Abbrev: "IMPL", Description: "Implied"},
	enumInstructionModes.IND:   {Abbrev: "IND", Description: "Indirect"},
	enumInstructionModes.IND_X: {Abbrev: "IND_X", Description: "Indirect,X"},
	enumInstructionModes.IND_Y: {Abbrev: "IND_Y", Description: "Indirect,Y"},
	enumInstructionModes.REL:   {Abbrev: "REL", Description: "Releative"},
	enumInstructionModes.ZP:    {Abbrev: "ZP", Description: "Zero Page"},
	enumInstructionModes.ZP_X:  {Abbrev: "ZP_X", Description: "Zero Page,X"},
	enumInstructionModes.ZP_Y:  {Abbrev: "ZP_Y", Description: "Zero Page,Y"},
}

var InstructionModeOperandRequiredBytes = map[enumInstructionModes.Def]int{
	enumInstructionModes.None:  -1,
	enumInstructionModes.A:     0,
	enumInstructionModes.ABS:   2,
	enumInstructionModes.ABS_X: 2,
	enumInstructionModes.ABS_Y: 2,
	enumInstructionModes.IMM:   1,
	enumInstructionModes.IMPL:  0,
	enumInstructionModes.IND:   2,
	enumInstructionModes.IND_X: 1,
	enumInstructionModes.IND_Y: 1,
	enumInstructionModes.REL:   1,
	enumInstructionModes.ZP:    1,
	enumInstructionModes.ZP_X:  1,
	enumInstructionModes.ZP_Y:  1,
}

type xy struct {
	X enumInstructionModes.Def
	Y enumInstructionModes.Def
}

// The X/Y indexing modes of supported modes
var InstructionModeEnumToXYModeEnum = map[enumInstructionModes.Def]xy{
	enumInstructionModes.ABS: {X: enumInstructionModes.ABS_X, Y: enumInstructionModes.ABS_Y},
	enumInstructionModes.IND: {X: enumInstructionModes.IND_X, Y: enumInstructionModes.IND_Y},
	enumInstructionModes.ZP:  {X: enumInstructionModes.ZP_X, Y: enumInstructionModes.ZP_Y},
}

// The zero page equivalent of an absolute mode
var InstructionABSEnumToZPEnum = map[instModes]instModes{
	enumInstructionModes.ABS:   enumInstructionModes.ZP,
	enumInstructionModes.ABS_X: enumInstructionModes.ZP_X,
	enumInstructionModes.ABS_Y: enumInstructionModes.ZP_Y,
}

var OpcodesAndSupportedModes = map[string]instructionOpcodesAndSupportedModes{}

func init() {
	assignModesAndOpcodes(&impliedInstructions)
	assignModesAndOpcodes(&branchInstructions)
	assignModesAndOpcodes(&compareXYInstructions)
	assignModesAndOpcodes(&incDecInstructions)
	assignModesAndOpcodes(&modifyAccumulatorInstructions)
	assignModesAndOpcodes(&shiftInstructions)

	assignModesAndOpcodes(&miscInstructions)
}

func assignModesAndOpcodes(modeMap *map[string]instructionOpcodesAndSupportedModes) {
	for instructionNameAsKey, opcodeModeSet := range *modeMap {
		baseOpcode := opcodeModeSet.BaseOpcode
		for _, mode := range *opcodeModeSet.SupportedModes {
			xyMod := instructionNameAsKey == "LDX" || instructionNameAsKey == "LDY" || instructionNameAsKey == "CPX" || instructionNameAsKey == "CPY"
			jmpMod := instructionNameAsKey == "JMP" || instructionNameAsKey == "JSR"
			adj := getAdjustedOpcode(mode, baseOpcode, xyMod, jmpMod)
			opcodeModeSet.ModeOpcodes[mode] = adj
		}
		OpcodesAndSupportedModes[instructionNameAsKey] = opcodeModeSet
	}
}

func getAdjustedOpcode(mode instModes, baseOpcode uint8, xyMod bool, jmpMod bool) uint8 {
	adjustedOpcode := baseOpcode
	switch mode {
	case enumInstructionModes.IMPL,
		enumInstructionModes.REL:
		break
	case enumInstructionModes.IMM:
		if !xyMod {
			adjustedOpcode += 0x04
		} else {
			adjustedOpcode -= 0x04
		}
	case enumInstructionModes.A:
		adjustedOpcode += 0x04
	case enumInstructionModes.ZP:
		break
	case enumInstructionModes.ZP_X, enumInstructionModes.ZP_Y:
		adjustedOpcode += 0x10
	case enumInstructionModes.ABS:
		if !jmpMod {
			adjustedOpcode += 0x08
		}
	case enumInstructionModes.ABS_X:
		adjustedOpcode += 0x18
	case enumInstructionModes.ABS_Y:
		if !xyMod {
			adjustedOpcode += 0x14
		} else {
			adjustedOpcode += 0x18
		}
	case enumInstructionModes.IND:
		adjustedOpcode += 0x20
	case enumInstructionModes.IND_X:
		adjustedOpcode -= 0x04
	case enumInstructionModes.IND_Y:
		adjustedOpcode += 0x0c
	default:
		panic("🛑 Something is terribly wrong with building instruction modes")
	}

	return adjustedOpcode
}
