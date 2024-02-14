package handlerDirective

import (
	"errors"
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
	currentStackOperation := blockStack.GetCurrentOperation()
	var trueStatementCapturedLines *[]blockStack.CapturedLine

	// Cycle through until finding a true block or a nil one (nothing is true)
	// Will change currentStackOp

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

	if currentStackOperation != nil {
		trueStatementCapturedLines = &currentStackOperation.CapturedLines
	} else {
		emptyCapturedLines := make([]blockStack.CapturedLine, 0)
		trueStatementCapturedLines = &emptyCapturedLines
	}

	blockStack.PopFromStackAndExtendCapturedLines(*trueStatementCapturedLines)
	return nil
}
