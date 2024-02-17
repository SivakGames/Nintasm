package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
)

func evalMacro(directiveName string, macroLabel string, operandList *[]Node) error {
	var err error

	if len(blockStack.Stack) > 0 {
		errMsg := fmt.Sprintf("Cannot define a macro when in another block statement!")
		return errors.New(errMsg) // ❌ Fails
	}

	err = environment.CheckIfSymbolAlreadyDefined(macroLabel)
	if err != nil {
		return err // ❌ Fails
	}

	err = blockStack.SetCurrentOperationLabel(macroLabel)
	if err != nil {
		return err // ❌ Fails
	}

	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnlyFlag()
	return nil
}

func evalEndMacro(directiveName string, operandList *[]Node) error {
	var noLines []blockStack.CapturedLine

	macroLabel := blockStack.GetCurrentOperationLabel()
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines
	environment.AddMacroToEnvironment(macroLabel, *capturedLines)

	blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendCapturedLines(noLines)
	return nil
}
