package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
)

var ikvKeys = map[string]bool{}

func evalIkv(directiveName string, operandList *[]Node) error {
	macroNameNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(macroNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
	}

	blockStack.PushCaptureBlock(directiveName, *operandList)
	macroTable.AppendToReplacementStack()
	symbolAsNodeTable.PushToSymbolTableStack()
	blockStack.SetInlineEval()

	return nil
}

// Middle operation - adding a key to IKV
func evalKv(operandList *[]Node) error {
	macroKeyNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsSubstitutionID(macroKeyNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
	}
	macroValueNode := &(*operandList)[1]

	_, exists := ikvKeys[macroKeyNode.NodeValue]
	if exists {
		return errorHandler.AddNew(enumErrorCodes.Other, "Key already exists")
	}
	ikvKeys[macroKeyNode.NodeValue] = true

	macroTable.AddToReplacementListOnTopOfStack(macroKeyNode.NodeValue, macroValueNode.NodeValue)
	symbolAsNodeTable.AddSymbolToTopTableStack(macroKeyNode.NodeValue, operandFactory.CreateBooleanLiteralNode(true))

	return nil
}

// Final operation
func evalEndIkv() error {
	// Get invoking macro name
	_, originalOperandList := blockStack.GetCurrentCaptureBlockCapturedLinesAndOperandList()
	macroNameNode := &(*originalOperandList)[0]
	ikvKeys = map[string]bool{}

	// Get macro's data
	macroData, _, err := macroTable.LookupAndGetMacroInEnvironment(macroNameNode.NodeValue)
	if err != nil {
		return err
	}

	var modifiedCapturedLines []blockStack.CapturedLine
	replacementList := macroTable.GetReplacementListOnTopOfStack()

	//Replace with substitutions
	for _, md := range macroData {
		for _, replacementListItem := range *replacementList {
			md.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllLiteralString(md.OriginalLine, replacementListItem.ReplaceString)
			md.OperationLabel = replacementListItem.ReplaceRegex.ReplaceAllLiteralString(md.OperationLabel, replacementListItem.ReplaceString)
		}
		modifiedCapturedLines = append(modifiedCapturedLines, md)
	}

	//macroTable.PopFromReplacementStack()

	blockStack.CopyPresetCapturedLinesToProcessedWithEmptyScope(&modifiedCapturedLines)

	//blockStack.ClearCaptureBlockListEvalFlag()
	//blockStack.PopCaptureBlockThenExtendCapturedLines(modifiedCapturedLines)
	return nil
}

func buildIkvMacro() {
	macroTable.PopFromReplacementStack()
}
