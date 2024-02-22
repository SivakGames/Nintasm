package directiveHandler

import (
	"errors"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
)

func evalExprmap(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalEndExprmap(directiveName string) error {
	exprmapLabel := blockStack.GetLabelAndDoEndBlockSetups()
	capturedLines := blockStack.GetTopOfStackCapturedLines()
	exprmapTable.AddExprmapToEnvironment(exprmapLabel)

	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}

func evalDefExpr(directiveName string, operandList *[]Node) error {
	exprNode := &(*operandList)[0]
	if !(operandFactory.ValidateNodeIsString(exprNode)) {
		return errors.New("Expression must be a string!")
	}
	exprValueNode := &(*operandList)[1]
	if !operandFactory.ValidateNodeIsNumeric(exprValueNode) ||
		!operandFactory.ValidateNumericNodeIs8BitValue(exprValueNode) {
		return errors.New("Expression value must be a number and be 8 bit!")
	}

	_, err := exprmapTable.CheckIfAlreadyExistsInExprmap(exprValueNode.NodeValue)
	if err != nil {
		return err
	}

	exprmapTable.AddExprToExprmap(exprNode.NodeValue, exprValueNode.AsNumber)
	return nil
}
