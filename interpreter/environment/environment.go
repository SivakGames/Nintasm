package environment

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/environment/macroTable"
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
	_, exists = macroTable.LookupMacroInEnvironment(symbolName, macroTable.Macro)
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
