package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/operandFactory"
)

func evalCharmap(directiveName string, charmapLabel string, operandList *[]Node) error {
	blockStack.PushCaptureBlock(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(charmapLabel, enumSymbolTableTypes.CharMap)
	return nil
}

func evalEndCharmap() error {
	blockStack.SetCurrentCaptureBlockPostFn(endCharmapCleanup)
	blockStack.CopyCapturedLinesToProcessedWithEmptyScope()
	return nil
}

func endCharmapCleanup() {
	blockStack.ClearCurrentOperationLabel()
}

// -----------------------------------------

func evalDefChar(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "DEFCHAR":
		textNode := &(*operandList)[0]

		targetRune, err := validateCharmapTextNodeGetRune(textNode)
		if err != nil {
			return err
		}

		_, err = charmapTable.CheckIfCharAlreadyExistsInCharmap(targetRune)
		if err != nil {
			return err
		}

		charNode := (*operandList)[1]

		charmapTable.AddCharToCharmap(targetRune, charNode)

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
			return errorHandler.AddNew(enumErrorCodes.DefCharRangeEndSmaller)
		}

		charBaseNode := (*operandList)[2]
		if !operandFactory.ValidateNodeIsNumeric(&charBaseNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
		}

		for i := targetStartRune; i <= targetEndRune; i++ {
			charmapTable.AddCharToCharmap(i, charBaseNode)
			charBaseNode.AsNumber += 1
			charBaseNode.NodeValue = fmt.Sprintf("%v", charBaseNode.AsNumber)
		}

	default:
		panic("🛑 Something is very wrong with charmap/defchar capturing!!!")
	}

	return nil
}

func validateCharmapTextNodeGetRune(runeNode *Node) (rune, error) {
	if !operandFactory.ValidateNodeIsString(runeNode) {
		return ' ', errorHandler.AddNew(enumErrorCodes.NodeTypeNotString) // ❌ Fails
	}
	runeArray := []rune(runeNode.NodeValue)
	if len(runeArray) != 1 {
		return ' ', errorHandler.AddNew(enumErrorCodes.DefCharTooLong, runeNode.NodeValue)
	}
	return runeArray[0], nil
}
