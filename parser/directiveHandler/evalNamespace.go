package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment/namespaceTable"
)

func evalNamespace(directiveName string, namespaceLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	namespaceTable.AddNamespaceToEnvironment(namespaceLabel)
	interpreter.AppendParentLabel(namespaceLabel)

	blockStack.SetCaptureParentOpOnlyFlag()
	blockStack.SetCurrentOperationEvaluatesFlag()
	return nil
}

// End the Namespace definition and add to environment
func evalEndNamespace() error {
	_ = blockStack.GetLabelAndDoEndBlockSetups()
	interpreter.PopParentLabel()
	blockStack.PopFromStackAndExtendNoLines()
	return nil
}
