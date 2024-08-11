package instructionHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumInstructionModes "misc/nintasm/constants/enums/instructionModes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/constants/instructionData"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/addDataToRom"
	"strings"
)

type instModeEnum = enumInstructionModes.Def
type tokenEnum = enumTokenTypes.Def
type Node = operandFactory.Node

// -------------------------------------------

// Handle the evaluation of the instruction
func EvaluateInstruction(instructionName string,
	operandList *[]Node,
	instructionMode instModeEnum,
	instructionXYTokenEnum tokenEnum) error {
	var operand Node
	var err error

	// Special reassign for branches
	if instructionMode == enumInstructionModes.ABS && checkIfBranchInstruction(instructionName) {
		instructionMode = enumInstructionModes.REL
	}

	if len(*operandList) == 1 {
		unevaluatedNode := (*operandList)[0]
		if instructionMode == enumInstructionModes.REL {
			branchNode := operandFactory.ConvertToBranchBinaryExpressionNode(unevaluatedNode, romBuilder.GetOrg())
			unevaluatedNode = branchNode
		}

		operand, err = interpreter.EvaluateNode(unevaluatedNode)
		if err != nil {
			err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
			if err != nil {
				return err // ❌❌ CONTINUES Failing!
			}
		}

	} else {
		operand = operandFactory.EmptyNode()
	}

	opcodesAndSupportedModes := instructionData.OpcodesAndSupportedModes[instructionName]
	useInstructionMode, useInstructionZPMode := enumInstructionModes.None, enumInstructionModes.None

	//If an index is present, see if it's usable with desired mode and reassign mode to it
	if instructionXYTokenEnum != enumTokenTypes.None {
		instructionMode, err = checkModeSupportsXYAndReassign(instructionMode, instructionXYTokenEnum)
		if err != nil {
			return err // ❌ Fails
		}
	}

	//Used for auto ZP convert if possible
	instructionZPModeEquivalent := getModeZeroPageEquivalent(instructionMode)

	// Check if instruction itself supports mode
	for _, supportedMode := range *opcodesAndSupportedModes.SupportedModes {
		if supportedMode == instructionMode {
			useInstructionMode = supportedMode
			continue
		}

		// Do ZP equivalent check and auto convert if possible
		if supportedMode == instructionZPModeEquivalent &&
			operand.Resolved &&
			operand.NodeType == enumNodeTypes.NumericLiteral {
			if operandFactory.ValidateNumericNodeIs8BitValue(&operand) {
				useInstructionZPMode = supportedMode
			}
		}
	}

	// See if the mode the instruction is going to use was valid
	if useInstructionMode == enumInstructionModes.None {
		modeDetails := instructionData.InstructionModeEnumDetails[instructionMode]
		suppModeStrings := make([]string, len(*opcodesAndSupportedModes.SupportedModes))
		for i, supportedMode := range *opcodesAndSupportedModes.SupportedModes {
			x := instructionData.InstructionModeEnumDetails[supportedMode]
			suppModeStrings[i] = x.Abbrev
		}
		result := strings.Join(suppModeStrings, ", ")

		errorHandler.AddHint(enumErrorCodes.InstUnsupportedMode, result)
		return errorHandler.AddNew(enumErrorCodes.InstUnsupportedMode, modeDetails.Abbrev) // ❌ Fails
	}
	//Overwrite mode with ZP version if possible
	if useInstructionZPMode != enumInstructionModes.None {
		instructionMode = useInstructionZPMode
	}

	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//Write data into ROM

	instructionOpcode := opcodesAndSupportedModes.ModeOpcodes[instructionMode]
	operandByteSize := instructionData.InstructionModeOperandRequiredBytes[instructionMode]

	//++++++++++++++++++++++++++++++++++
	//Add the opcode to ROM
	opcodeByteToInsert := make([]uint8, 1)
	opcodeByteToInsert[0] = instructionOpcode
	err = romBuilder.AddBytesToRom(opcodeByteToInsert)
	if err != nil {
		return err // ❌ Fails
	}

	//Add the operand (if any) to ROM
	if operandByteSize > 0 {
		err := addDataToRom.AddInstructionOperandToRom(operand, operandByteSize, instructionMode)
		if err != nil {
			return err // ❌ Fails
		}
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
func checkModeSupportsXYAndReassign(instructionMode instModeEnum, instructionXYTokenEnum tokenEnum) (instModeEnum, error) {
	if instructionXYTokenEnum != enumTokenTypes.REGISTER_X && instructionXYTokenEnum != enumTokenTypes.REGISTER_Y {
		panic("🛑 XY index checks didn't yield an X nor a Y!!!")
	}

	xyMode, ok := instructionData.InstructionModeEnumToXYModeEnum[instructionMode]
	if ok {
		if instructionXYTokenEnum == enumTokenTypes.REGISTER_X {
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
