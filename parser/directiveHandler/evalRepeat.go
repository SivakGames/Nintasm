package directiveHandler

import (
	"misc/nintasm/assemble/blockStack2"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"regexp"
	"strconv"
)

func evalRepeat(directiveName string, operandList *[]Node) error {

	// Check and validate repeat amount
	numRepeatsNode, err := interpreter.EvaluateNode((*operandList)[0])
	if err != nil {
		return err
	}
	if !operandFactory.ValidateNodeIsNumeric(&numRepeatsNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIsGTEValue(&numRepeatsNode, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
	}

	evaluatedNodes := []Node{numRepeatsNode}

	// Check if iterator node was set
	if len(*operandList) > 1 {
		iterNameNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsSubstitutionID(&iterNameNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[1])
	}

	//See if this is the bottom entry (i.e. creating a new stack)
	// Set the eval instead of capture to false

	//blockStack.PushOntoStack(directiveName, evaluatedNodes)
	//blockStack.ClearCaptureParentOpOnlyFlag()
	//blockStack.ClearCurrentOperationEvaluatesCapturedNodesFlag()
	blockStack2.PushOntoTopEntry(directiveName, evaluatedNodes)

	return nil
}

func evalEndRepeat() error {
	capturedLines, operandList := blockStack2.GetTopBlockEntryData()
	//currentStackOperation := blockStack.GetTopOfStackOperation()
	//operandList := &currentStackOperation.OperandList
	//capturedLines := &currentStackOperation.CapturedLines

	//Extract repeatAmount
	repeatAmount := (*operandList)[0].AsNumber
	iteratorName := ""

	if len(*operandList) > 1 {
		iteratorName = (*operandList)[1].NodeValue
	}

	replacedLines := make([]blockStack2.CapturedLine, len(*capturedLines)*repeatAmount)
	replacedIndex := 0

	if iteratorName == "" {
		for i := 0; i < repeatAmount; i++ {
			for _, j := range *capturedLines {
				replacedLines[replacedIndex] = j
				replacedIndex++
			}
		}
	} else {
		iterNameAsRegex := regexp.MustCompile(`\` + iteratorName + `\b`)
		for i := 0; i < repeatAmount; i++ {
			replaceNum := strconv.Itoa(i)
			for _, j := range *capturedLines {
				replaced := iterNameAsRegex.ReplaceAllString(j.OriginalLine, replaceNum)
				j.OriginalLine = replaced
				replacedLines[replacedIndex] = j
				replacedIndex++
			}
		}
	}
	if len(replacedLines) == 0 {
		errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
	}

	//blockStack.PopFromStackAndExtendCapturedLines(replacedLines)
	blockStack2.PopTopEntryThenExtendCapturedLines(replacedLines)
	return nil
}
