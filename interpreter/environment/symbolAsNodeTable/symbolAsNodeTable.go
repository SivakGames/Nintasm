package symbolAsNodeTable

import (
	"misc/nintasm/interpreter/operandFactory"
)

type Node = operandFactory.Node

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

var symbolTable = map[string]Node{
	"temp":        generateNumericNodeForEnvironment(100),
	"PPUCTRL":     generateNumericNodeForEnvironment(0x02000),
	"PPUMASK":     generateNumericNodeForEnvironment(0x02001),
	"PPUADDR":     generateNumericNodeForEnvironment(0x02006),
	"PPUADDR.aba": generateNumericNodeForEnvironment(0b00000001),
	"bank":        generateAssemblerReservedWordNode("bank"),
	"high":        generateAssemblerReservedWordNode("high"),
	"low":         generateAssemblerReservedWordNode("low"),
}

var labalAsBankTable = map[string]int{}

// -----------------------------------------

func AddIdentifierKeyToSymbolAsNodeTable(symbolName string, node Node) {
	symbolTable[symbolName] = node
}

func AddIdentifierKeyToLabelAsBankTable(symbolName string, bankId int) {
	labalAsBankTable[symbolName] = bankId
}

func GetNodeFromSymbolAsNodeTable(symbolName string) (Node, bool) {
	node, exists := symbolTable[symbolName]
	return node, exists
}

// +++++++++++++++++++++++++++++++++++++++++

func generateAssemblerReservedWordNode(funcName string) Node {
	return operandFactory.CreateAssemblerReservedWordNode(funcName)
}

func generateNumericNodeForEnvironment(number int) Node {
	return operandFactory.CreateNumericLiteralNode(number)
}
