package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func evalIf(directiveName string, operandList *[]Node) error {
	blockStack.PushOntoTopEntry(directiveName, *operandList)
	return nil
}

func evalElseIf(directiveName string, operandList *[]Node) error {
	lastOpName := blockStack.GetCurrentBlockEntryOperationName()
	//See if else has been declared - can't do else if afterwards
	if lastOpName == "ELSE" {
		return errorHandler.AddNew(enumErrorCodes.IfStatementElseIfAfterElse)
	}
	blockStack.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalElse(directiveName string, operandList *[]Node) error {
	lastOpName := blockStack.GetCurrentBlockEntryOperationName()
	//See if else has been declared - can't do duplicate elses
	if lastOpName == "ELSE" {
		return errorHandler.AddNew(enumErrorCodes.IfStatementDuplicateElse)
	}

	*operandList = append(*operandList, operandFactory.CreateBooleanLiteralNode(true))
	blockStack.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalEndIf(operandList *[]Node) error {
	currentStackOperation := blockStack.GetCurrentBlockEntry()
	var trueStatementCapturedLines *[]blockStack.CapturedLine

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
		emptyCapturedLines := make([]blockStack.CapturedLine, 0)
		trueStatementCapturedLines = &emptyCapturedLines
	}

	blockStack.PopTopEntryThenExtendCapturedLines(*trueStatementCapturedLines)

	currentStack := blockStack.GetCurrentBlockEntries()
	if blockStack.GoToProcessingFlag {
		if currentStackOperation != nil {
			(*currentStack)[0] = *currentStackOperation
			(*currentStack)[0].AlternateStackBlock = nil
		} else {
			blockStack.SetBottomOfStackToEmptyBlock()
		}
	}

	return nil
}
