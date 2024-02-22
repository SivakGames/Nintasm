package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment/namespaceTable"
)

func evalNamespace(directiveName string, NamespaceLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnlyFlag()
	return nil
}

// End the Namespace definition and add to environment
func evalEndNamespace(directiveName string) error {
	namespaceLabel := blockStack.GetCurrentOperationLabel()
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines
	namespaceTable.AddNamespaceToEnvironment(namespaceLabel)

	interpreter.PopParentLabelWhenBlockOpDone = true
	interpreter.AppendParentLabel(namespaceLabel)

	if len(*capturedLines) == 0 {
		fmt.Println("Warning: Namespace is empty!")
	}

	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}
