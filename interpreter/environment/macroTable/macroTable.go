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

// ----------------------------------
func AddIdentifierKeyToMacroTable(macroName string) {
	macroSymbolTable[macroName] = MacroXX{}
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func RemoveIdentifierKeyFromMacroTable(macroName string) {
	delete(macroSymbolTable, macroName)
}

// ----------------------------------

func AddCapturedLinesToMacro(macroName string, capturedLines MacroTableType) {
	macro := macroSymbolTable[macroName]
	macro.Lines = capturedLines
	macroSymbolTable[macroName] = macro

}

func AddArgumentsToMacroTable(macroName string, arguments *[]string) {
	macro := macroSymbolTable[macroName]
	macro.Arguments = arguments
	macroSymbolTable[macroName] = macro
}

// ----------------------------------

func LookupMacroInEnvironment(macroName string) (MacroXX, bool) {
	var macro MacroXX
	var ok bool

	macro, ok = macroSymbolTable[macroName]

	return macro, ok
}

func LookupAndGetMacroInEnvironment(symbolName string) (MacroTableType, *[]string, error) {
	macro, ok := LookupMacroInEnvironment(symbolName)
	if ok {
		return macro.Lines, macro.Arguments, nil
	} else {
		return nil, nil, errorHandler.AddNew(enumErrorCodes.MacroNotExist, symbolName)
	}
}
