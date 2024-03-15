package blockStack

import (
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/util"
	"strings"
)

// +++++++++++++++++++++++++++++++++++++++++++++

var currentBlockOperationLabel string = ""
var invokeOperations InvokeOperation = newInvokeOperation()
var GoToProcessingFlag bool = false

// +++++++++++++++++++++++++++++++++++++++++++++

func getCurrentInvokeOperation() *InvokeOperation {
	var invokeOp *InvokeOperation = &invokeOperations
	for invokeOp.nextCollection != nil {
		invokeOp = invokeOp.nextCollection
	}
	return invokeOp
}

func getCurrentInvokeOperationBlockEntries() *[]blockEntry {
	currentInvokeOp := getCurrentInvokeOperation()
	return &currentInvokeOp.blockEntries
}

func getCurrentInvokeOperationTopBlockEntry() *blockEntry {
	blockEntries := getCurrentInvokeOperationBlockEntries()
	return &(*blockEntries)[len(*blockEntries)-1]
}

func getCurrentInvokeOperationTopBlockEntryFurthestAlternate() *blockEntry {
	blockEntry := getCurrentInvokeOperationTopBlockEntry()
	for blockEntry.AlternateStackBlock != nil {
		blockEntry = blockEntry.AlternateStackBlock
	}
	return blockEntry
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func pushOntoCurrentInvokeOperationEntries(blockOperationName string, operandList []Node) {
	blockEntries := getCurrentInvokeOperationBlockEntries()
	*blockEntries = append(*blockEntries, newBlockEntry(blockOperationName, operandList))
}
func popFromCurrentInvokeOperationEntries() {
	blockEntries := getCurrentInvokeOperationBlockEntries()
	*blockEntries = (*blockEntries)[:len(*blockEntries)-1]
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
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
		panic(fmt.Sprintf("🛑 Somehow entering another label block operation while first (%v) is not done...", currentBlockOperationLabel))
	}
	currentBlockOperationLabel = label
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCurrentInvokeOperationEvalFlag() {
	currentInvokeOp := getCurrentInvokeOperation()
	currentInvokeOp.evalutesInsteadOfCapturing = false
}
func GetCurrentInvokeOperationEvalFlag() bool {
	currentInvokeOp := getCurrentInvokeOperation()
	return currentInvokeOp.evalutesInsteadOfCapturing
}
func SetCurrentInvokeOperationEvalFlag() {
	currentInvokeOp := getCurrentInvokeOperation()
	currentInvokeOp.evalutesInsteadOfCapturing = true
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCurrentInvokeOperationForcedCapturingFlag() {
	currentInvokeOp := getCurrentInvokeOperation()
	currentInvokeOp.forcedCapturing = false
}
func GetCurrentInvokeOperationForcedCapturingFlag() bool {
	currentInvokeOp := getCurrentInvokeOperation()
	return currentInvokeOp.forcedCapturing
}
func SetCurrentInvokeOperationForcedCapturingFlag() {
	currentInvokeOp := getCurrentInvokeOperation()
	currentInvokeOp.forcedCapturing = true
}

func EndLabeledDirective() {
	ClearCurrentOperationLabel()
	ForcePopTopEntry()
}

func SetBottomOfStackToEmptyBlock() {
	currentStack := getCurrentInvokeOperationBlockEntries()
	(*currentStack)[0] = newBlockEntry("nil", nil)
}

// ====================================================
// ====================================================

// When processing, will add another potentially usable collection
// in case an operation would make a new stack
func AddNewInvokeOperationCollection() {
	highestOp := getCurrentInvokeOperation()
	newOp := newInvokeOperation()
	(*highestOp).nextCollection = &newOp
}

func DestroyTempCollection(ptr *InvokeOperation) {
	var invokeOp *InvokeOperation = &invokeOperations

	for invokeOp.nextCollection != nil {
		if invokeOp.nextCollection == ptr {
			invokeOp.nextCollection = nil
			break
		}
		invokeOp = invokeOperations.nextCollection
	}
}

func PushOntoTopEntry(blockOperationName string, operandList []Node) {
	pushOntoCurrentInvokeOperationEntries(blockOperationName, operandList)
	flags := getStartOperationFlags(blockOperationName)
	currentOp := getCurrentInvokeOperation()
	currentOp.evalutesInsteadOfCapturing = flags.ForcedEval
	currentOp.forcedCapturing = flags.ForcedCapture
}

func ForcePopTopEntry() {
	popFromCurrentInvokeOperationEntries()
}

func CreateNewAlternateForTopEntry(blockOperationName string, operandList []Node) {
	curr := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
	altBlock := newBlockEntry(blockOperationName, operandList)
	curr.AlternateStackBlock = &altBlock
}

func PopTopEntryThenExtendCapturedLines(extendedLines []CapturedLine) {
	blockEntries := getCurrentInvokeOperationBlockEntries()

	// More than 1 will
	if len(*blockEntries) > 1 {
		popFromCurrentInvokeOperationEntries()
		blockEntry := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
		for _, line := range extendedLines {
			blockEntry.CapturedLines = append(blockEntry.CapturedLines, line)
		}

	} else if len(*blockEntries) == 1 {
		//Set eval operands to true
		blockEntry := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
		blockEntry.CapturedLines = extendedLines
		GoToProcessingFlag = true

	} else {
		panic("🛑 Popping nothing/extending nothing!!!")
	}
}

func GetCurrentBlockEntries() *[]blockEntry {
	return getCurrentInvokeOperationBlockEntries()
}

func GetCurrentBlockEntry() *blockEntry {
	currentStackOp := getCurrentInvokeOperationTopBlockEntry()
	return currentStackOp
}

func GetCurrentBlockEntryOperationName() string {
	curr := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
	return curr.BlockOperationName
}

func GetCurrentBlockEntryCapturedLines() *[]CapturedLine {
	curr := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
	return &curr.CapturedLines
}

func CheckIfEndOpMatchesOpeningOp(desiredEndOpName string) bool {
	currentStackOp := getCurrentInvokeOperationTopBlockEntry()
	endOpName, _ := correspondingEndBlockOperations[currentStackOp.BlockOperationName]
	return endOpName == strings.ToUpper(desiredEndOpName)
}

func GetTopBlockEntryData() (*[]CapturedLine, *[]Node) {
	topEntry := getCurrentInvokeOperationTopBlockEntry()
	return &topEntry.CapturedLines, &topEntry.OperandList
}

// ***************************************************
func GetCurrentOpPtr() *InvokeOperation {
	return getCurrentInvokeOperation()
}
func GetLinesWithPtr(pointer *InvokeOperation) *[]CapturedLine {
	return &pointer.blockEntries[0].CapturedLines
}
func GetBlockEntriesWithPtr(pointer *InvokeOperation) *[]blockEntry {
	return &pointer.blockEntries
}
func ClearBlockEntriesWithPtr(pointer *InvokeOperation) {
	pointer.blockEntries = (*pointer).blockEntries[:0]
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++

func CheckIfNewStartEndOperation(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	isStartEndEnum := (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockStart ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd)

	// If op isn't start/end enum no need to proceed...
	if !isStartEndEnum {
		return false
	}

	//If in forced evaluate mode, see if there is a pair to force-close it

	if GetCurrentInvokeOperationForcedCapturingFlag() &&
		!CheckIfEndOpMatchesOpeningOp(lineOperationParsedValues.OperationTokenValue) {
		return false
	}

	return true
}

//-----------------------------------------------------

func CheckOperationIsCapturableAndAppend(
	originalLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	err := checkOperationIsCapturableByCurrentBlockOperation(lineOperationParsedValues)
	if err != nil {
		return err
	}
	currentStackOp := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
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

func CheckIfEndOperationAndGoesToProcessing(lineOperationParsedValues *util.LineOperationParsedValues) bool {
	if (lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		lineOperationParsedValues.OperationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd) &&
		GoToProcessingFlag {
		GoToProcessingFlag = false
		return true
	}
	return false
}
