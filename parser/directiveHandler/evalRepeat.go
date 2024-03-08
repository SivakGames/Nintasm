package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"regexp"
	"strconv"
)

func evalRepeat(directiveName string, operandList *[]Node) error {
	numRepeatsNode, err := interpreter.EvaluateNode((*operandList)[0])
	if err != nil {
		return err
	}

	blockStack.SetCurrentOperationEvaluatesCapturedNodesFlag()

	if !operandFactory.ValidateNodeIsNumeric(&numRepeatsNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIsGTEValue(&numRepeatsNode, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
	}

	evaluatedNodes := []Node{numRepeatsNode}

	if len(*operandList) > 1 {
		iterNameNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsSubstitutionID(&iterNameNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[1])
	}

	blockStack.PushOntoStack(directiveName, evaluatedNodes)

	return nil
}

func evalEndRepeat(operandList *[]Node) error {
	currentStackOperation := blockStack.GetTopOfStackOperation()
	currentStackOperationOperandList := &currentStackOperation.OperandList
	currentStackOperationCapturedLines := &currentStackOperation.CapturedLines

	repeatAmount := (*currentStackOperationOperandList)[0].AsNumber
	iteratorName := ""
	if len(*currentStackOperationOperandList) > 1 {
		iteratorName = (*currentStackOperationOperandList)[1].NodeValue
	}

	replacedLines := make([]blockStack.CapturedLine, len(*currentStackOperationCapturedLines)*repeatAmount)
	replacedIndex := 0

	if iteratorName == "" {
		for i := 0; i < repeatAmount; i++ {
			for _, j := range *currentStackOperationCapturedLines {
				replacedLines[replacedIndex] = j
				replacedIndex++
			}
		}
	} else {
		iterNameAsRegex := regexp.MustCompile(`\` + iteratorName + `\b`)
		for i := 0; i < repeatAmount; i++ {
			replaceNum := strconv.Itoa(i)
			for _, j := range *currentStackOperationCapturedLines {
				replaced := iterNameAsRegex.ReplaceAllString(j.OriginalLine, replaceNum)
				j.OriginalLine = replaced
				replacedLines[replacedIndex] = j
				replacedIndex++
			}
		}
	}

	blockStack.PopFromStackAndExtendCapturedLines(replacedLines)
	return nil
}
