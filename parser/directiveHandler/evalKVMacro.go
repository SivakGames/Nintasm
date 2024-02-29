package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/macroTable"
)

func evalKVMacro(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnlyFlag()
	return nil
}

// End the macro definition and add to environment
func evalEndKVMacro() error {
	macroLabel := blockStack.GetLabelAndDoEndBlockSetups()
	capturedLines := blockStack.GetTopOfStackCapturedLines()
	macroTable.AddMacroToEnvironment(macroLabel, macroTable.KVMacro, *capturedLines)

	if len(*capturedLines) == 0 {
		errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
	}

	blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendNoLines()
	return nil
}
