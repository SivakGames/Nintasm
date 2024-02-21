package directiveHandler

import (
	"errors"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder/romBuildingSettings"
)

// ---------------------------------

func evalSettingChange(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "SETCHARMAP":
		changeToCharmapNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsIdentifier(&changeToCharmapNode) {
			return errors.New("Must use an identifier!")
		}
		newCharmapName := changeToCharmapNode.NodeValue
		err := charmapTable.SetCharmapTo__(newCharmapName)
		if err != nil {
			return err
		}
	case "SETEXPRMAP":
		changeToExprmapNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsIdentifier(&changeToExprmapNode) {
			return errors.New("Must use an identifier!")
		}
		newExprmapName := changeToExprmapNode.NodeValue
		err := exprmapTable.SetExprmapTo__(newExprmapName)
		if err != nil {
			return err
		}

	case "RSSET":
		EmptyFillSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&EmptyFillSettingNode) ||
			!operandFactory.ValidateNumericNodeIsPositive(&EmptyFillSettingNode) ||
			!operandFactory.ValidateNumericNodeIs16BitValue(&EmptyFillSettingNode) {
			return errors.New("Bad rsset value. Must be a 16-bit positive number")
		}

	case "AUTOZP":
		autoZPSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&autoZPSettingNode) ||
			!operandFactory.ValidateNumericNodeIsGTEandLTEValues(&autoZPSettingNode, 0, 1) {
			return errors.New("Bad auto ZP value. Must be 0 or 1!")
		}
		autoZPOn := autoZPSettingNode.AsNumber == 1
		romBuildingSettings.SetAutoZeroPage(autoZPOn)
		return nil

	case "EMPTYROMFILL":
		EmptyFillSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&EmptyFillSettingNode) ||
			!operandFactory.ValidateNumericNodeIsPositive(&EmptyFillSettingNode) ||
			!operandFactory.ValidateNumericNodeIs8BitValue(&EmptyFillSettingNode) {
			return errors.New("Bad empty fill value. Must be an 8-bit positive number")
		}
		emptyFillValue := uint8(EmptyFillSettingNode.AsNumber)
		romBuildingSettings.SetEmptyRomFillValue(emptyFillValue)
		return nil

	default:
		panic("Unknown reset setting!")
	}
	return nil
}

// ---------------------------------

func evalSettingReset(directiveName string) error {
	switch directiveName {
	case "RESETCHARMAP":
		err := charmapTable.SetCharmapToDefault()
		if err != nil {
			return err
		}
	case "RESETEXPRMAP":
		err := exprmapTable.SetExprmapToDefault()
		if err != nil {
			return err
		}
	default:
		panic("Unknown reset setting!")
	}
	return nil
}
