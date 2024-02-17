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

var StackWillClearFlag bool = false
var stackCapturesParentOpOnlyFlag bool = false
var currentBlockOperationLabel string = ""

// -----------------------------

func PushOntoStack(op string, operandList []Node) {
	Stack = append(Stack, newStackBlock(op, operandList))
	return
}

func popFromStack() {
	Stack = Stack[:len(Stack)-1]
}

func ClearStack() {
	Stack = Stack[:0]
}

// -----------------

func AppendToTopOfStackAlternateBlock(op string, operandList []Node) {
	currentStackOp := GetTopOfStackLastAlternateOperation()
	altBlock := newStackBlock(op, operandList)
	currentStackOp.AlternateStackBlock = &altBlock
	return
}

// Get whatever's on top (and thus current)
func GetTopOfStackOperation() *StackBlock {
	return &Stack[len(Stack)-1]
}

func GetTopOfStackLastAlternateOperation() *StackBlock {
	currentStackOp := GetTopOfStackOperation()
	for currentStackOp.AlternateStackBlock != nil {
		currentStackOp = currentStackOp.AlternateStackBlock
	}
	return currentStackOp
}

// -----------------

func SetBottomOfStackToEmptyBlock() {
	Stack[0] = newStackBlock("nil", nil)
}

func ClearBottomOfStackCapturedLines() {
	Stack[0].CapturedLines = Stack[0].CapturedLines[:0]
}

// -----------------

// Stack will only handle the parent op. No nested ops.
func SetCaptureParentOpOnlyFlag() {
	stackCapturesParentOpOnlyFlag = true
}
func ClearCaptureParentOpOnlyFlag() {
	stackCapturesParentOpOnlyFlag = false
}

// -----------------

// Will set the label of the labeled operation that will be captured.
// If one was previously set then error because it hasn't finished.
func SetCurrentOperationLabel(label string) error {
	if currentBlockOperationLabel != "" {
		return errors.New("Somehow entering another label block operation while first is not done...")
	}

	currentBlockOperationLabel = label
	return nil
}
func ClearCurrentOperationLabel() {
	currentBlockOperationLabel = ""
	return
}

func GetCurrentOperationLabel() string {
	return currentBlockOperationLabel
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
		if stackCapturesParentOpOnlyFlag {
			return false
		}
		return true
	case enumTokenTypes.DIRECTIVE_blockEnd, enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		currentStackOp := GetTopOfStackOperation()
		endOpName, _ := correspondingEndBlockOperations[currentStackOp.BlockOperationName]
		if endOpName == strings.ToUpper(lineOperationParsedValues.OperationTokenValue) {
			return true
		} else if stackCapturesParentOpOnlyFlag {
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
	currentStackOp := GetTopOfStackLastAlternateOperation()
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
		newCurrentStackOperation := GetTopOfStackLastAlternateOperation()
		for _, line := range extendLines {
			newCurrentStackOperation.CapturedLines = append(newCurrentStackOperation.CapturedLines, line)
		}

	} else if len(Stack) == 1 {
		newCurrentStackOperation := GetTopOfStackLastAlternateOperation()
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
