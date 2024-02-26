package errorHandler

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"strings"
)

type ErrorTableEntry struct {
	description string
	severity    enumErrorCodes.Severity
}

func newErrorTableEntry(severity enumErrorCodes.Severity, description string) ErrorTableEntry {
	return ErrorTableEntry{
		description: description,
		severity:    severity,
	}
}

var errorTable = map[enumErrorCodes.Def]ErrorTableEntry{
	enumErrorCodes.IncludeFileNotExist: newErrorTableEntry(enumErrorCodes.Fatal, "Source file %v does not exist!"),
	enumErrorCodes.FailOpenFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to open source file: %v"),
	enumErrorCodes.FailScanFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to scan file!\n%v"),

	enumErrorCodes.INESPRGSet:          newErrorTableEntry(enumErrorCodes.Error, "INES PRG has already been set!!!"),
	enumErrorCodes.INESPRGBadValue:     newErrorTableEntry(enumErrorCodes.Error, "INES PRG must be >= 1 or use a valid alias"),
	enumErrorCodes.INESPRGUnacceptable: newErrorTableEntry(enumErrorCodes.Error, "Unacceptable INES PRG size declared!"),
	enumErrorCodes.INESCHRSet:          newErrorTableEntry(enumErrorCodes.Error, "INES CHR has already been set!!!"),
	enumErrorCodes.INESCHRBadValue:     newErrorTableEntry(enumErrorCodes.Error, "INES CHR must be >= 1 or use a valid alias"),
	enumErrorCodes.INESCHRUnacceptable: newErrorTableEntry(enumErrorCodes.Error, "Unacceptable INES CHR size declared!"),
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
		code:       code,
		message:    message,
		fileName:   "NO FILE",
		lineNumber: 0,
	}
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) {
	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry := NewErrorEntry(errorTableKey, errMsg, errData.severity)

		fmt.Println("\x1b[41m", entry.fileName, "\x1b[0m")
		if entry.lineNumber > 0 {
			lineData := fmt.Sprintf("%d %v", entry.lineNumber, entry.lineContent)
			fmt.Println(lineData)
		}
		paddedStr := padStringLeft("1", 5, '░')
		fmt.Println(paddedStr)

	}
}

func padStringLeft(s string, length int, char rune) string {
	padding := length - len(s)
	if padding <= 0 {
		return s // No padding needed or negative padding
	}
	return strings.Repeat(string(char), padding) + s
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func ProcessError(err error) {
	errorCode := err.Error()
	errorTableKey := enumErrorCodes.Def(errorCode)
	errData, ok := errorTable[errorTableKey]
	if ok {
		fmt.Println(errData)

	} else {
		fmt.Println("Non-coded error has occurred!")
	}
}
