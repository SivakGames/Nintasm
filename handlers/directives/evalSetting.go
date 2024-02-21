package handlerDirective

import (
	"errors"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/operandFactory"
)

// ---------------------------------

func evalSettingChange(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "SETCHARMAP":
		changeToCharmapName := (*operandList)[0]
		if !operandFactory.ValidateNodeIsIdentifier(&changeToCharmapName) {
			return errors.New("Must use an identifier!")
		}

		newCharmapName := changeToCharmapName.NodeValue
		err := charmapTable.SetCharmapTo__(newCharmapName)
		if err != nil {
			return err
		}
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
	default:
		panic("Unknown reset setting!")
	}
	return nil
}
