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
type LookupMacroType int

const (
	Macro LookupMacroType = iota + 0
	KVMacro
)

// ----------------------------------

func AddMacroToEnvironment(symbolName string, macroEnum LookupMacroType, capturedLines MacroTableType) error {
	if macroEnum == Macro {
		macroSymbolTable[symbolName] = capturedLines
	} else {
		kvMacroSymbolTable[symbolName] = capturedLines
	}
	return nil
}

func LookupMacroInEnvironment(symbolName string, macroEnum LookupMacroType) (MacroTableType, bool) {
	var macro MacroTableType
	var ok bool

	if macroEnum == Macro {
		macro, ok = macroSymbolTable[symbolName]
	} else {
		macro, ok = kvMacroSymbolTable[symbolName]
	}
	return macro, ok

}

func LookupAndGetMacroInEnvironment(symbolName string, macroEnum LookupMacroType) (MacroTableType, error) {
	macro, ok := LookupMacroInEnvironment(symbolName, macroEnum)
	if ok {
		return macro, nil
	} else {
		return nil, errorHandler.AddNew(enumErrorCodes.InterpreterSymbolNotFound, symbolName)
	}
}
