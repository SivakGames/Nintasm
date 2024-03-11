package directiveHandler

import (
	"misc/nintasm/assemble/blockStack2"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
)

func evalNamespace(directiveName string, namespaceLabel string, operandList *[]Node) error {
	blockStack2.PushOntoTopEntry(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(namespaceLabel, enumSymbolTableTypes.Namespace)
	interpreter.AppendParentLabel(namespaceLabel)
	return nil
}

// End the Namespace definition and add to environment
func evalEndNamespace() error {
	blockStack2.ClearCurrentOperationLabel() //TODO - Change to new implementation
	blockStack2.ForcePopTopEntry()
	interpreter.PopParentLabel()
	return nil
}
