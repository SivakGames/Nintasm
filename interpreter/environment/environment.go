package environment

import (
	"errors"
	"fmt"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

type Environment struct {
	name   string
	record map[string]Node
	parent string
}

func NewEnvironment(envName string) Environment {
	return Environment{
		name: envName,
	}
}

// ----------------------------------

func generateNumericNodeForEnvironment(number int) Node {
	return operandFactory.CreateNumericLiteralNode(enumTokenTypes.NUMBER_decimal, fmt.Sprintf("%d", number), number)
}

func generateAssemblerBuiltInFunctionNode(funcName string) Node {
	return operandFactory.CreateAssemblerBuiltInFunctionNode(funcName)
}

// ----------------------------------
func decodeHigh(node Node) error {
	if node.NodeType != enumNodeTypes.NumericLiteral {
		return errors.New("High node must be a number")
	}
	return nil //node.AsNumber
}

// ----------------------------------

var GlobalEnvironment Environment

func init() {
	GlobalEnvironment = NewEnvironment("-GLOBAL-")
	GlobalEnvironment.record = *&GlobalEnvironmentValues
	GlobalEnvironment.parent = ""
	return
}

var GlobalEnvironmentValues = map[string]Node{
	"temp":        generateNumericNodeForEnvironment(100),
	"PPUCTRL":     generateNumericNodeForEnvironment(0x02000),
	"PPUMASK":     generateNumericNodeForEnvironment(0x02001),
	"PPUADDR":     generateNumericNodeForEnvironment(0x02006),
	"PPUADDR.aba": generateNumericNodeForEnvironment(0b00000001),
	"high":        generateAssemblerBuiltInFunctionNode("high"),
	"low":         generateAssemblerBuiltInFunctionNode("low"),
}

// ----------------------------------

func AddToEnvironment(symbolName string, node Node) (Node, error) {
	_, exists := GlobalEnvironment.record[symbolName]
	if exists {
		errMsg := fmt.Sprintf("%v : SYMBOL IS ALREADY DEFINED!", symbolName)
		return node, errors.New(errMsg)

	} else {
		GlobalEnvironment.record[symbolName] = node
	}
	return GlobalEnvironment.record[symbolName], nil
}

func LookupInEnvironment(symbolName string) (Node, error) {
	return resolveInEnvironment(symbolName)
}

func resolveInEnvironment(symbolName string) (Node, error) {

	value, ok := GlobalEnvironment.record[symbolName]
	if ok {
		return value, nil
	} else {
		errMsg := fmt.Sprintf("\x1b[33m%v\x1b[0m was not found", symbolName)
		return operandFactory.EmptyNode(), errors.New(errMsg)
	}
	//return operandFactory.EmptyNode(), nil
}
