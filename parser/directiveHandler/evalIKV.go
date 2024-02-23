package directiveHandler

import (
	"errors"
	"fmt"
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

	return nil
}

func evalEndIkv(directiveName string, operandList *[]Node) error {
	currentStackOperation := blockStack.GetTopOfStackOperation()
	currentStackOperationOperandList := &currentStackOperation.OperandList
	macroNameNode := &(*currentStackOperationOperandList)[0]
	macroData, err := macroTable.LookupAndGetMacroInEnvironment(macroNameNode.NodeValue, macroTable.KVMacro)
	if err != nil {
		return err
	}

	fmt.Println(macroData)

	blockStack.ClearCaptureParentOpOnlyFlag()
	blockStack.ClearCurrentOperationEvaluatesFlag()
	blockStack.PopFromStackAndExtendNoLines()
	//blockStack.PopFromStackAndExtendCapturedLines(*currentStackOperationCapturedLines)

	return nil
}

func evalKv(directiveName string, operandList *[]Node) error {
	macroKeyNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsSubstitutionID(macroKeyNode) {
		return errors.New("Must use a substitution type node for KV")
	}
	macroValueNode := &(*operandList)[1]
	_ = macroValueNode

	//ikvKeys[macroKeyNode.NodeValue] = ""

	fmt.Println("macroData", operandList)
	return nil
}
