package environment

import (
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

var literalNodeSymbolTable = map[string]Node{
	"temp":        generateNumericNodeForEnvironment(100),
	"PPUCTRL":     generateNumericNodeForEnvironment(0x02000),
	"PPUMASK":     generateNumericNodeForEnvironment(0x02001),
	"PPUADDR":     generateNumericNodeForEnvironment(0x02006),
	"PPUADDR.aba": generateNumericNodeForEnvironment(0b00000001),
	"bank":        generateAssemblerReservedWordNode("bank"),
	"high":        generateAssemblerReservedWordNode("high"),
	"low":         generateAssemblerReservedWordNode("low"),
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type MacroTableType = []blockStack.CapturedLine

var macroSymbolTable = map[string]MacroTableType{
	"__PPU__": {{
		OriginalLine:         " LDA #$20",
		OperationLabel:       "",
		OperationTokenEnum:   enumTokenTypes.INSTRUCTION,
		OperationTokenValue:  "LDA",
		OperandStartPosition: 5,
		ParentParserEnum:     enumParserTypes.Instruction,
	}},
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type CharmapTableType = map[rune][]Node

var charmapSymbolTable = map[string]CharmapTableType{}
