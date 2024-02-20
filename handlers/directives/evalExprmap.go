package handlerDirective

import (
	"errors"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
)

func evalExprmap(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalEndExprmap(directiveName string) error {

	exprmapLabel := blockStack.GetCurrentOperationLabel()
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines

	environment.AddExprmapToEnvironment(exprmapLabel)
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

	_, err := environment.CheckIfAlreadyExistsInExprmap(exprValueNode.NodeValue)
	if err != nil {
		return err
	}

	environment.AddExprToExprmap(exprNode.NodeValue, exprValueNode.AsNumber)
	return nil
}