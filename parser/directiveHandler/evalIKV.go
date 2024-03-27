package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
)

var ikvKeys map[string]string

func evalIkv(directiveName string, operandList *[]Node) error {
	macroNameNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(macroNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotIdentifier) // ❌ Fails
	}

	blockStack.PushCaptureBlock(directiveName, *operandList)
	macroTable.AppendToReplacementStack()
	symbolAsNodeTable.PushToSymbolTableStack()

	return nil
}

// Middle operation - adding a key to IKV
func evalKv(operandList *[]Node) error {
	macroKeyNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsSubstitutionID(macroKeyNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
	}
	macroValueNode := &(*operandList)[1]
	macroTable.AddToReplacementListOnTopOfStack(macroKeyNode.NodeValue, macroValueNode.NodeValue)
	symbolAsNodeTable.AddSymbolToTopTableStack(macroKeyNode.NodeValue, operandFactory.CreateBooleanLiteralNode(true))

	return nil
}

// Final operation
func evalEndIkv(operandList *[]Node) error {
	// Get invoking macro name
	_, originalOperandList := blockStack.GetCurrentCaptureBlockCapturedLinesAndOperandList()
	macroNameNode := &(*originalOperandList)[0]

	// Get macro's data
	macroData, err := macroTable.LookupAndGetMacroInEnvironment(macroNameNode.NodeValue, macroTable.KVMacro)
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

	macroTable.PopFromReplacementStack()
	blockStack.ClearCaptureBlockListEvalFlag()
	blockStack.PopCaptureBlockThenExtendCapturedLines(modifiedCapturedLines)
	return nil
}
