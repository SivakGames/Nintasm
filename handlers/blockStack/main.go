package blockStack

import (
	"errors"
	"fmt"
	enumParserTypes "misc/nintasm/enums/parserTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/util"
	"strings"
)

type Node = operandFactory.Node

type CapturedLine struct {
	OriginalLine         string
	OperationLabel       string
	OperationTokenEnum   enumTokenTypes.Def
	OperationTokenValue  string
	OperandStartPosition int
	ParentParserEnum     enumParserTypes.Def
}

func newCapturedLine(originalLine string,
	operationLabel string,
	operationTokenEnum enumTokenTypes.Def,
	operationTokenValue string,
	operandStartPosition int,
	parentParserEnum enumParserTypes.Def,
) CapturedLine {
	return CapturedLine{
		OriginalLine:         originalLine,
		OperationLabel:       operationLabel,
		OperationTokenEnum:   operationTokenEnum,
		OperationTokenValue:  operationTokenValue,
		OperandStartPosition: operandStartPosition,
		ParentParserEnum:     parentParserEnum,
	}
}

type StackBlock struct {
	blockOperationName  string
	OperandList         []Node
	CapturedLines       []CapturedLine
	AlternateStackBlock *StackBlock
}

func newStackBlock(operationName string, operandList []Node) StackBlock {
	return StackBlock{
		blockOperationName: operationName,
		OperandList:        operandList,
	}
}

var Stack []StackBlock

var stackWillClear bool = false

// -----------------------------

func PushOntoStack(op string, operandList []Node) {
	Stack = append(Stack, newStackBlock(op, operandList))
	return
}

func PushIntoAlternateStackBlock(op string, operandList []Node) {
	currentStackOp := &Stack[len(Stack)-1]
	for currentStackOp.AlternateStackBlock != nil {
		currentStackOp = currentStackOp.AlternateStackBlock
	}
	altBlock := newStackBlock(op, operandList)
	currentStackOp.AlternateStackBlock = &altBlock

	fmt.Println(Stack)

	return
}

//+++++++++++++++++++++++++++++++

var correspondingEndBlockOperations = map[string]string{
	"REPEAT": "ENDREPEAT",
	"IF":     "ENDIF",
}

//--------------------------------

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	switch lineOperationParsedValues.OperationTokenEnum {
	case enumTokenTypes.DIRECTIVE_blockStart:
		return true
	case enumTokenTypes.DIRECTIVE_blockEnd:
		currentStackOp := &Stack[len(Stack)-1]
		endOpName, _ := correspondingEndBlockOperations[currentStackOp.blockOperationName]
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
	"IF": {
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
	currentStackOpValue := currentStackOp.blockOperationName
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
		lineOperationParsedValues.OperationLabel,
		lineOperationParsedValues.OperationTokenEnum,
		lineOperationParsedValues.OperationTokenValue,
		lineOperationParsedValues.OperandStartPosition,
		lineOperationParsedValues.ParentParserEnum,
	))

	return nil
}

//--------------------------------

func popFromStack() {
	Stack = Stack[:len(Stack)-1]
	return
}

func GetCurrentOperation() *StackBlock {
	return &Stack[len(Stack)-1]
}

//--------------------------------

// Take top of the stack and append all of it to the next unit down and pop the top
func PopFromStackAndExtendCapturedLines(extendLines []CapturedLine) {
	if len(Stack) > 1 {
		popFromStack()
		newCurrentStackOperation := GetCurrentOperation()
		for _, line := range extendLines {
			newCurrentStackOperation.CapturedLines = append(newCurrentStackOperation.CapturedLines, line)
		}

	} else if len(Stack) == 1 {
		newCurrentStackOperation := GetCurrentOperation()
		newCurrentStackOperation.CapturedLines = extendLines
		stackWillClear = true

	} else {
		panic("Popping nothing from stack!!!")
	}
}

//--------------------------------

func CheckIfEndOperationAndClearStack(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	if lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd &&
		stackWillClear {
		stackWillClear = false
		return true
	}
	return false
}

func ClearStack() {
	Stack = Stack[:0]
	return
}
