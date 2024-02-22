package macroTable

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/blockStack"
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
		errMsg := fmt.Sprintf("Macro %v does not exist!", symbolName)
		return nil, errors.New(errMsg)
	}
}
