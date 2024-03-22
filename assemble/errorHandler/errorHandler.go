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
const ERROR_STACK_CAPACITY = 10

// For underlining parts of the line
var highlightStart int = -1
var highlightEnd int = -1
var totalErrors uint = 0

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
	code        enumErrorCodes.Def
	lineNumber  uint
	lineContent string
	fileName    string
	message     string
	hint        string
	severity    enumErrorCodes.Severity
}

func newErrorEntry(
	code enumErrorCodes.Def,
	message string,
	severity enumErrorCodes.Severity,
) ErrorEntry {
	fileData := fileStack.GetTopOfFileStack()

	if fileData != nil {
		return ErrorEntry{
			code:        code,
			message:     message,
			fileName:    fileData.FileName,
			lineNumber:  fileData.CurrentLineNumber,
			lineContent: fileData.ProcessedLines[fileData.CurrentLineNumber-1],
			severity:    severity,
		}
	}
	return ErrorEntry{
		code:        code,
		message:     message,
		fileName:    noFileDefaults.fileName,
		lineNumber:  noFileDefaults.lineNumber,
		lineContent: noFileDefaults.lineContent,
		severity:    severity,
	}
}

var errorStack = make([]ErrorEntry, ERROR_STACK_CAPACITY+1)

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) error {
	defer func() {
		totalErrors++
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
	errorStack[totalErrors] = entry

	if totalErrors+2 > ERROR_STACK_CAPACITY {
		totalErrors++
		tooManyErrorsEntry := newErrorEntry(errorTableKey, "Too many errors!", enumErrorCodes.Fatal)
		tooManyErrorsEntry.fileName = "OUTPUT LIMIT REACHED"
		tooManyErrorsEntry.lineNumber = 0
		tooManyErrorsEntry.lineContent = ""
		errorStack[totalErrors] = tooManyErrorsEntry
		return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, tooManyErrorsEntry.severity))
	}

	return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, entry.severity))
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
	return totalErrors
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
