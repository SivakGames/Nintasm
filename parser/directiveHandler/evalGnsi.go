package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"strings"
)

func evalGnsi(operationLabel string, operandList *[]Node) error {
	err := environment.AddOtherIdentifierToMasterTable(operationLabel, enumSymbolTableTypes.Namespace)
	if err != nil {
		return err
	}

	unresFlag := environment.GetUnresolvedSilentErrorFlag()
	environment.ClearUnresolvedSilentErrorFlag()

	defer environment.SetUnresolvedSilentErrorFlagTo(unresFlag)

	gnsiLabelNode := (*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(&gnsiLabelNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
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

	var gnsiResolveOpNode *Node = nil
	if len(*operandList) > 1 {
		gnsiResolveOpNode = &(*operandList)[1]
	}

	for _, localLabel := range localLabelsFromParent {
		localLabelNode, _, _ := environment.LookupIdentifierInSymbolAsNodeTable(localLabel)
		localLabelStartingIndex := strings.Index(localLabel, ".")
		isolatedLocalLabel := localLabel[localLabelStartingIndex+1:]

		var finalNode Node

		difference := localLabelNode.AsNumber - parentNode.AsNumber

		if gnsiResolveOpNode == nil {
			finalNode = operandFactory.CreateNumericLiteralNode(difference)
		} else {
			finalNode, err = evalGnsiChild(difference, gnsiResolveOpNode)
			if err != nil {
				return err
			}
		}

		newName := operationLabel + "." + isolatedLocalLabel
		err = environment.AddIdentifierToSymbolAsNodeTable(newName, finalNode, enumSymbolTableTypes.SymbolAsNode)
		if err != nil {
			return err
		}
	}

	return nil
}

func evalGnsiChild(difference float64, gnsiResolveOpNode *Node) (Node, error) {
	differenceNode := operandFactory.CreateNumericLiteralNode(difference)

	symbolAsNodeTable.PushToSymbolTableStack()
	defer symbolAsNodeTable.PopFromSymbolTableStack()
	symbolAsNodeTable.AddSymbolToTopTableStack("\\1", differenceNode)
	finalNode, err := interpreter.EvaluateNode(*gnsiResolveOpNode)

	return finalNode, err
}
