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
	BlockOperationName  string
	OperandList         []Node
	CapturedLines       []CapturedLine
	AlternateStackBlock *StackBlock
}

func newStackBlock(operationName string, operandList []Node) StackBlock {
	return StackBlock{
		BlockOperationName: operationName,
		OperandList:        operandList,
	}
}

var Stack []StackBlock

var StackCapturesParentOpOnlyFlag bool = false
var StackWillClearFlag bool = false

// -----------------------------

func PushOntoStack(op string, operandList []Node) {
	Stack = append(Stack, newStackBlock(op, operandList))
	return
}

func PushIntoAlternateStackBlock(op string, operandList []Node) {
	currentStackOp := GetLastAlternateOperation()
	altBlock := newStackBlock(op, operandList)
	currentStackOp.AlternateStackBlock = &altBlock
	return
}

func GetCurrentOperation() *StackBlock {
	return &Stack[len(Stack)-1]
}

func GetLastAlternateOperation() *StackBlock {
	currentStackOp := GetCurrentOperation()
	for currentStackOp.AlternateStackBlock != nil {
		currentStackOp = currentStackOp.AlternateStackBlock
	}
	return currentStackOp
}

func popFromStack() {
	Stack = Stack[:len(Stack)-1]
}

func ClearStack() {
	Stack = Stack[:0]
}

func SetBottomOfStackToEmptyBlock() {
	Stack[0] = newStackBlock("nil", nil)
}

func ClearBottomOfStackCapturedLines() {
	Stack[0].CapturedLines = Stack[0].CapturedLines[:0]
}

func SetCaptureParentOpOnly() {
	StackCapturesParentOpOnlyFlag = true
}
func ClearCaptureParentOpOnly() {
	StackCapturesParentOpOnlyFlag = false
}

//+++++++++++++++++++++++++++++++

var correspondingEndBlockOperations = map[string]string{
	"REPEAT": "ENDREPEAT",
	"IF":     "ENDIF",
	"MACRO":  "ENDM",
}

//--------------------------------

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	switch lineOperationParsedValues.OperationTokenEnum {
	case enumTokenTypes.DIRECTIVE_blockStart:
		if StackCapturesParentOpOnlyFlag {
			return false
		}
		return true
	case enumTokenTypes.DIRECTIVE_blockEnd, enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		currentStackOp := GetCurrentOperation()
		endOpName, _ := correspondingEndBlockOperations[currentStackOp.BlockOperationName]
		if endOpName == strings.ToUpper(lineOperationParsedValues.OperationTokenValue) {
			return true
		} else if StackCapturesParentOpOnlyFlag {
			return false
		}
	}
	return false
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++

type captureableOpMap = map[enumTokenTypes.Def]bool

var sharedCapturableOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
}

var allowedOperationsForParentOps = map[string]captureableOpMap{
	"REPEAT": sharedCapturableOps,
	"IF":     sharedCapturableOps,
	"ELSEIF": sharedCapturableOps,
	"ELSE":   sharedCapturableOps,
	"MACRO": func() captureableOpMap {
		m := make(captureableOpMap)
		// Copy shared operations
		for k, v := range sharedCapturableOps {
			m[k] = v
		}
		m[enumTokenTypes.DIRECTIVE_blockStart] = true
		m[enumTokenTypes.DIRECTIVE_blockEnd] = true
		return m
	}(),
}

//-----------------------------------------------------

func CheckOperationIsCapturableAndAppend(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	currentStackOp := GetLastAlternateOperation()
	currentStackOpValue := currentStackOp.BlockOperationName
	checka, ok := allowedOperationsForParentOps[currentStackOpValue]
	if !ok {
		errMsg := fmt.Sprintf("Very bad stack op! Got: %v", currentStackOpValue)
		panic(errMsg)
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

// Take top of the stack and append all of it to the next unit down and pop the top
func PopFromStackAndExtendCapturedLines(extendLines []CapturedLine) {
	if len(Stack) > 1 {
		popFromStack()
		newCurrentStackOperation := GetLastAlternateOperation()
		for _, line := range extendLines {
			newCurrentStackOperation.CapturedLines = append(newCurrentStackOperation.CapturedLines, line)
		}

	} else if len(Stack) == 1 {
		newCurrentStackOperation := GetLastAlternateOperation()
		newCurrentStackOperation.CapturedLines = extendLines
		StackWillClearFlag = true

	} else {
		panic("Popping nothing from stack!!!")
	}
}

//--------------------------------

func CheckIfEndOperationAndClearStack(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	if (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd) &&
		StackWillClearFlag {
		StackWillClearFlag = false
		return true
	}
	return false
}
