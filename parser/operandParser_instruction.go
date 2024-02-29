package parser

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumInstructionModes "misc/nintasm/constants/enums/instructionModes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/constants/instructionData"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
	"strings"
)

type InstructionOperandParser struct {
	OperandParser
}

// Constructor
func NewInstructionOperandParser() InstructionOperandParser {
	return InstructionOperandParser{}
}

func (p *InstructionOperandParser) Process(operationValue string) error {
	var err error = nil
	var instructionMode instModeEnum
	var operand Node = operandFactory.EmptyNode()

	instructionName := strings.ToUpper(operationValue)

	isBranch := false
	instructionXYIndex := enumTokenTypes.None
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]
	operandList := []Node{}
	useInstructionMode, useInstructionZPMode := enumInstructionModes.None, enumInstructionModes.None

	switch p.lookaheadType {
	//Try no operand
	case enumTokenTypes.None:
		instructionMode = enumInstructionModes.IMPL
	//Try A (for shifts)
	case enumTokenTypes.REGISTER_A:
		instructionMode = enumInstructionModes.A
		err = p.eatFreelyAndAdvance(enumTokenTypes.REGISTER_A)
		if err != nil {
			return err // ❌ Fails
		}
		if p.lookaheadType != enumTokenTypes.None {
			return errorHandler.AddNew(enumErrorCodes.InstBadAccumMode) // ❌ Fails
		}

	default:
		isBranch = checkIfBranchInstruction(instructionName)
		operandList, err = p.GetOperandList(0, 1, false, []string{"instruction"})
		if err != nil {
			return err // ❌ Fails
		}
		instructionXYIndex = p.instructionXYIndex
		instructionMode = p.instructionMode
		if isBranch && instructionMode == enumInstructionModes.ABS {
			instructionMode = enumInstructionModes.REL
		}
	}

	// -----------------------------

	//If an index is present, see if it's usable with desired mode and set it
	if instructionXYIndex == enumTokenTypes.REGISTER_X || instructionXYIndex == enumTokenTypes.REGISTER_Y {
		instructionMode, err = checkModeSupportsXY(instructionMode, instructionXYIndex)
		if err != nil {
			return err // ❌ Fails
		}
	} else if instructionXYIndex != enumTokenTypes.None {
		panic("MAJOR ERROR with XY index checks!!!")
	}

	if len(operandList) == 1 {
		operand = operandList[0]
	}

	//Used for auto ZP convert if possible
	instructionZPModeEquivalent := getModeZeroPageEquivalent(instructionMode)

	// Check if instruction itself supports mode
	for _, mode := range *opcodesAndSupportedModes.SupportedModes {
		if mode == instructionMode {
			useInstructionMode = mode
			continue
		}
		if mode == instructionZPModeEquivalent && operand.Resolved {
			if operandFactory.ValidateNumericNodeIs8BitValue(&operand) {
				useInstructionZPMode = mode
			}
		}
	}

	if useInstructionMode == enumInstructionModes.None {
		return errorHandler.AddNew(enumErrorCodes.InstUnsupportedMode) // ❌ Fails
	}
	//Overwrite mode with ZP version if possible
	if useInstructionZPMode != enumInstructionModes.None {
		instructionMode = useInstructionZPMode
	}

	//Write data into ROM

	instructionOpcode := opcodesAndSupportedModes.ModeOpcodes[instructionMode]
	operandNeedsNBytes := instructionData.InstructionModeOperandRequiredBytes[instructionMode]
	bytesToInsert := make([]uint8, 0)
	bytesToInsert = append(bytesToInsert, instructionOpcode)

	asRomData, err := nodesToBytes.ConvertNodeValueToUInts(operand, operandNeedsNBytes, false)
	if err != nil {
		return err
	}
	bytesToInsert = append(bytesToInsert, asRomData...)

	err = romBuilder.AddBytesToRom(bytesToInsert)
	if err != nil {
		return err
	}

	return nil
}

// +++++++++++++++++

// Helper to see if instruction is a branch instruction. Only branch instructions contain REL mode
func checkIfBranchInstruction(instructionName string) bool {
	// What modes this instruction can use
	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]

	for _, v := range *opcodesAndSupportedModes.SupportedModes {
		if v == enumInstructionModes.REL {
			return true
		}
	}
	return false
}

// +++++++++++++++++

// Operand wants to use an X/Y index. See if index is used properly and eat it
func checkModeSupportsXY(instructionMode instModeEnum, instructionIndex tokenEnum) (instModeEnum, error) {
	xyMode, ok := instructionData.InstructionModeEnumToXYModeEnum[instructionMode]

	if ok {
		if instructionIndex == enumTokenTypes.REGISTER_X {
			return xyMode.X, nil
		}
		return xyMode.Y, nil
	}

	return instructionMode, errorHandler.AddNew(enumErrorCodes.InstXYUnusableMode) // ❌ Fails
}

// If absolute mode (X,Y too) get the ZP version for auto convert. Get None otherwise
func getModeZeroPageEquivalent(instructionMode instModeEnum) instModeEnum {
	zpMode, ok := instructionData.InstructionABSEnumToZPEnum[instructionMode]
	if ok {
		return zpMode
	}
	return enumInstructionModes.None
}
