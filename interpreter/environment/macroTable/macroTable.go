package macroTable

import (
	"errors"
	"fmt"
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
)

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

// ----------------------------------

func AddMacroToEnvironment(symbolName string, capturedLines MacroTableType) error {
	macroSymbolTable[symbolName] = capturedLines
	return nil
}

func LookupMacroInEnvironment(symbolName string) (MacroTableType, bool) {
	macro, ok := macroSymbolTable[symbolName]
	return macro, ok
}

func LookupAndGetMacroInEnvironment(symbolName string) (MacroTableType, error) {
	macro, ok := LookupMacroInEnvironment(symbolName)
	if ok {
		return macro, nil
	} else {
		errMsg := fmt.Sprintf("Macro %v does not exist!", symbolName)
		return nil, errors.New(errMsg)
	}
}
