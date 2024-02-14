package handlerDirective

import (
	"errors"
	"fmt"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/parser/operandFactory"
)

func evalIf(directiveName string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalElseIf(directiveName string, operandList *[]Node) error {
	blockStack.PushIntoAlternateStackBlock(directiveName, *operandList)
	return nil
}

func evalEndIf(directiveName string, operandList *[]Node) error {
	currentStackOperation := &blockStack.Stack[len(blockStack.Stack)-1]

	for currentStackOperation != nil {
		ifData := &currentStackOperation.OperandList[0]
		if !operandFactory.ValidateNodeIsBoolean(ifData) {
			return errors.New("UNEXPECTED NON-BOOLEAN FOR IF!")
		}
		if ifData.AsBool {
			break
		}
		currentStackOperation = currentStackOperation.AlternateStackBlock
	}

	fmt.Println("At the end!")

	//blockStack.PopFromStackAndExtendCapturedLines(replacedLines)
	return nil
}
