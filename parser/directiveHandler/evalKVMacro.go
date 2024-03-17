package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/macroTable"
)

func evalKVMacro(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushCaptureBlock(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(macroLabel, enumSymbolTableTypes.KVMacro)
	return nil
}

// End the macro definition and add to environment
func evalEndKVMacro() error {
	macroLabel := blockStack.GetCurrentOperationLabel()
	capturedLines := blockStack.GetCurrentCaptureBlockCapturedLines()
	if len(*capturedLines) == 0 {
		errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
	}

	macroTable.AddCapturedLinesToMacro(macroLabel, macroTable.KVMacro, *capturedLines)
	blockStack.ProcessEndLabeledDirective()
	return nil
}
