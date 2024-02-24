package directiveHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/fileStack"
	"misc/nintasm/interpreter/operandFactory"
)

func evalInclude(operandList *[]Node) error {
	fmt.Println(operandList)
	return nil
}

func evalIncbin(operandList *[]Node) error {
	fileNameNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsString(&fileNameNode) {
		return errors.New("Input file has to be a string!!!")
	}
	newBinFileName := fileStack.AddRelativePathIncludeFileName(fileNameNode.NodeValue)
	_, err := fileStack.OpenBinFile(newBinFileName, 1, 10)
	if err != nil {
		return err
	}

	fmt.Println(newBinFileName)
	return nil
}
