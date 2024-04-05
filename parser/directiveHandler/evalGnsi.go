package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
	"strings"
)

func evalGnsi(operationLabel string, operandList *[]Node) error {
	err := environment.AddOtherIdentifierToMasterTable(operationLabel, enumSymbolTableTypes.Namespace)
	if err != nil {
		return err
	}

	gnsiLabelNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(&gnsiLabelNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
	}

	gnsiResolveSize := 0
	if len(*operandList) > 1 {
		gnsiResolveSizeNode := (*operandList)[1]
		if !operandFactory.ValidateNumericNodeIsGTEandLTEValues(&gnsiResolveSizeNode, 0, 1) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTEandLTE, 0, 1)
		}
		gnsiResolveSize = int(gnsiResolveSizeNode.AsNumber)
	}

	gnsiLabelName := gnsiLabelNode.NodeValue

	symbolTableEnum, exists := environment.CheckIfIdentifierExistsInMasterTable(gnsiLabelName)
	if !exists {
		return errorHandler.AddNew(enumErrorCodes.Other, "Not found")
	}
	if symbolTableEnum != enumSymbolTableTypes.SymbolAsNode {
		return errorHandler.AddNew(enumErrorCodes.Other, "Must be a label")
	}
	if strings.Contains(gnsiLabelName, ".") {
		return errorHandler.AddNew(enumErrorCodes.Other, "No local or parent/locals allowed")
	}
	localLabelsFromParent := environment.GetLocalLabelsOfParent(gnsiLabelName)
	if len(localLabelsFromParent) == 0 {
		return errorHandler.AddNew(enumErrorCodes.Other, "No local labels...")
	}

	parentNode, _, _ := environment.LookupIdentifierInSymbolAsNodeTable(gnsiLabelName)

	for _, localLabel := range localLabelsFromParent {
		localLabelNode, _, _ := environment.LookupIdentifierInSymbolAsNodeTable(localLabel)
		index := strings.Index(localLabel, ".")
		isolatedLabel := localLabel[index+1:]
		difference := int(localLabelNode.AsNumber-parentNode.AsNumber) >> gnsiResolveSize
		differenceNode := operandFactory.CreateNumericLiteralNode(float64(difference))
		newName := operationLabel + "." + isolatedLabel
		err := environment.AddIdentifierToSymbolAsNodeTable(newName, differenceNode)
		if err != nil {
			return err
		}
	}

	return nil
}
