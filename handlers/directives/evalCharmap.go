package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/parser/operandFactory"
)

func evalCharmap(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalEndCharmap(directiveName string) error {

	charmapLabel := blockStack.GetCurrentOperationLabel()
	_ = charmapLabel
	blockStack.ClearCurrentOperationLabel()
	blockStack.ClearCaptureParentOpOnlyFlag()

	currentStackOp := blockStack.GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines

	environment.AddCharmapToEnvironment(charmapLabel)
	blockStack.PopFromStackAndExtendCapturedLines(*capturedLines)
	return nil
}

func evalDefChar(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "DEFCHAR":
		runeNode := &(*operandList)[0]
		if !operandFactory.ValidateNodeIsString(runeNode) {
			return errors.New("First operand must be a string!")
		}
		runeArray := []rune(runeNode.NodeValue)
		if len(runeArray) != 1 {
			return errors.New("Character definition must be 1 character long!")
		}
		targetRune := runeArray[0]

		charValue := &(*operandList)[1]
		if !operandFactory.ValidateNodeIsNumeric(charValue) ||
			!operandFactory.ValidateNumericNodeIsPositive(charValue) ||
			!operandFactory.ValidateNumericNodeIs8BitValue(charValue) {
			return errors.New("Node must be positive, 8 bit, and numeric...")
		}

		environment.AddCharToCharmap(targetRune, {charValue})

		fmt.Println(directiveName, operandList)

	case "DEFCHARRANGE":
	default:
		panic("Something is very wrong with charmap/defchar capturing!!!")
	}

	return nil
}
