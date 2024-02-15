package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
)

func evalMacro(directiveName string, operandList *[]Node) error {
	if len(blockStack.Stack) > 0 {
		errMsg := fmt.Sprintf("Cannot define a macro when in another block statement!")
		return errors.New(errMsg)
	}
	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnly()
	return nil
}

func evalEndMacro(directiveName string, operandList *[]Node) error {
	var noLines []blockStack.CapturedLine
	blockStack.ClearCaptureParentOpOnly()
	blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendCapturedLines(noLines)
	return nil
}
