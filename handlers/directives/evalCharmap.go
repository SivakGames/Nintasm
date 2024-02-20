package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
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
		textNode := &(*operandList)[0]

		targetRune, err := validateCharmapTextNodeGetRune(textNode)
		if err != nil {
			return err
		}

		_, err = environment.CheckIfAlreadyExistsInCharmap(targetRune)
		if err != nil {
			return err
		}

		charNodes := []Node{}
		for _, charNode := range (*operandList)[1:] {
			err = validateCharmapNumberNode(&charNode)
			if err != nil {
				return err
			}
			charNodes = append(charNodes, charNode)
		}

		environment.AddCharToCharmap(targetRune, charNodes)

	case "DEFCHARRANGE":
		textNodeStart := &(*operandList)[0]
		targetStartRune, err := validateCharmapTextNodeGetRune(textNodeStart)
		if err != nil {
			return err
		}
		textNodeEnd := &(*operandList)[1]
		targetEndRune, err := validateCharmapTextNodeGetRune(textNodeEnd)
		if err != nil {
			return err
		}

		if targetStartRune >= targetEndRune {
			return errors.New("End bigger than start")
		}

		charBaseNode := (*operandList)[2]

		for i := targetStartRune; i <= targetEndRune; i++ {
			environment.AddCharToCharmap(i, []Node{charBaseNode})
			charBaseNode.AsNumber += 1
			charBaseNode.NodeValue = fmt.Sprintf("%v", charBaseNode.AsNumber)
		}

	default:
		panic("Something is very wrong with charmap/defchar capturing!!!")
	}

	return nil
}

func validateCharmapTextNodeGetRune(runeNode *Node) (rune, error) {
	if !operandFactory.ValidateNodeIsString(runeNode) {
		return ' ', errors.New("First operand must be a string!")
	}
	runeArray := []rune(runeNode.NodeValue)
	if len(runeArray) != 1 {
		errMsg := fmt.Sprintf("%v Character definition must be 1 character long!", runeNode.NodeValue)
		return ' ', errors.New(errMsg)
	}
	return runeArray[0], nil
}

func validateCharmapNumberNode(node *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(node) ||
		!operandFactory.ValidateNumericNodeIsPositive(node) ||
		!operandFactory.ValidateNumericNodeIs8BitValue(node) {
		return errors.New("Node must be positive, 8 bit, and numeric...")
	}
	return nil
}
