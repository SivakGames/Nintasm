package handlerDirective

import (
	"fmt"
	"misc/nintasm/handlers/blockStack"
)

func evalCharmap(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalEndCharmap(directiveName string) error {

	charmapLabel := blockStack.GetCurrentOperationLabel()
	_ = charmapLabel
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines

	//environment.AddMacroToEnvironment(charmapLabel, *capturedLines)
	//blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}

func evalDefChar(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "DEFCHAR":
		fmt.Println(directiveName, operandList)
	case "DEFCHARRANGE":
	default:
		panic("Somethis is very wrong with charmap/defchar capturing!!!")
	}

	return nil
}
