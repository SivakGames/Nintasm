package handlerDirective

import (
	"errors"
	"misc/nintasm/assemble/blockStack"
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

	if !(operandFactory.ValidateNodeIsNumeric(&numRepeatsNode) &&
		operandFactory.ValidateNumericNodeIsGTEValue(&numRepeatsNode, 1)) {
		return errors.New("Repeat must be numeric and >= 1") // ❌ Fails
	}

	evaluatedNodes := []Node{numRepeatsNode}

	if len(*operandList) > 1 {
		iterNameNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsSubstitutionID(&iterNameNode) {
			return errors.New("Bad iterator name for repeat. Must be an \\iter") // ❌ Fails
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[1])
	}

	blockStack.PushOntoStack(directiveName, evaluatedNodes)

	return nil
}

func evalEndRepeat(directiveName string, operandList *[]Node) error {
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
