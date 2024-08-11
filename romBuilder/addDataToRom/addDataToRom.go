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
	isBranch := instructionMode == enumInstructionModes.REL
	asRomData, err := nodesToBytes.ConvertNodeValueToUInts(operand, operandByteSize, isBranch, false, true)
	if err != nil {
		return err // ❌ Fails
	}

	// Will try and resolve again in pass 2
	if !operand.Resolved {
		unresolvedTable.AddUnresolvedRomEntry(operand, operandByteSize, isBranch, false)
	}

	//Warning about string operands used for modes that aren't immediate
	if operand.NodeType == enumNodeTypes.StringLiteral && instructionMode != enumInstructionModes.IMM {
		errorHandler.AddNew(enumErrorCodes.ResolvedValueIsStringForInst) // ⚠️ Warns
	}

	err = romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}

//------------------------------------------

func AddRawBytesToRom(operand Node, operandByteSize int, isBigEndian bool) error {
	asRomData, err := nodesToBytes.ConvertNodeValueToUInts(operand, operandByteSize, false, isBigEndian, false)
	if err != nil {
		return err // ❌ Fails
	}

	// Will try and resolve again in pass 2
	if !operand.Resolved {
		unresolvedTable.AddUnresolvedRomEntry(operand, operandByteSize, false, isBigEndian)
	}

	err = romBuilder.AddBytesToRom(asRomData)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}
