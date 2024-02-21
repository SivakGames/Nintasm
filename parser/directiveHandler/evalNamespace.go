package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
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
	blockStack.SetTemporaryOverwritingParentLabel(namespaceLabel)

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines

	if len(*capturedLines) == 0 {
		fmt.Println("Warning: Namespace is empty!")
	}

	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}
