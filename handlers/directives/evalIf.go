package handlerDirective

import (
	"errors"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/parser/operandFactory"
)

func evalIf(directiveName string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalElseIf(directiveName string, operandList *[]Node) error {
	lastOp := blockStack.GetLastAlternateOperation()
	if lastOp.BlockOperationName == "ELSE" {
		return errors.New("Cannot have elseif after else")
	}
	blockStack.PushIntoAlternateStackBlock(directiveName, *operandList)
	return nil
}

func evalElse(directiveName string, operandList *[]Node) error {
	lastOp := blockStack.GetLastAlternateOperation()
	if lastOp.BlockOperationName == "ELSE" {
		return errors.New("Cannot only have 1 else in this block")
	}

	*operandList = append(*operandList, operandFactory.CreateBooleanLiteralNode(enumTokenTypes.NUMBER_decimal, "1", true))
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

	if blockStack.StackWillClearFlag {
		if currentStackOperation != nil {
			blockStack.Stack[0] = *currentStackOperation
			blockStack.Stack[0].AlternateStackBlock = nil
		} else {
			blockStack.SetBottomOfStackToEmpty()
		}
	}

	return nil
}
