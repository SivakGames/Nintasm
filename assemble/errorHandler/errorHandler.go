package errorHandler

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/util"
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
	enumErrorCodes.IncludeFileNotExist: newErrorTableEntry(enumErrorCodes.Fatal, "Source file \x1b[92m%v\x1b[0m does not exist!"),
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
		code:        code,
		message:     message,
		fileName:    "NO FILE",
		lineNumber:  0,
		lineContent: "",
		severity:    severity,
	}
}

/*

░ >> D:\Emulate\NES\Disassemblies\Lolo 3\prg\fixed.6502
▓   2711   .include "prg/music-engine/dpcm-samples.6502a"
▓  FATAL ERROR  Source  .INCLUDE  file prg/music-engine/dpcm-samples.6502a does not exist!!!

 >>> Assembly WILL NOT continue due to fatal errors! <<<
Assembly could not be completed due to errors!
Total Error Count: 1 / Total Warning Count: 0

*/

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) {
	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry := NewErrorEntry(errorTableKey, errMsg, errData.severity)

		colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", entry.fileName), "red", true)
		fmt.Println("░", colorizedFileName)

		//Line number and content
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), 7, ' '), "blue", true)
		fmt.Println("▓", colorizedLineNumber)

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

	}
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
