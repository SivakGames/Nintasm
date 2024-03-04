package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
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
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
		}
		newCharmapName := changeToCharmapNode.NodeValue
		err := charmapTable.SetCharmapTo__(newCharmapName)
		if err != nil {
			return err
		}

	case "SETEXPRMAP":
		changeToExprmapNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsIdentifier(&changeToExprmapNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
		}
		newExprmapName := changeToExprmapNode.NodeValue
		err := exprmapTable.SetExprmapTo__(newExprmapName)
		if err != nil {
			return err
		}

	case "RSSET":
		RSSetSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&RSSetSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIsPositive(&RSSetSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIs16BitValue(&RSSetSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNot16Bit) // ❌ Fails
		}
		rssetNumber := RSSetSettingNode.AsNumber
		romBuildingSettings.SetRSValue(uint(rssetNumber))

	case "AUTOZP":
		autoZPSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&autoZPSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIsGTEandLTEValues(&autoZPSettingNode, 0, 1) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTEandLTE, 0, 1) // ❌ Fails
		}
		autoZPOn := autoZPSettingNode.AsNumber == 1
		romBuildingSettings.SetAutoZeroPage(autoZPOn)
		return nil

	case "EMPTYROMFILL":
		EmptyFillSettingNode := (*operandList)[0]
		if !operandFactory.ValidateNodeIsNumeric(&EmptyFillSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIsPositive(&EmptyFillSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive) // ❌ Fails
		} else if !operandFactory.ValidateNumericNodeIs8BitValue(&EmptyFillSettingNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNot8Bit) // ❌ Fails
		}
		emptyFillValue := uint8(EmptyFillSettingNode.AsNumber)
		romBuildingSettings.SetEmptyRomFillValue(emptyFillValue)
		return nil

	default:
		panic("🛑 Unknown setting!")
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
		panic("🛑 Unknown reset setting!")
	}
	return nil
}
