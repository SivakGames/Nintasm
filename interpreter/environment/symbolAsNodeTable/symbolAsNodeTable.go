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

var prevLabelNextLabelTable = map[string]string{}
var prevLocalLabelNextLocalLabelTable = map[string]string{}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type SymbolTableType = map[string]Node

// When in block ops, certain things can be scoped to the block
var localBlockScopes []*SymbolTableType

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

// Used when doing function calls
var symbolTableStack = []SymbolTableType{}

var symbolTable = SymbolTableType{
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

	"PPUSTATUS": generateNumericNodeForEnvironment(0x02002),
	"PPUSCROLL": generateNumericNodeForEnvironment(0x02005),

	"PPUADDR": generateNumericNodeForEnvironment(0x02006),
	"PPUDATA": generateNumericNodeForEnvironment(0x02007),

	"bank": generateAssemblerReservedWordNode("bank"),
	"high": generateAssemblerReservedWordNode("high"),
	"low":  generateAssemblerReservedWordNode("low"),
}

// -----------------------------------------
func init() {
	generateKeys := func(baseKey string, startValue int, endValue, step int) {
		for value := startValue; value < endValue; value += step {
			key := fmt.Sprintf("%s%d", baseKey, (value-startValue)/step)
			genNode := generateNumericNodeForEnvironment(value)
			symbolTable[key] = genNode
		}
	}

	for i := 0; i < 4; i++ {
		start := i*0x00400 + 0x02000
		end := start + 0x003c0
		name := fmt.Sprintf("PPUADDR.nt%dline", i)
		generateKeys(name, start, end, 0x20)
	}
	for i := 0; i < 4; i++ {
		start := i*0x00400 + 0x023c0
		end := start + 0x00040
		name := fmt.Sprintf("PPUADDR.nt%dattLine", i)
		generateKeys(name, start, end, 0x08)
	}
	for i := 0; i < 4; i++ {
		value := i*0x00004 + 0x03f00
		name := fmt.Sprintf("PPUADDR.palBg%d", i)
		symbolTable[name] = generateNumericNodeForEnvironment(value)
	}
	for i := 0; i < 4; i++ {
		value := i*0x00004 + 0x03f10
		name := fmt.Sprintf("PPUADDR.palSprite%d", i)
		symbolTable[name] = generateNumericNodeForEnvironment(value)
	}
}

// -----------------------------------------

func AddIdentifierKeyToSymbolAsNodeTable(symbolName string, node Node) {
	symbolTable[symbolName] = node
}
func AddIdentifierKeyToLabelAsBankTable(symbolName string, bankId int) {
	labalAsBankTable[symbolName] = bankId
}
func AddIdentifierKeyToPrevLabelNextLabelTable(prevLabel string, nextLabel string) {
	prevLabelNextLabelTable[prevLabel] = nextLabel
}
func AddIdentifierKeyToPrevLocalLabelNextLocalLabelTable(prevLabel string, nextLabel string) {
	prevLocalLabelNextLocalLabelTable[prevLabel] = nextLabel
}

// -----------------------------------------

// Go through the symbol as node table and check sub-scopes before checking the global scope
func GetNodeFromSymbolAsNodeTable(symbolName string) (Node, bool) {
	for lsi := len(localBlockScopes) - 1; lsi >= 0; lsi-- {
		scope := localBlockScopes[lsi]
		node, exists := (*scope)[symbolName]
		if exists {
			return node, exists
		}
	}

	node, exists := symbolTable[symbolName]
	return node, exists
}

// -----------------------------------------

func GetValueFromLabelAsBankTable(symbolName string) (int, bool) {
	bankValue, exists := labalAsBankTable[symbolName]
	return bankValue, exists
}
func GetValueFromPrevLabelNextLabelTable(symbolName string) (string, bool) {
	nextLabelName, exists := prevLabelNextLabelTable[symbolName]
	return nextLabelName, exists
}
func GetValueFromPrevLocalLabelNextLocalLabelTable(symbolName string) (string, bool) {
	nextLocalLabelName, exists := prevLocalLabelNextLocalLabelTable[symbolName]
	return nextLocalLabelName, exists
}

// ------------------------------------------

func PushToSymbolTableStack() {
	symbolTableStack = append(symbolTableStack, SymbolTableType{})
}

func PopFromSymbolTableStack() {
	symbolTableStack = symbolTableStack[:len(symbolTableStack)-1]
}

func CheckTopOfSymbolTableStackHasLength() bool {
	return len(symbolTableStack) > 0
}

func GetTopOfSymbolTableStack() *SymbolTableType {
	if len(symbolTableStack) > 0 {
		return &symbolTableStack[len(symbolTableStack)-1]
	}
	return nil
}

func AddSymbolToTopTableStack(symbolName string, node Node) {
	topStack := GetTopOfSymbolTableStack()
	if topStack != nil {
		(*topStack)[symbolName] = node
		return
	}
	errorHandler.AddNew(enumErrorCodes.Other, "Function stack is empty!")
	return
}

// Checks if there is anything on top of the stack and if the symbol exists.
func LookupSymbolInTopOfSymbolTableStack(symbolName string) (Node, bool) {
	topStack := GetTopOfSymbolTableStack()
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
	return operandFactory.CreateNumericLiteralNode(float64(number))
}

// ------------------------------------------

func DeepCopyLocalBlockScopes() []*SymbolTableType {
	copyScopes := make([]*SymbolTableType, len(localBlockScopes))
	for i, scope := range localBlockScopes {
		copyScope := make(SymbolTableType)
		for key, value := range *scope {
			copyScope[key] = value
		}
		copyScopes[i] = &copyScope
	}
	return copyScopes
}

func AddChildBlockScope(scope SymbolTableType) {
	localBlockScopes = append(localBlockScopes, &scope)
}

func PopChildBlockScope() {
	localBlockScopes = localBlockScopes[:len(localBlockScopes)-1]
}

func SetCurrentLocalBlockScopes(blockScopes []*SymbolTableType) {
	localBlockScopes = blockScopes
}

func ClearCurrentLocalBlockScopes() {
	localBlockScopes = localBlockScopes[:0]
}
