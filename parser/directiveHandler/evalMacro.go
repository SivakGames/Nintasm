package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/macroTable"
)

func evalMacro(directiveName string, macroLabel string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(macroLabel, enumSymbolTableTypes.Macro)
	blockStack.SetCaptureParentOpOnlyFlag()
	return nil
}

// End the macro definition and add to environment
func evalEndMacro() error {
	macroLabel := blockStack.GetLabelAndDoEndBlockSetups()
	capturedLines := blockStack.GetTopOfStackCapturedLines()
	if len(*capturedLines) == 0 {
		errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
	}

	macroTable.AddCapturedLinesToMacro(macroLabel, macroTable.Macro, *capturedLines)
	blockStack.ClearBottomOfStackCapturedLines()
	blockStack.PopFromStackAndExtendNoLines()
	return nil
}
