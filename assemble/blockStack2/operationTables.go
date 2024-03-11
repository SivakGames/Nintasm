package blockStack2

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
	"KVMACRO":   newStartOpFlags(false, true),
	"MACRO":     newStartOpFlags(false, true),
	"NAMESPACE": newStartOpFlags(true, false),
	"REPEAT":    newStartOpFlags(false, false),
}

func getStartOperationFlags(directiveName string) startOpFlags {
	flags, ok := startBlockOperationFlags[directiveName]
	if !ok {
		panic("Bad/undefined start op flag name")
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
	"KVMACRO":   "ENDKVM",
	"MACRO":     "ENDM",
	"NAMESPACE": "ENDNAMESPACE",
	"REPEAT":    "ENDREPEAT",
}

//++++++++++++++++++++++++++++++++++++++++++++++++++

type captureableOpMap = map[enumTokenTypes.Def]bool

var sharedCapturableOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
}

var sharedCapturableMacroOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
	enumTokenTypes.DIRECTIVE_blockStart: true,
	enumTokenTypes.DIRECTIVE_blockEnd:   true,
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
	"KVMACRO": sharedCapturableMacroOps,
	"MACRO":   sharedCapturableMacroOps,
	"NAMESPACE": {
		enumTokenTypes.ASSIGN_simple: true,
	},
	"REPEAT": sharedCapturableOps,
}

func getAllowedOperationsForCurrentBlockOperation() captureableOpMap {
	blockEntry := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
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
