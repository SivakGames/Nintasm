package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder/romBuildingSettings"
)

func evalRs(operationLabel string, operandList *[]Node) error {
	rsValueNode := (*operandList)[0]

	if !operandFactory.ValidateNodeIsNumeric(&rsValueNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIsGTEValue(&rsValueNode, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
	}

	rsCurrentValue, rserr := romBuildingSettings.GetRSValue()
	if rserr != nil {
		return rserr
	}

	fmt.Println(rsCurrentValue, operationLabel, rsValueNode)

	rsIdentifierNode := operandFactory.CreateIdentifierNode(operationLabel)
	rsCurrentValueNode := operandFactory.CreateNumericLiteralNode(int(rsCurrentValue))
	rsAssignmentNode := operandFactory.CreateAssignmentNode(rsIdentifierNode, rsCurrentValueNode)
	_, err := interpreter.EvaluateNode(rsAssignmentNode)
	if err != nil {
		return err
	}

	romBuildingSettings.AddToRSValue(uint(rsValueNode.AsNumber))

	return nil
}
