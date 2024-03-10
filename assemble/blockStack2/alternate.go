package blockStack2

import (
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/util"
	"strings"
)

// ++++++++++++++++++++++++++++++++++++

type blockEntry struct {
	blockOperationName  string
	capturedLines       []CapturedLine
	operandList         []Node
	alternateStackBlock *blockEntry
}

func newBlockEntry(blockOperationName string, operandList []Node) blockEntry {
	return blockEntry{
		blockOperationName:  blockOperationName,
		capturedLines:       []CapturedLine{},
		operandList:         operandList,
		alternateStackBlock: nil,
	}
}

// ++++++++++++++++++++++++++++++++++++

type InvokeOperation struct {
	blockEntries []blockEntry
	//Setting where the operation evaluates things while capturing
	evalutesInsteadOfCapturing bool
	//Mainly for macros - Will always capture nodes except for a corresponding ending block
	forcedCapturing bool
	nextCollection  *InvokeOperation
}

func newInvokeOperation() InvokeOperation {
	return InvokeOperation{
		blockEntries:               []blockEntry{},
		evalutesInsteadOfCapturing: false,
		forcedCapturing:            false,
		nextCollection:             nil,
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++

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
	for blockEntry.alternateStackBlock != nil {
		blockEntry = blockEntry.alternateStackBlock
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
	var exited bool = false

	for invokeOp.nextCollection != nil {
		if invokeOp.nextCollection == ptr {
			invokeOp.nextCollection = nil
			exited = true
			break
		}
		invokeOp = invokeOperations.nextCollection
	}
	if exited {
		fmt.Println("Exited legit")
	} else {
		fmt.Println("Exited not legit")
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

func PopTopEntryThenExtendCapturedLines(extendedLines []CapturedLine) {
	blockEntries := getCurrentInvokeOperationBlockEntries()

	// More than 1 will
	if len(*blockEntries) > 1 {
		popFromCurrentInvokeOperationEntries()
		blockEntry := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
		for _, line := range extendedLines {
			blockEntry.capturedLines = append(blockEntry.capturedLines, line)
		}

	} else if len(*blockEntries) == 1 {
		//Set eval operands to true
		blockEntry := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
		blockEntry.capturedLines = extendedLines
		GoToProcessingFlag = true
	} else {
		panic("🛑 Popping nothing/extending nothing!!!")
	}
}

func GetCurrentBlockEntries() *[]blockEntry {
	return getCurrentInvokeOperationBlockEntries()
}

func GetCurrentBlockEntryCapturedLines() *[]CapturedLine {
	curr := getCurrentInvokeOperationTopBlockEntryFurthestAlternate()
	return &curr.capturedLines
}

func CheckIfEndOpMatchesOpeningOp(desiredEndOpName string) bool {
	currentStackOp := getCurrentInvokeOperationTopBlockEntry()
	endOpName, _ := correspondingEndBlockOperations[currentStackOp.blockOperationName]
	return endOpName == strings.ToUpper(desiredEndOpName)
}

func GetCurrentOpPtr() *InvokeOperation {
	return getCurrentInvokeOperation()
}

// ***************************************************

func GetLinesWithPtr(pointer *InvokeOperation) *[]CapturedLine {
	return &pointer.blockEntries[0].capturedLines
}
func GetBlockEntriesWithPtr(pointer *InvokeOperation) *[]blockEntry {
	return &pointer.blockEntries
}
func ClearBlockEntriesWithPtr(pointer *InvokeOperation) {
	pointer.blockEntries = (*pointer).blockEntries[:0]
}

func GetTopBlockEntryData() (*[]CapturedLine, *[]Node) {
	topEntry := getCurrentInvokeOperationTopBlockEntry()
	return &topEntry.capturedLines, &topEntry.operandList
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
	currentStackOp.capturedLines = append(currentStackOp.capturedLines, newCapturedLine(
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
