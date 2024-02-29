package environment

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
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

	_, exists = literalNodeSymbolTable[symbolName]
	if exists {
		return errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, "literal")
	}
	_, exists = macroTable.LookupMacroInEnvironment(symbolName, macroTable.Macro)
	if exists {
		return errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, "macro")

	}
	return nil
}

// ----------------------------------

func AddToEnvironment(symbolName string, node Node) (Node, error) {
	_, exists := literalNodeSymbolTable[symbolName]
	if exists {
		return node, errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, "literal")

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
		return operandFactory.EmptyNode(), errorHandler.AddNew(enumErrorCodes.InterpreterSymbolNotFound, symbolName)
	}
}
