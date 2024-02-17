package environment

import (
	enumParserTypes "misc/nintasm/enums/parserTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
)

type MacroTableType = []blockStack.CapturedLine

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
