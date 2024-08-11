package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func evalThrow(operandList *[]Node) error {
	throwNode := &(*operandList)[0]
	if !(operandFactory.ValidateNodeIsString(throwNode)) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotString) // ❌ Fails
	}
	return errorHandler.AddNew(enumErrorCodes.Other, throwNode.NodeValue)
}
