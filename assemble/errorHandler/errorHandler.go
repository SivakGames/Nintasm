package errorHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/util"
	"strconv"
	"strings"
)

// ++++++++++++++++++++++++++++++++++++++++

const ERROR_CAPTION_MIN_WIDTH = 7
const ERROR_STACK_CAPACITY = 3

// For underlining parts of the line
var highlightStart int = -1
var highlightEnd int = -1
var totalErrors uint = 0

// ++++++++++++++++++++++++++++++++++++++++

type resolveSymbolData struct {
	lineNumber  uint
	lineContent string
	fileName    string
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

func NewErrorEntry(code enumErrorCodes.Def, message string, severity enumErrorCodes.Severity) ErrorEntry {
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
		entry = NewErrorEntry(errorTableKey, errMsg, errData.severity)
	} else {
		entry = NewErrorEntry(errorTableKey, "Non-error-code error???", enumErrorCodes.Fatal)
	}
	errorStack[totalErrors] = entry

	if totalErrors+2 > ERROR_STACK_CAPACITY {
		totalErrors++
		tooManyErrorsEntry := NewErrorEntry(errorTableKey, "Too many errors!", enumErrorCodes.Fatal)
		errorStack[totalErrors] = tooManyErrorsEntry
		return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, tooManyErrorsEntry.severity))
	}

	return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, entry.severity))
}

func CheckAndPrintErrors() {
	if totalErrors == 0 {
		return
	}
	for i := 0; i < int(totalErrors); i++ {
		PrintError(&errorStack[i])
	}

}

func PrintError(entry *ErrorEntry) {
	//File name (if any)
	colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", entry.fileName), "red", true)
	fmt.Println("░", colorizedFileName)

	//Line number and content
	colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), ERROR_CAPTION_MIN_WIDTH, ' '), "blue", true)
	fmt.Println("▓", colorizedLineNumber, entry.lineContent)

	subOp := fileStack.GetSubOp()
	if subOp != nil {
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", subOp.SubOpLineNumber), ERROR_CAPTION_MIN_WIDTH-2, ' '), "blue", true)
		fmt.Println("▓▓▓", colorizedLineNumber, subOp.LineContent)
	}

	if highlightEnd > highlightStart {
		marginSpacer := util.Colorize(util.PadStringLeft("", ERROR_CAPTION_MIN_WIDTH+2, '░'), "cyan", false)
		leadingSpace := util.PadStringLeft("", highlightStart, ' ')
		arrows := util.Colorize(util.PadStringLeft("", highlightEnd-highlightStart, '~'), "lightred", false)
		fmt.Println(marginSpacer, fmt.Sprintf("%v%v", leadingSpace, arrows))
	}

	//Severity and description
	severityDescription, severityColor := "", ""

	switch entry.severity {
	case enumErrorCodes.Warning:
		severityColor = "yellow"
		severityDescription = "WARN"
	case enumErrorCodes.Error:
		severityColor = "red"
		severityDescription = "ERROR"
	case enumErrorCodes.Fatal:
		severityColor = "magenta"
		severityDescription = "FATAL ERROR"
	}
	colorizedSeverity := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %v ", severityDescription), 7, ' '), severityColor, true)
	fmt.Println("▓", colorizedSeverity, entry.message)

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
