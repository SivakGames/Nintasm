package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder/romBuildingSettings"
)

func evalRs(operandList *[]Node) error {
	rsValueNode := (*operandList)[0]

	if !operandFactory.ValidateNodeIsNumeric(&rsValueNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIsGTEValue(&rsValueNode, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
	}

	rsCurrentValue, err := romBuildingSettings.GetRSValue()
	if err != nil {
		return err
	}

	_ = rsCurrentValue
	fmt.Println(rsCurrentValue)

	romBuildingSettings.AddToRSValue(uint(rsValueNode.AsNumber))

	return nil
}
