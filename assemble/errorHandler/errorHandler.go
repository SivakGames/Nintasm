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
var highlightStart int = -1
var highlightEnd int = -1
var TotalErrors uint = 0

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

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddHighlights(start int, end int) {
	highlightStart = start
	highlightEnd = end
}

func Resetighlights() {
	highlightEnd = -1
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) error {
	defer Resetighlights()

	TotalErrors++

	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry := NewErrorEntry(errorTableKey, errMsg, errData.severity)

		//File name (if any)
		colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", entry.fileName), "red", true)
		fmt.Println("░", colorizedFileName)

		//Line number and content
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), ERROR_CAPTION_MIN_WIDTH, ' '), "blue", true)
		fmt.Println("▓", colorizedLineNumber, entry.lineContent)

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
		fmt.Println("▓", colorizedSeverity, errMsg)
		return errors.New(fmt.Sprintf("%v%d", SEVERITY_PREFIX, entry.severity))
	}
	return errors.New("Non-error-code error???")
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
