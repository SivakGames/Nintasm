package symbolAsNodeTable

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

type Node = operandFactory.Node

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

var labalAsBankTable = map[string]int{}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type symbolTableType = map[string]Node

// Used when doing function calls
var symbolTableStack = []symbolTableType{}

var symbolTable = symbolTableType{
	"CTRLBTN.right":  generateNumericNodeForEnvironment(0x01),
	"CTRLBTN.left":   generateNumericNodeForEnvironment(0x02),
	"CTRLBTN.down":   generateNumericNodeForEnvironment(0x04),
	"CTRLBTN.up":     generateNumericNodeForEnvironment(0x08),
	"CTRLBTN.start":  generateNumericNodeForEnvironment(0x10),
	"CTRLBTN.select": generateNumericNodeForEnvironment(0x20),
	"CTRLBTN.b":      generateNumericNodeForEnvironment(0x40),
	"CTRLBTN.a":      generateNumericNodeForEnvironment(0x80),

	"PPUCTRL":                    generateNumericNodeForEnvironment(0x02000),
	"PPUCTRL.nameTable0":         generateNumericNodeForEnvironment(0x00),
	"PPUCTRL.nameTable1":         generateNumericNodeForEnvironment(0x01),
	"PPUCTRL.nameTable2":         generateNumericNodeForEnvironment(0x02),
	"PPUCTRL.nameTable3":         generateNumericNodeForEnvironment(0x03),
	"PPUCTRL.drawDirection":      generateNumericNodeForEnvironment(0x04),
	"PPUCTRL.spritePatternTable": generateNumericNodeForEnvironment(0x08),
	"PPUCTRL.bgPatternTable":     generateNumericNodeForEnvironment(0x10),
	"PPUCTRL.use8x16Sprites":     generateNumericNodeForEnvironment(0x20),
	"PPUCTRL.masterSlave":        generateNumericNodeForEnvironment(0x40),
	"PPUCTRL.enableNMI":          generateNumericNodeForEnvironment(0x80),

	"PPUMASK":                   generateNumericNodeForEnvironment(0x02001),
	"PPUMASK.grayscale":         generateNumericNodeForEnvironment(0x01),
	"PPUMASK.disableBgClip":     generateNumericNodeForEnvironment(0x02),
	"PPUMASK.disableSpriteClip": generateNumericNodeForEnvironment(0x04),
	"PPUMASK.showBg":            generateNumericNodeForEnvironment(0x08),
	"PPUMASK.showSprites":       generateNumericNodeForEnvironment(0x10),
	"PPUMASK.emphasizeRed":      generateNumericNodeForEnvironment(0x20),
	"PPUMASK.emphasizeGreen":    generateNumericNodeForEnvironment(0x40),
	"PPUMASK.emphasizeBlue":     generateNumericNodeForEnvironment(0x80),

	"PPUADDR": generateNumericNodeForEnvironment(0x02006),

	"bank": generateAssemblerReservedWordNode("bank"),
	"high": generateAssemblerReservedWordNode("high"),
	"low":  generateAssemblerReservedWordNode("low"),
}

// -----------------------------------------
func init() {
	generateKeys := func(baseKey string, startValue int, endValue, step int) {
		for value := startValue; value < endValue; value += step {
			key := fmt.Sprintf("%s%d", baseKey, (value-startValue)/step)
			symbolTable[key] = generateNumericNodeForEnvironment(value)
		}
	}

	for i := 0; i < 4; i++ {
		start := i*0x00400 + 0x02000
		end := start + 0x003c0
		name := fmt.Sprintf("PPUADDR.nt%dline", i)
		generateKeys(name, start, end, 0x20)
	}
	for i := 0; i < 4; i++ {
		start := i*0x00400 + 0x03f00
		end := start + 0x00004
		name := fmt.Sprintf("PPUADDR.palBg%d", i)
		generateKeys(name, start, end, 0x04)
	}
}

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

// ------------------------------------------

func PushToSymbolTableStack() {
	symbolTableStack = append(symbolTableStack, symbolTableType{})
}

func PopFromSymbolTableStack() {
	symbolTableStack = symbolTableStack[:len(symbolTableStack)-1]
}

func CheckTopOfSymbolTableStack() *symbolTableType {
	if len(symbolTableStack) > 0 {
		return &symbolTableStack[len(symbolTableStack)-1]
	}
	return nil
}

func AddSymbolToTopTableStack(symbolName string, node Node) {
	topStack := CheckTopOfSymbolTableStack()
	if topStack != nil {
		(*topStack)[symbolName] = node
		return
	}
	errorHandler.AddNew(enumErrorCodes.Other, "Function stack is empty!")
	return
}

func LookupTopOfSymbolTableStack(symbolName string) (Node, bool) {
	topStack := CheckTopOfSymbolTableStack()
	if topStack != nil {
		node, exists := (*topStack)[symbolName]
		return node, exists
	}
	return operandFactory.ErrorNode(symbolName), false
}

// +++++++++++++++++++++++++++++++++++++++++

func generateAssemblerReservedWordNode(funcName string) Node {
	return operandFactory.CreateAssemblerReservedWordNode(funcName)
}

func generateNumericNodeForEnvironment(number int) Node {
	return operandFactory.CreateNumericLiteralNode(number)
}
