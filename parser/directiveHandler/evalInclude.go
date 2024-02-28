package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

func evalInclude(operandList *[]Node) error {
	fileNameNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsString(&fileNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotString) // ❌ Fails
	}
	newFileName := fileHandler.AddRelativePathIncludeFileName(fileNameNode.NodeValue)
	err := fileHandler.OpenInputFileAndPushLinesToStack(newFileName)
	if err != nil {
		return err
	}
	fileHandler.TriggerNewStackCall = true

	return nil
}

// ------------------------------------------------------

func evalIncbin(operandList *[]Node) error {
	var seekValue int = 0
	var readValue int = 0

	fileNameNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsString(&fileNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotString) // ❌ Fails
	}
	newBinFileName := fileHandler.AddRelativePathIncludeFileName(fileNameNode.NodeValue)

	if len(*operandList) >= 2 {
		seekNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsNumeric(&seekNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIsPositive(&seekNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive) // ❌ Fails
		}
		seekValue = seekNode.AsNumber
	}

	if len(*operandList) == 3 {
		readNode := (*operandList)[2]
		if !operandFactory.ValidateNodeIsNumeric(&readNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIsGTEValue(&readNode, 1) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
		}
		readValue = readNode.AsNumber
	}

	byteBuffer, err := fileHandler.ProcessBinFile(newBinFileName, seekValue, readValue)
	if err != nil {
		return err
	}
	err = romBuilder.AddBytesToRom(byteBuffer)
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}
