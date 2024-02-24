package directiveHandler

import (
	"errors"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/operandFactory"
)

var ikvKeys map[string]string

func evalIkv(directiveName string, operandList *[]Node) error {
	macroNameNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsIdentifier(macroNameNode) {
		return errors.New("IKV name must be an Identifier!")
	}

	blockStack.PushOntoStack(directiveName, *operandList)
	blockStack.SetCaptureParentOpOnlyFlag()
	blockStack.SetCurrentOperationEvaluatesFlag()

	macroTable.AppendToReplacementStack()

	return nil
}

func evalEndIkv(operandList *[]Node) error {
	currentStackOperation := blockStack.GetTopOfStackOperation()
	currentStackOperationOperandList := &currentStackOperation.OperandList
	macroNameNode := &(*currentStackOperationOperandList)[0]
	macroData, err := macroTable.LookupAndGetMacroInEnvironment(macroNameNode.NodeValue, macroTable.KVMacro)
	if err != nil {
		return err
	}

	replacementList := macroTable.GetReplacementListOnTopOfStack()
	var modifiedCapturedLines []blockStack.CapturedLine

	for _, md := range macroData {
		for _, replacementListItem := range *replacementList {
			md.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllString(md.OriginalLine, replacementListItem.ReplaceString)
		}
		modifiedCapturedLines = append(modifiedCapturedLines, md)
	}

	macroTable.PopFromReplacementStack()

	blockStack.ClearCaptureParentOpOnlyFlag()
	blockStack.ClearCurrentOperationEvaluatesFlag()
	blockStack.PopFromStackAndExtendCapturedLines(modifiedCapturedLines)

	return nil
}

func evalKv(operandList *[]Node) error {
	macroKeyNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsSubstitutionID(macroKeyNode) {
		return errors.New("Must use a substitution type node for KV")
	}
	macroValueNode := &(*operandList)[1]
	macroTable.AddToReplacementListOnTopOfStack(macroKeyNode.NodeValue, macroValueNode.NodeValue)
	return nil
}
