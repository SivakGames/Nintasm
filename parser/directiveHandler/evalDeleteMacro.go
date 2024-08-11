package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
)

func evalDeleteMacro(operandList *[]Node) error {
	macroNameNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(&macroNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
	}
	macroName := macroNameNode.NodeValue

	environment.RemoveOtherIdentifierFromMasterTable(macroName, enumSymbolTableTypes.Macro)
	return nil
}
