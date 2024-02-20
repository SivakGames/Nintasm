package environment

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
