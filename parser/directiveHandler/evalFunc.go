package directiveHandler

import (
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/funcTable"
)

func evalFunc(operationLabel string, operandList *[]Node) error {
	functionNode := (*operandList)[0]
	environment.AddOtherIdentifierToMasterTable(operationLabel, enumSymbolTableTypes.Function)

	funcTable.AddNodeToFunction(operationLabel, &functionNode)

	return nil
}
