package macroTable

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type MacroTableType = []blockStack.CapturedLine

type MacroXX struct {
	Lines     MacroTableType
	Arguments *[]string
}

var macroSymbolTable = map[string]MacroXX{}
var kvMacroSymbolTable = map[string]MacroXX{}

// The possible values for simple operations
type LookupMacroEnumType int

const (
	Macro LookupMacroEnumType = iota + 0
	KVMacro
)

// ----------------------------------
func AddIdentifierKeyToMacroTable(macroName string) {
	macroSymbolTable[macroName] = MacroXX{}
}
func AddIdentifierKeyToKVMacroTable(macroName string) {
	kvMacroSymbolTable[macroName] = MacroXX{}
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func RemoveIdentifierKeyFromMacroTable(macroName string) {
	delete(macroSymbolTable, macroName)
}
func RemoveIdentifierKeyFromKVMacroTable(macroName string) {
	delete(kvMacroSymbolTable, macroName)
}

// ----------------------------------

func AddCapturedLinesToMacro(macroName string, macroEnum LookupMacroEnumType, capturedLines MacroTableType) {
	if macroEnum == Macro {
		macro := macroSymbolTable[macroName]
		macro.Lines = capturedLines
		macroSymbolTable[macroName] = macro
	} else {
		macro := kvMacroSymbolTable[macroName]
		macro.Lines = capturedLines
		kvMacroSymbolTable[macroName] = macro
	}
}

func AddArgumentsToMacroTable(macroName string, arguments *[]string) {
	macro := macroSymbolTable[macroName]
	macro.Arguments = arguments
	macroSymbolTable[macroName] = macro
}

// ----------------------------------

func LookupMacroInEnvironment(macroName string, macroEnum LookupMacroEnumType) (MacroXX, bool) {
	var macro MacroXX
	var ok bool

	if macroEnum == Macro {
		macro, ok = macroSymbolTable[macroName]
	} else {
		macro, ok = kvMacroSymbolTable[macroName]
	}
	return macro, ok
}

func LookupAndGetMacroInEnvironment(symbolName string, macroEnum LookupMacroEnumType) (MacroTableType, *[]string, error) {
	macro, ok := LookupMacroInEnvironment(symbolName, macroEnum)
	if ok {
		return macro.Lines, macro.Arguments, nil
	} else {
		return nil, nil, errorHandler.AddNew(enumErrorCodes.MacroNotExist, symbolName)
	}
}
