package directiveHandler

import (
	"errors"
	"misc/nintasm/assemble/fileHandler"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

func evalInclude(operandList *[]Node) error {
	fileNameNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsString(&fileNameNode) {
		return errors.New("Input file has to be a string!!!")
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
		return errors.New("Input file has to be a string!!!")
	}
	newBinFileName := fileHandler.AddRelativePathIncludeFileName(fileNameNode.NodeValue)

	if len(*operandList) >= 2 {
		seekNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsNumeric(&seekNode) ||
			!operandFactory.ValidateNumericNodeIsPositive(&seekNode) {
			return errors.New("If setting seek value, value must be numeric and positive!")
		}
		seekValue = seekNode.AsNumber
	}

	if len(*operandList) == 3 {
		readNode := (*operandList)[2]
		if !operandFactory.ValidateNodeIsNumeric(&readNode) ||
			!operandFactory.ValidateNumericNodeIsGTEValue(&readNode, 1) {
			return errors.New("If setting read value, value must be numeric and >= 1!")
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
