package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/parser/operandFactory"
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
		return errors.New("Repeat must be numeric and >= 1")
	}

	evaluatedNodes := []Node{numRepeatsNode}

	if len(*operandList) > 1 {
		iterNameNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsSubstitutionID(&iterNameNode) {
			return errors.New("Bad iterator name for repeat. Must be an \\iter")
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[1])
	}

	blockStack.PushOntoStack(directiveName, evaluatedNodes)
	return nil
}

func evalEndRepeat(directiveName string, operandList *[]Node) error {
	currentStackOperation := &blockStack.Stack[len(blockStack.Stack)-1]
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
	if len(blockStack.Stack) > 1 {
		fmt.Println("RUNNING")
		blockStack.Stack = blockStack.Stack[:len(blockStack.Stack)-1]
		newCurrentStackOperation := &blockStack.Stack[len(blockStack.Stack)-1]
		for _, line := range replacedLines {
			newCurrentStackOperation.CapturedLines = append(newCurrentStackOperation.CapturedLines, line)
		}
	}

	fmt.Println(blockStack.Stack)

	return nil
}
