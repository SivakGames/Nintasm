package blockStack

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/operandFactory"
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

// Will evaluate the node rather than capturing it
var currentOperationEvaluatesCapturedNodesFlag bool = false
var stackCapturesParentOpOnlyFlag bool = false
var currentBlockOperationLabel string = ""

// Mainly used with namespace. Will temporarily act as the parent label for such operations
var temporaryOverwritingParentLabel string = ""

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

func GetTopOfStackCapturedLines() *[]CapturedLine {
	currentStackOp := GetTopOfStackOperation()
	capturedLines := &currentStackOp.CapturedLines
	return capturedLines
}

// -----------------

func SetBottomOfStackToEmptyBlock() {
	Stack[0] = newStackBlock("nil", nil)
}

func ClearBottomOfStackCapturedLines() {
	Stack[0].CapturedLines = Stack[0].CapturedLines[:0]
}

// -----------------

func ClearCaptureParentOpOnlyFlag() {
	stackCapturesParentOpOnlyFlag = false
}

// Stack will only handle the parent op. No nested ops.
func SetCaptureParentOpOnlyFlag() {
	stackCapturesParentOpOnlyFlag = true
}

// -----------------

func ClearCurrentOperationLabel() {
	currentBlockOperationLabel = ""
}
func GetCurrentOperationLabel() string {
	return currentBlockOperationLabel
}

// Will set the label of the labeled operation that will be captured.
// If one was previously set then error because it hasn't finished.
func SetCurrentOperationLabel(label string) error {
	if currentBlockOperationLabel != "" {
		panic("🛑 Somehow entering another label block operation while first is not done...")
	}

	currentBlockOperationLabel = label
	return nil
}

// -----------------

func ClearCurrentOperationEvaluatesCapturedNodesFlag() {
	currentOperationEvaluatesCapturedNodesFlag = false
}
func GetCurrentOperationEvaluatesCapturedNodesFlag() bool {
	return currentOperationEvaluatesCapturedNodesFlag
}
func SetCurrentOperationEvaluatesCapturedNodesFlag() {
	currentOperationEvaluatesCapturedNodesFlag = true
}

// --------------------------------

// Will clear the current labeled op's label and the capture parent op flag,
func GetLabelAndDoEndBlockSetups() string {
	blockLabel := GetCurrentOperationLabel()
	ClearCurrentOperationLabel()
	ClearCaptureParentOpOnlyFlag()
	return blockLabel
}

//--------------------------------

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	isStartEndEnum := (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd)

	// If op isn't start/end enum no need to proceed...
	if !isStartEndEnum {
		return false
	}

	//If capturing for only the parent op, only a valid closing block for the parent can end it

	if stackCapturesParentOpOnlyFlag &&
		!CheckIfEndOpMatchesOpeningOp(lineOperationParsedValues.OperationTokenValue) {
		return false
	}

	return true
}

//+++++++++++++++++++++++++++++++

var correspondingEndBlockOperations = map[string]string{
	"CHARMAP":   "ENDCHARMAP",
	"EXPRMAP":   "ENDEXPRMAP",
	"IF":        "ENDIF",
	"IKV":       "ENDIKV",
	"KVMACRO":   "ENDKVM",
	"MACRO":     "ENDM",
	"NAMESPACE": "ENDNAMESPACE",
	"REPEAT":    "ENDREPEAT",
}

// --------------------------------
func CheckIfEndOpMatchesOpeningOp(desiredEndOpName string) bool {
	currentStackOp := GetTopOfStackOperation()
	endOpName, _ := correspondingEndBlockOperations[currentStackOp.BlockOperationName]
	return endOpName == strings.ToUpper(desiredEndOpName)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++

type captureableOpMap = map[enumTokenTypes.Def]bool

var sharedCapturableOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
}

var sharedCapturableMacroOps = captureableOpMap{
	enumTokenTypes.INSTRUCTION:          true,
	enumTokenTypes.DIRECTIVE_dataBytes:  true,
	enumTokenTypes.DIRECTIVE_dataSeries: true,
	enumTokenTypes.DIRECTIVE_mixedData:  true,
	enumTokenTypes.DIRECTIVE_blockStart: true,
	enumTokenTypes.DIRECTIVE_blockEnd:   true,
}

var allowedOperationsForParentOps = map[string]captureableOpMap{
	"CHARMAP": {
		enumTokenTypes.DIRECTIVE_defCharMap: true,
	},
	"EXPRMAP": {
		enumTokenTypes.DIRECTIVE_defExprMap: true,
	},
	"IF":     sharedCapturableOps,
	"ELSEIF": sharedCapturableOps,
	"ELSE":   sharedCapturableOps,
	"IKV": {
		enumTokenTypes.DIRECTIVE_invokeKeyVal: true,
	},
	"KVMACRO": sharedCapturableMacroOps,
	"MACRO":   sharedCapturableMacroOps,
	"NAMESPACE": {
		enumTokenTypes.ASSIGN_simple: true,
	},
	"REPEAT": sharedCapturableOps,
}

//-----------------------------------------------------

func CheckOperationIsCapturable(
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
		return errorHandler.AddNew(enumErrorCodes.BlockOpUncapturableByParent, lineOperationParsedValues.OperationTokenValue)
	}
	return nil
}

func CheckOperationIsCapturableAndAppend(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	err := CheckOperationIsCapturable(originalLine, lineOperationParsedValues)
	if err != nil {
		return err
	}
	currentStackOp := GetTopOfStackLastAlternateOperation()
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
		panic("🛑 Popping nothing from stack!!!")
	}
}

// Helper to pop from the stack but won't supply any captured lines to extends.
// (Used for labeled block directives)
func PopFromStackAndExtendNoLines() {
	var emptyCapturedLines []CapturedLine
	PopFromStackAndExtendCapturedLines(emptyCapturedLines)
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
