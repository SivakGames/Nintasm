package handlerDirective

import (
	"errors"
	"misc/nintasm/assemble/blockStack"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/operandFactory"
)

func evalIf(directiveName string, operandList *[]Node) error {
	blockStack.PushOntoStack(directiveName, *operandList)
	return nil
}

func evalElseIf(directiveName string, operandList *[]Node) error {
	lastOp := blockStack.GetTopOfStackLastAlternateOperation()
	if lastOp.BlockOperationName == "ELSE" {
		return errors.New("Cannot have elseif after else")
	}
	blockStack.AppendToTopOfStackAlternateBlock(directiveName, *operandList)
	return nil
}

func evalElse(directiveName string, operandList *[]Node) error {
	lastOp := blockStack.GetTopOfStackLastAlternateOperation()
	if lastOp.BlockOperationName == "ELSE" {
		return errors.New("Cannot only have 1 else in this block")
	}

	*operandList = append(*operandList, operandFactory.CreateBooleanLiteralNode(enumTokenTypes.NUMBER_decimal, "1", true))
	blockStack.AppendToTopOfStackAlternateBlock(directiveName, *operandList)
	return nil
}

func evalEndIf(directiveName string, operandList *[]Node) error {
	currentStackOperation := blockStack.GetTopOfStackOperation()
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

	// Nil signifies no if/elseif/else are true so nothing will be

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
			blockStack.SetBottomOfStackToEmptyBlock()
		}
	}

	return nil
}