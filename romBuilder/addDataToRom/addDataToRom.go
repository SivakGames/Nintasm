package addDataToRom

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumInstructionModes "misc/nintasm/constants/enums/instructionModes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment/unresolvedTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
)

type Node = operandFactory.Node

//------------------------------------------

func AddInstructionOperandToRom(operand Node, operandByteSize int, instructionMode enumInstructionModes.Def) error {
	asRomData, err := nodesToBytes.ConvertNodeValueToUInts(operand, operandByteSize, false)
	if err != nil {
		return err // ❌ Fails
	}

	// Will try and resolve again in pass 2
	if !operand.Resolved {
		unresolvedTable.AddUnresolvedRomEntry(operand, operandByteSize)
	}

	//Warning about string operands used for modes that aren't immediate
	if operand.NodeType == enumNodeTypes.StringLiteral && instructionMode != enumInstructionModes.IMM {
		errorHandler.AddNew(enumErrorCodes.ResolvedValueIsStringForInst) // ⚠️ Warns
	}

	//Multibytes are technially OK, but have to fit
	if len(asRomData) > operandByteSize {
		return errorHandler.AddNew(enumErrorCodes.ResolvedValueTooBig, len(asRomData), operandByteSize) // ❌ Fails
	}

	err = romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}

//------------------------------------------

func AddRawBytesToRom(operand Node, operandByteSize int, isBigEndian bool, byteSizeIsExact bool) error {
	asRomData, err := nodesToBytes.ConvertNodeValueToUInts(operand, operandByteSize, isBigEndian)
	if err != nil {
		return err // ❌ Fails
	}

	// Will try and resolve again in pass 2
	if !operand.Resolved {
		unresolvedTable.AddUnresolvedRomEntry(operand, operandByteSize)
	}

	err = romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}
