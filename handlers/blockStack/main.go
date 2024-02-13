package blockStack

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/util"
	"strings"
)

type Node = operandFactory.Node

type CapturedLine struct {
	OriginalLine         string
	operationTokenEnum   enumTokenTypes.Def
	operationTokenValue  string
	operandStartPosition int
}

func newCapturedLine(originalLine string,
	operationTokenEnum enumTokenTypes.Def,
	operationTokenValue string,
	operandStartPosition int) CapturedLine {
	return CapturedLine{
		OriginalLine:         originalLine,
		operationTokenEnum:   operationTokenEnum,
		operationTokenValue:  operationTokenValue,
		operandStartPosition: operandStartPosition,
	}
}

type Bloxx struct {
	blockOperationValue string
	OperandList         []Node
	CapturedLines       []CapturedLine
}

func newBloxx(op string, operandList []Node) Bloxx {
	return Bloxx{
		blockOperationValue: op,
		OperandList:         operandList,
	}
}

var Stack []Bloxx

var stackWillClear bool = false

// -----------------------------

func PushOntoStack(op string, operandList []Node) {
	Stack = append(Stack, newBloxx(op, operandList))
	return
}

//+++++++++++++++++++++++++++++++

var correspondingEndBlockOperations = map[string]string{
	"REPEAT": "ENDREPEAT",
}

//--------------------------------

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	switch lineOperationParsedValues.OperationTokenEnum {
	case enumTokenTypes.DIRECTIVE_blockStart:
		return true
	case enumTokenTypes.DIRECTIVE_blockEnd:
		currentStackOp := &Stack[len(Stack)-1]
		endOpName, _ := correspondingEndBlockOperations[currentStackOp.blockOperationValue]
		return lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd &&
			endOpName == strings.ToUpper(lineOperationParsedValues.OperationTokenValue)
	}
	return false
}

//+++++++++++++++++++++++++++++++

var allowedOperationsForParentOps = map[string]map[enumTokenTypes.Def]bool{
	"REPEAT": {
		enumTokenTypes.INSTRUCTION:          true,
		enumTokenTypes.DIRECTIVE_dataBytes:  true,
		enumTokenTypes.DIRECTIVE_dataSeries: true,
		enumTokenTypes.DIRECTIVE_mixedData:  true,
	},
}

//--------------------------------

func CheckOperationIsCapturableAndAppend(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	currentStackOp := &Stack[len(Stack)-1]
	currentStackOpValue := currentStackOp.blockOperationValue
	checka, ok := allowedOperationsForParentOps[currentStackOpValue]
	if !ok {
		panic("Very bad stack op!")
	}

	_, ok = checka[lineOperationParsedValues.OperationTokenEnum]
	if !ok {
		errMsg := fmt.Sprintf("Uncapturable! %v", lineOperationParsedValues.OperationTokenValue)
		return errors.New(errMsg)
	}
	currentStackOp.CapturedLines = append(currentStackOp.CapturedLines, newCapturedLine(
		originalLine,
		lineOperationParsedValues.OperationTokenEnum,
		lineOperationParsedValues.OperationTokenValue,
		lineOperationParsedValues.OperandStartPosition,
	))

	return nil
}

//--------------------------------

// Take top of the stack and append all of it to the next unit down and pop the top
func PopFromStackAndExtendCapturedLines(extendLines []CapturedLine) {
	if len(Stack) > 1 {
		Stack = Stack[:len(Stack)-1]
		newCurrentStackOperation := &Stack[len(Stack)-1]
		for _, line := range extendLines {
			newCurrentStackOperation.CapturedLines = append(newCurrentStackOperation.CapturedLines, line)
		}

		// This is the last stack entry being popped

	} else if len(Stack) == 1 {
		newCurrentStackOperation := &Stack[len(Stack)-1]
		newCurrentStackOperation.CapturedLines = extendLines
		stackWillClear = true
	} else {
		panic("Popping nothing from stack!!!")
	}
}

//--------------------------------

func CheckIfEndOperationAndClearStack(lineOperationParsedValues *util.LineOperationParsedValues) {
	if lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd &&
		stackWillClear {
		Stack = Stack[:0]
		stackWillClear = false
	}
}
