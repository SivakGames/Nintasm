package macroTable

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type MacroTableType = []blockStack.CapturedLine

var macroSymbolTable = map[string]MacroTableType{}
var kvMacroSymbolTable = map[string]MacroTableType{}

// The possible values for simple operations
type LookupMacroEnumType int

const (
	Macro LookupMacroEnumType = iota + 0
	KVMacro
)

// ----------------------------------
func AddIdentifierKeyToMacroTable(macroName string) {
	macroSymbolTable[macroName] = MacroTableType{}
}
func AddIdentifierKeyToKVMacroTable(macroName string) {
	kvMacroSymbolTable[macroName] = MacroTableType{}
}

// ----------------------------------

func AddCapturedLinesToMacro(macroName string, macroEnum LookupMacroEnumType, capturedLines MacroTableType) error {
	if macroEnum == Macro {
		macroSymbolTable[macroName] = capturedLines
	} else {
		kvMacroSymbolTable[macroName] = capturedLines
	}
	return nil
}

func LookupMacroInEnvironment(macroName string, macroEnum LookupMacroEnumType) (MacroTableType, bool) {
	var macro MacroTableType
	var ok bool

	if macroEnum == Macro {
		macro, ok = macroSymbolTable[macroName]
	} else {
		macro, ok = kvMacroSymbolTable[macroName]
	}
	return macro, ok

}

func LookupAndGetMacroInEnvironment(symbolName string, macroEnum LookupMacroEnumType) (MacroTableType, error) {
	macro, ok := LookupMacroInEnvironment(symbolName, macroEnum)
	if ok {
		return macro, nil
	} else {
		return nil, errorHandler.AddNew(enumErrorCodes.MacroNotExist, symbolName)
	}
}
