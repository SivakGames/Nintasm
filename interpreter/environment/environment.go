package environment

import (
	"errors"
	"fmt"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
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
func decodeHigh(node Node) error {
	if node.NodeType != enumNodeTypes.NumericLiteral {
		return errors.New("High node must be a number")
	}
	return nil //node.AsNumber
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
var lastAddedCharmapName = ""

func AddCharmapToEnvironment(symbolName string) error {
	charmapSymbolTable[symbolName] = CharmapTableType{}
	lastAddedCharmapName = symbolName

	if defaultCharmapName == "" {
		defaultCharmapName = lastAddedCharmapName
	}
	return nil
}

func CheckIfDefinedInCharmap(lookupChar rune) ([]Node, bool) {
	target, exists := charmapSymbolTable[lastAddedCharmapName][lookupChar]
	return target, exists
}

func CheckIfAlreadyExistsInCharmap(lookupChar rune) ([]Node, error) {
	target, exists := CheckIfDefinedInCharmap(lookupChar)
	if exists {
		return target, errors.New("Target char already defined!")
	}
	return target, nil
}

func AddCharToCharmap(newChar rune, charNodes []Node) error {
	charmapSymbolTable[lastAddedCharmapName][newChar] = charNodes
	return nil
}
