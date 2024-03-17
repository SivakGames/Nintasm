package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/namespaceTable"
)

func evalNamespace(directiveName string, namespaceLabel string, operandList *[]Node) error {
	blockStack.PushCaptureBlock(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(namespaceLabel, enumSymbolTableTypes.Namespace)
	interpreter.AppendParentLabel(namespaceLabel)
	namespaceTable.IsDefiningNamespace = true
	return nil
}

// End the Namespace definition and add to environment
func evalEndNamespace() error {
	blockStack.ProcessEndLabeledDirective()
	interpreter.PopParentLabel()
	namespaceTable.IsDefiningNamespace = false
	return nil
}
