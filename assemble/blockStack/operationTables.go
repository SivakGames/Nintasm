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
	ForcedEval    bool
	ForcedCapture bool
}

func newStartOpFlags(forcedEval bool, forcedCapture bool) startOpFlags {
	return startOpFlags{
		ForcedEval:    forcedEval,
		ForcedCapture: forcedCapture,
	}
}

var startBlockOperationFlags = map[string]startOpFlags{
	"CHARMAP":   newStartOpFlags(true, false),
	"EXPRMAP":   newStartOpFlags(true, false),
	"IF":        newStartOpFlags(false, false),
	"IKV":       newStartOpFlags(true, false),
	"IM":        newStartOpFlags(true, false), //Invoke macro
	"KVMACRO":   newStartOpFlags(false, true),
	"MACRO":     newStartOpFlags(false, true),
	"NAMESPACE": newStartOpFlags(true, false),
	"REPEAT":    newStartOpFlags(false, false),
	"SWITCH":    newStartOpFlags(false, false),
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
	enumTokenTypes.TEMPLATE_STRING:      true,
}

var sharedCapturableMacroOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
	enumTokenTypes.DIRECTIVE_blockStart: true,
	enumTokenTypes.DIRECTIVE_blockEnd:   true,
	enumTokenTypes.DIRECTIVE_throw:      true,
	enumTokenTypes.IDENTIFIER:           true,
	enumTokenTypes.TEMPLATE_STRING:      true,
}

// Dictates which operations can be captured by the parent
// when within a block
var allowedOperationsForParentOps = map[string]captureableOpMap{
	"CHARMAP": {
		enumTokenTypes.DIRECTIVE_defCharMap: true,
	},
	"EXPRMAP": {
		enumTokenTypes.DIRECTIVE_defExprMap: true,
	},
	"IF":     sharedCapturableOps,
	"ELSEIF": sharedCapturableOps,
	"ELSE":   sharedCapturableOps,
	"IKV": {
		enumTokenTypes.DIRECTIVE_invokeKeyVal: true,
	},
	"IM":      {},
	"KVMACRO": sharedCapturableMacroOps,
	"MACRO":   sharedCapturableMacroOps,
	"NAMESPACE": {
		enumTokenTypes.ASSIGN_simple: true,
	},
	"REPEAT": sharedCapturableOps,
	"SWITCH": {
		enumTokenTypes.DIRECTIVE_blockStart: true,
	},
	"CASE":    sharedCapturableOps,
	"DEFAULT": sharedCapturableOps,
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
