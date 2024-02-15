package environment

import (
	"errors"
	"fmt"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
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
func AddMacroToEnvironment(symbolName string, capturedLines []blockStack.CapturedLine) error {
	macroSymbolTable[symbolName] = capturedLines
	return nil
}

func LookupMacroInEnvironment(symbolName string) error {
	_, ok := macroSymbolTable[symbolName]
	if ok {
		return nil
	} else {
		return errors.New("No Macro!")
	}
}
