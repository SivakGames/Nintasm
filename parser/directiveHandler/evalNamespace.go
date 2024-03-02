package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
)

func evalNamespace(directiveName string, namespaceLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(namespaceLabel, enumSymbolTableTypes.Namespace)
	interpreter.AppendParentLabel(namespaceLabel)

	blockStack.SetCaptureParentOpOnlyFlag()
	blockStack.SetCurrentOperationEvaluatesCapturedNodesFlag()
	return nil
}

// End the Namespace definition and add to environment
func evalEndNamespace() error {
	_ = blockStack.GetLabelAndDoEndBlockSetups()
	blockStack.ClearCurrentOperationEvaluatesCapturedNodesFlag()
	interpreter.PopParentLabel()
	blockStack.PopFromStackAndExtendNoLines()
	return nil
}
