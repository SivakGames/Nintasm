package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
)

func evalExprmap(directiveName string, exprmapLabel string, operandList *[]Node) error {
	blockStack.PushOntoTopEntry(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(exprmapLabel, enumSymbolTableTypes.ExprMap)
	return nil
}

func evalEndExprmap() error {
	blockStack.EndLabeledDirective()
	return nil
}

func evalDefExpr(directiveName string, operandList *[]Node) error {
	exprNode := &(*operandList)[0]
	if !(operandFactory.ValidateNodeIsString(exprNode)) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotString) // ❌ Fails
	}
	exprValueNode := &(*operandList)[1]
	if !operandFactory.ValidateNodeIsNumeric(exprValueNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIs8BitValue(exprValueNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNot8Bit) // ❌ Fails
	}

	_, err := exprmapTable.CheckIfAlreadyExistsInExprmap(exprValueNode.NodeValue)
	if err != nil {
		return err
	}

	exprmapTable.AddExprToExprmap(exprNode.NodeValue, exprValueNode.AsNumber)
	return nil
}
