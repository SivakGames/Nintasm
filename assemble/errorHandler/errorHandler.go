package errorHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"strconv"
	"strings"
)

// ++++++++++++++++++++++++++++++++++++++++

const ERROR_CAPTION_MIN_WIDTH = 7
const ERROR_STACK_CAPACITY = 100

// For underlining parts of the line
var highlightStart int = -1
var highlightEnd int = -1
var currentHint string = ""
var totalQueuedErrors uint = 0
var totalErrors uint = 0
var totalWarnings uint = 0

// ++++++++++++++++++++++++++++++++++++++++

type resolveSymbolData struct {
	fileName    string
	lineNumber  uint
	lineContent string
}

var noFileDefaults = resolveSymbolData{
	fileName:    "NO FILE",
	lineNumber:  0,
	lineContent: "",
}

// ++++++++++++++++++++++++++++++++++++++++

type ErrorEntry struct {
	code             enumErrorCodes.Def
	severity         enumErrorCodes.Severity
	fileName         string
	lineNumber       uint
	lineContent      string
	message          string
	hint             string
	subOpLineNumber  uint
	subOpLineContent string
}

func newErrorEntry(
	code enumErrorCodes.Def,
	message string,
	severity enumErrorCodes.Severity,
) ErrorEntry {
	entry := ErrorEntry{
		code:             code,
		severity:         severity,
		message:          message,
		hint:             currentHint,
		fileName:         noFileDefaults.fileName,
		lineNumber:       noFileDefaults.lineNumber,
		lineContent:      noFileDefaults.lineContent,
		subOpLineContent: "",
		subOpLineNumber:  0,
	}
	fileData := fileStack.GetTopOfFileStack()
	if fileData != nil {
		entry.fileName = fileData.FileName
		entry.lineNumber = fileData.CurrentLineNumber
		entry.lineContent = fileData.ProcessedLines[fileData.CurrentLineNumber-1]
	}
	subOp := fileStack.GetSubOp()
	if subOp != nil {
		entry.subOpLineContent = subOp.LineContent
		entry.subOpLineNumber = subOp.SubOpLineNumber
	}
	return entry
}

var errorStack = make([]ErrorEntry, ERROR_STACK_CAPACITY+1)

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) error {
	defer func() {
		totalQueuedErrors++
		currentHint = ""
		Resetighlights()
	}()

	var entry ErrorEntry

	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry = newErrorEntry(errorTableKey, errMsg, errData.severity)
	} else {
		entry = newErrorEntry(errorTableKey, "Non-error-code error???", enumErrorCodes.Fatal)
	}
	errorStack[totalQueuedErrors] = entry

	if totalQueuedErrors+2 > ERROR_STACK_CAPACITY {
		totalQueuedErrors++
		tooManyErrorsEntry := newErrorEntry(errorTableKey, "Too many errors!", enumErrorCodes.Fatal)
		tooManyErrorsEntry.fileName = "OUTPUT LIMIT REACHED"
		tooManyErrorsEntry.lineNumber = 0
		tooManyErrorsEntry.lineContent = ""
		errorStack[totalQueuedErrors] = tooManyErrorsEntry
		return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, tooManyErrorsEntry.severity))
	}

	switch entry.severity {
	case enumErrorCodes.Fatal,
		enumErrorCodes.Error:
		totalErrors++
	case enumErrorCodes.Warning:
		totalWarnings++
	}

	return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, entry.severity))
}

func AddHint(errorTableKey enumErrorCodes.Def, args ...interface{}) {
	hint, hintExists := errorHintTable[errorTableKey]
	if !hintExists {
		AddNew(enumErrorCodes.Other, "Specified hint doesn't exist!")
		return
	}
	currentHint = fmt.Sprintf(hint, args...)
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

// A silent error initially...
func AddUnresolved(symbolName string) error {
	return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, enumErrorCodes.UnresolvedIdentifier))
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func OverwriteNoFileDefaults(fileName string, lineNumber uint, lineContent string) {
	noFileDefaults.fileName = fileName
	noFileDefaults.lineNumber = lineNumber
	noFileDefaults.lineContent = lineContent
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func IsErrorCode(err error) (bool, enumErrorCodes.Def) {
	errorCode := err.Error()
	errorTableKey := enumErrorCodes.Def(errorCode)
	_, isValidErrorCode := errorTable[errorTableKey]
	return isValidErrorCode, errorTableKey
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddHighlights(start int, end int) {
	highlightStart = start
	highlightEnd = end
}

func Resetighlights() {
	highlightEnd = -1
}

// ------------------------------------------------

func GetErrorCount() uint {
	return totalQueuedErrors
}

// ------------------------------------------------

// If severity is <= threshold it should STOP propagating up
func CheckErrorContinuesUpwardPropagation(err error, threshold enumErrorCodes.Severity) error {
	severityValue := err.Error()

	modded, ok := strings.CutPrefix(severityValue, SEVERITY_PREFIX)
	if ok {
		severityAmt, err := strconv.Atoi(modded)
		if err != nil {
			return err
		}
		if severityAmt <= int(threshold) {
			return nil
		}
	}
	return err
}
