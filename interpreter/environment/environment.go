package environment

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/operandFactory"
)

type Node = operandFactory.Node

// ----------------------------------

func generateAssemblerReservedWordNode(funcName string) Node {
	return operandFactory.CreateAssemblerReservedWordNode(funcName)
}

func generateNumericNodeForEnvironment(number int) Node {
	return operandFactory.CreateNumericLiteralNode(enumTokenTypes.NUMBER_decimal, fmt.Sprintf("%d", number), number)
}

// ----------------------------------

func CheckIfSymbolAlreadyDefined(symbolName string) error {
	var exists bool
	errMsgBase := "Symbol %v has been previously defined! (Defined as %v)"

	_, exists = literalNodeSymbolTable[symbolName]
	if exists {
		errMsg := fmt.Sprintf(errMsgBase, symbolName, "literal")
		return errors.New(errMsg)
	}
	_, exists = macroSymbolTable[symbolName]
	if exists {
		errMsg := fmt.Sprintf(errMsgBase, symbolName, "macro")
		return errors.New(errMsg)
	}
	return nil
}

// ----------------------------------

func AddToEnvironment(symbolName string, node Node) (Node, error) {
	_, exists := literalNodeSymbolTable[symbolName]
	if exists {
		errMsg := fmt.Sprintf("%v : SYMBOL IS ALREADY DEFINED!", symbolName)
		return node, errors.New(errMsg)

	} else {
		literalNodeSymbolTable[symbolName] = node
	}
	return literalNodeSymbolTable[symbolName], nil
}

func LookupInEnvironment(symbolName string) (Node, error) {
	value, ok := literalNodeSymbolTable[symbolName]
	if ok {
		return value, nil
	} else {
		errMsg := fmt.Sprintf("\x1b[33m%v\x1b[0m was not found", symbolName)
		return operandFactory.EmptyNode(), errors.New(errMsg)
	}
}

// ----------------------------------

func AddMacroToEnvironment(symbolName string, capturedLines MacroTableType) error {
	macroSymbolTable[symbolName] = capturedLines
	return nil
}

func LookupMacroInEnvironment(symbolName string) (MacroTableType, bool) {
	macro, ok := macroSymbolTable[symbolName]
	return macro, ok
}

func LookupAndGetMacroInEnvironment(symbolName string) (MacroTableType, error) {
	macro, ok := LookupMacroInEnvironment(symbolName)
	if ok {
		return macro, nil
	} else {
		errMsg := fmt.Sprintf("Macro %v does not exist!", symbolName)
		return nil, errors.New(errMsg)
	}
}

// ----------------------------------

var defaultCharmapName = ""
var currentCharmapName = ""
var lastAddedCharmapName = ""

func AddCharmapToEnvironment(symbolName string) error {
	charmapSymbolTable[symbolName] = CharmapTableType{}
	lastAddedCharmapName = symbolName

	if defaultCharmapName == "" {
		defaultCharmapName = lastAddedCharmapName
		currentCharmapName = defaultCharmapName
	}
	return nil
}

func GetCurrentCharmap() (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errors.New("No charmaps have been defined!!!!")
	}

	return charmapSymbolTable[currentCharmapName], nil
}

func CheckIfDefinedInCharmap(lookupChar rune) ([]Node, bool) {
	target, exists := charmapSymbolTable[lastAddedCharmapName][lookupChar]
	return target, exists
}

func CheckIfAlreadyExistsInCharmap(lookupChar rune) ([]Node, error) {
	target, exists := CheckIfDefinedInCharmap(lookupChar)
	if exists {
		errMsg := fmt.Sprintf("Target char %c already defined!", lookupChar)
		return target, errors.New(errMsg)
	}
	return target, nil
}

func AddCharToCharmap(newChar rune, charNodes []Node) error {
	charmapSymbolTable[lastAddedCharmapName][newChar] = charNodes
	return nil
}

// ----------------------------------

var defaultExprmapName = ""
var currentExprmapName = ""
var lastAddedExprmapName = ""

func AddExprmapToEnvironment(symbolName string) error {
	exprmapSymbolTable[symbolName] = ExprmapTableType{}
	lastAddedExprmapName = symbolName

	if defaultExprmapName == "" {
		defaultExprmapName = lastAddedExprmapName
		currentExprmapName = defaultExprmapName
	}
	return nil
}

func GetCurrentExprmap() (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errors.New("No exprmaps have been defined!!!!")
	}

	return exprmapSymbolTable[currentExprmapName], nil
}

func CheckIfDefinedInExprmap(lookupExpr string) (int, bool) {
	target, exists := exprmapSymbolTable[lastAddedExprmapName][lookupExpr]
	return target, exists
}

func CheckIfAlreadyExistsInExprmap(lookupExpr string) (int, error) {
	target, exists := CheckIfDefinedInExprmap(lookupExpr)
	if exists {
		errMsg := fmt.Sprintf("Target expr %v already defined in exprmap!", lookupExpr)
		return target, errors.New(errMsg)
	}
	return target, nil
}

func AddExprToExprmap(newExpr string, exprValue int) error {
	exprmapSymbolTable[lastAddedExprmapName][newExpr] = exprValue
	return nil
}
