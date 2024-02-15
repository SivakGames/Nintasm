package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
)

func evalMacro(directiveName string, operandList *[]Node) error {
	if len(blockStack.Stack) > 0 {
		errMsg := fmt.Sprintf("Cannot define a macro when in another block statement!")
		return errors.New(errMsg)
	}
	environment.LookupMacroInEnvironment("GSG")

	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnlyFlag()
	return nil
}

func evalEndMacro(directiveName string, operandList *[]Node) error {
	var noLines []blockStack.CapturedLine

	blockStack.ClearCaptureParentOpOnlyFlag()
	currentStackOp := blockStack.GetCurrentOperation()
	capturedLines := &currentStackOp.CapturedLines
	environment.AddMacroToEnvironment("GSG", *capturedLines)
	blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendCapturedLines(noLines)
	return nil
}
