package instructionData

import (
	"fmt"
	enumInstructionModes "misc/nintasm/enums/instructionModes"
)

type instModes = enumInstructionModes.Def

type modeDetailsKeys struct {
	abbrev      string
	description string
}

var InstructionModeEnumDetails = map[enumInstructionModes.Def]modeDetailsKeys{
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
	for key, value := range *modeMap {
		baseOpcode := value.BaseOpcode
		for _, mode := range *value.SupportedModes {
			xyMod := key == "LDX" || key == "LDY"
			adj := getAdjustedOpcode(mode, baseOpcode, xyMod)
			value.ModeOpcodes[mode] = adj
		}
		fmt.Println(value)

		OpcodesAndSupportedModes[key] = value
	}
}

func getAdjustedOpcode(mode instModes, baseOpcode uint8, xyMod bool) uint8 {
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
		adjustedOpcode += 0x08
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
		panic("Something is terribly wrong with building instruction modes")
	}

	return adjustedOpcode
}
