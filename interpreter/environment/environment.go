package environment

import (
	"fmt"
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
}

// ----------------------------------

func LookupInEnvironment(symbolName string) Node {
	return resolveInEnvironment(symbolName)
}

func resolveInEnvironment(symbolName string) Node {

	value, ok := GlobalEnvironment.record[symbolName]
	if ok {
		return value
	} else {
		fmt.Println("Not found")
	}
	return operandFactory.EmptyNode()
}
