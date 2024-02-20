package handlerDirective

import (
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
)

func evalExprmap(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalEndExprmap(directiveName string) error {

	charmapLabel := blockStack.GetCurrentOperationLabel()
	_ = charmapLabel
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines

	environment.AddCharmapToEnvironment(charmapLabel)
	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}
