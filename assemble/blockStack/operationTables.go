package blockStack

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/util"
)

//++++++++++++++++++++++++++++++++++++++++++++++++++

// When opening a block, what to set the flags to

type startOpFlags struct {
	ForcedEval             bool
	ForcedCapture          bool
	OverwriteForcedCapture map[string]bool
}

func newStartOpFlags(forcedEval bool, forcedCapture bool, overwriteForcedCapture map[string]bool) startOpFlags {
	return startOpFlags{
		ForcedEval:             forcedEval,
		ForcedCapture:          forcedCapture,
		OverwriteForcedCapture: overwriteForcedCapture,
	}
}

var startBlockOperationFlags = map[string]startOpFlags{
	"CHARMAP":   newStartOpFlags(true, false, map[string]bool{}),
	"EXPRMAP":   newStartOpFlags(false, true, map[string]bool{}),
	"IF":        newStartOpFlags(false, false, map[string]bool{}),
	"IKV":       newStartOpFlags(true, false, map[string]bool{}),
	"IM":        newStartOpFlags(true, false, map[string]bool{}), //Invoke macro
	"KVMACRO":   newStartOpFlags(false, true, map[string]bool{}),
	"MACRO":     newStartOpFlags(false, true, map[string]bool{}),
	"NAMESPACE": newStartOpFlags(true, false, map[string]bool{}),
	"REPEAT":    newStartOpFlags(false, true, map[string]bool{"REPEAT": true}),
	"SWITCH":    newStartOpFlags(false, false, map[string]bool{}),
}

// The default values for flags when a new operation is pushed
func getStartOperationFlags(directiveName string) startOpFlags {
	flags, ok := startBlockOperationFlags[directiveName]
	if !ok {
		panic(fmt.Sprintf("Bad/undefined start op flag name: %v", directiveName))
	}

	return flags
}

//++++++++++++++++++++++++++++++++++++++++++++++++++

// How to end an opening block
var correspondingEndBlockOperations = map[string]string{
	"CHARMAP":   "ENDCHARMAP",
	"EXPRMAP":   "ENDEXPRMAP",
	"IF":        "ENDIF",
	"IKV":       "ENDIKV",
	"IM":        "ENDIM",
	"KVMACRO":   "ENDKVM",
	"MACRO":     "ENDM",
	"NAMESPACE": "ENDNAMESPACE",
	"REPEAT":    "ENDREPEAT",
	"SWITCH":    "ENDSWITCH",
}

//++++++++++++++++++++++++++++++++++++++++++++++++++

type captureableOpMap = map[enumTokenTypes.Def]bool

var sharedCapturableOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
	enumTokenTypes.DIRECTIVE_throw:      true,
	enumTokenTypes.IDENTIFIER:           true,
	enumTokenTypes.DYNAMIC_LABEL:        true,
}

var ifRepeatCapturableOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.ASSIGN_EQU:           true,
	enumTokenTypes.ASSIGN_simple:        true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_defExprMap: true,
	enumTokenTypes.DIRECTIVE_include:    true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
	enumTokenTypes.DIRECTIVE_throw:      true,
	enumTokenTypes.IDENTIFIER:           true,
	enumTokenTypes.DYNAMIC_LABEL:        true,
}

var sharedCapturableMacroOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:            true,
	enumTokenTypes.DIRECTIVE_dataBytes:    true,
	enumTokenTypes.DIRECTIVE_dataSeries:   true,
	enumTokenTypes.DIRECTIVE_mixedData:    true,
	enumTokenTypes.DIRECTIVE_blockStart:   true,
	enumTokenTypes.DIRECTIVE_blockEnd:     true,
	enumTokenTypes.DIRECTIVE_invokeKeyVal: true,
	enumTokenTypes.DIRECTIVE_throw:        true,
	enumTokenTypes.IDENTIFIER:             true,
	enumTokenTypes.DYNAMIC_LABEL:          true,
}

// Dictates which operations can be captured by the parent
// when within a block
var allowedOperationsForParentOps = map[string]captureableOpMap{
	"CHARMAP": {
		enumTokenTypes.DIRECTIVE_defCharMap: true,
	},
	"EXPRMAP": {
		enumTokenTypes.DIRECTIVE_defExprMap: true,
		enumTokenTypes.DIRECTIVE_blockStart: true,
		enumTokenTypes.DIRECTIVE_blockEnd:   true,
	},
	"IF":     ifRepeatCapturableOps,
	"ELSEIF": ifRepeatCapturableOps,
	"ELSE":   ifRepeatCapturableOps,
	"IKV": {
		enumTokenTypes.DIRECTIVE_invokeKeyVal: true,
	},
	"IM":      {},
	"KVMACRO": sharedCapturableMacroOps,
	"MACRO":   sharedCapturableMacroOps,
	"NAMESPACE": {
		enumTokenTypes.ASSIGN_simple: true,
	},
	"REPEAT": ifRepeatCapturableOps,
	"SWITCH": {
		enumTokenTypes.DIRECTIVE_blockStart: true,
	},
	"CASE":    ifRepeatCapturableOps,
	"DEFAULT": ifRepeatCapturableOps,
}

func getAllowedOperationsForCurrentBlockOperation() captureableOpMap {
	blockEntry := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	blockOperationName := blockEntry.BlockOperationName
	allowedOperations, ok := allowedOperationsForParentOps[blockOperationName]
	if !ok {
		errMsg := fmt.Sprintf("Very bad stack op! Got: %v", blockOperationName)
		panic(errMsg)
	}
	return allowedOperations
}

func checkOperationIsCapturableByCurrentBlockOperation(
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	allowedOperations := getAllowedOperationsForCurrentBlockOperation()
	_, ok := allowedOperations[lineOperationParsedValues.OperationTokenEnum]
	if !ok {
		return errorHandler.AddNew(enumErrorCodes.BlockOpUncapturableByParent, lineOperationParsedValues.OperationTokenValue)
	}

	return nil
}
