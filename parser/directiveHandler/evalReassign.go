package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
)

func evalReassign(operandList *[]Node) error {
	identifierNode := &(*operandList)[0]
	reassignValueNode := &(*operandList)[1]
	reassignedIdentifier := identifierNode.NodeValue

	if !operandFactory.ValidateNodeIsIdentifier(identifierNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier)
	}
	symbolEnum, exists := environment.CheckIfIdentifierExistsInMasterTable(reassignedIdentifier)
	if !exists {
		return errorHandler.AddNew(enumErrorCodes.ReassignmentIdentifierNotFound, reassignedIdentifier)
	}

	if symbolEnum != enumSymbolTableTypes.SymbolAsNode {
		return errorHandler.AddNew(enumErrorCodes.ReassignmentIdentifierNotSymbol, reassignedIdentifier, environment.GetSymbolDescriptionFromEnum(symbolEnum))
	}

	evaluatedReassignNode, err := interpreter.EvaluateNode(*reassignValueNode)
	if err != nil {
		return err
	}

	symbolAsNodeTable.AddIdentifierKeyToSymbolAsNodeTable(reassignedIdentifier, evaluatedReassignNode)

	return nil
}
