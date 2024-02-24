package directiveHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder/romBuildingSettings"
)

func evalRs(operandList *[]Node) error {
	rsValueNode := (*operandList)[0]

	if !(operandFactory.ValidateNodeIsNumeric(&rsValueNode) &&
		operandFactory.ValidateNumericNodeIsGTEValue(&rsValueNode, 1)) {
		return errors.New("RS value must be numeric and >= 1") // ❌ Fails
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
