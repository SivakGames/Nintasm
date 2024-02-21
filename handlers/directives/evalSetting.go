package handlerDirective

import (
	"errors"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/operandFactory"
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
