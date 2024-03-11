package directiveHandler

import (
	"misc/nintasm/assemble/blockStack2"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func evalIf(directiveName string, operandList *[]Node) error {
	blockStack2.PushOntoTopEntry(directiveName, *operandList)
	return nil
}

func evalElseIf(directiveName string, operandList *[]Node) error {
	lastOpName := blockStack2.GetCurrentBlockEntryOperationName()
	//See if else has been declared - can't do else if afterwards
	if lastOpName == "ELSE" {
		return errorHandler.AddNew(enumErrorCodes.IfStatementElseIfAfterElse)
	}
	blockStack2.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalElse(directiveName string, operandList *[]Node) error {
	lastOpName := blockStack2.GetCurrentBlockEntryOperationName()
	//See if else has been declared - can't do duplicate elses
	if lastOpName == "ELSE" {
		return errorHandler.AddNew(enumErrorCodes.IfStatementDuplicateElse)
	}

	*operandList = append(*operandList, operandFactory.CreateBooleanLiteralNode(true))
	blockStack2.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalEndIf(operandList *[]Node) error {
	currentStackOperation := blockStack2.GetCurrentBlockEntry()
	var trueStatementCapturedLines *[]blockStack2.CapturedLine

	// Cycle through until finding a true block or a nil one (nothing is true)
	// Will change currentStackOp

	for currentStackOperation != nil {
		ifData := &currentStackOperation.OperandList[0]
		if !operandFactory.ValidateNodeIsBoolean(ifData) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotBool)
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
		emptyCapturedLines := make([]blockStack2.CapturedLine, 0)
		trueStatementCapturedLines = &emptyCapturedLines
	}

	blockStack2.PopTopEntryThenExtendCapturedLines(*trueStatementCapturedLines)

	//currentStack := blockStack.GetCurrentStack()
	/*
		currentStack := blockStack2.GetCurrentBlockEntries()
		if blockStack.StackWillClearFlag {
			if currentStackOperation != nil {
				(*currentStack)[0] = *currentStackOperation
				(*currentStack)[0].AlternateStackBlock = nil
			} else {
				blockStack.SetBottomOfStackToEmptyBlock()
			}
		} */

	return nil
}
