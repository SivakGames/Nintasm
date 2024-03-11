package directiveHandler

import (
	"misc/nintasm/assemble/blockStack2"
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

	blockStack2.PushOntoTopEntry(directiveName, *operandList)
	//blockStack.PushOntoStack(directiveName, *operandList)
	//blockStack.SetCaptureParentOpOnlyFlag()
	//blockStack.SetCurrentOperationEvaluatesCapturedNodesFlag()

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
	//currentStackOperation := blockStack.GetTopOfStackOperation()
	//currentStackOperationOperandList := &currentStackOperation.OperandList
	//macroNameNode := &(*currentStackOperationOperandList)[0]

	// Get invoking macro name
	_, originalOperandList := blockStack2.GetTopBlockEntryData()
	macroNameNode := &(*originalOperandList)[0]

	// Get macro's data
	macroData, err := macroTable.LookupAndGetMacroInEnvironment(macroNameNode.NodeValue, macroTable.KVMacro)
	if err != nil {
		return err
	}

	//var modifiedCapturedLines []blockStack.CapturedLine

	var modifiedCapturedLines []blockStack2.CapturedLine
	replacementList := macroTable.GetReplacementListOnTopOfStack()

	//Replace with substitutions
	for _, md := range macroData {
		for _, replacementListItem := range *replacementList {
			md.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllString(md.OriginalLine, replacementListItem.ReplaceString)
		}
		modifiedCapturedLines = append(modifiedCapturedLines, md)
	}

	macroTable.PopFromReplacementStack()

	//blockStack.ClearCaptureParentOpOnlyFlag()

	//Clearing this flag was a reset from the old method.
	//New method should just leave the flag alone at the end op and all opening ops
	// should set/clear it instead

	//blockStack.ClearCurrentOperationEvaluatesCapturedNodesFlag()
	blockStack2.ClearCurrentInvokeOperationEvalFlag()
	//blockStack.PopFromStackAndExtendCapturedLines(modifiedCapturedLines)
	blockStack2.PopTopEntryThenExtendCapturedLines(modifiedCapturedLines)

	return nil
}
