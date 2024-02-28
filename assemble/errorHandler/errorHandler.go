package errorHandler

import (
	"errors"
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

func coloredSymbol(s string) string {
	return util.Colorize(s, "lightcyan", false)
}
func coloredNumber(s string) string {
	return util.Colorize(s, "yellow", false)
}

var errorTable = map[enumErrorCodes.Def]ErrorTableEntry{
	enumErrorCodes.IncludeFileNotExist: newErrorTableEntry(enumErrorCodes.Fatal, "Source file \x1b[92m%v\x1b[0m does not exist!"),
	enumErrorCodes.FailOpenFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to open source file: %v"),
	enumErrorCodes.FailScanFile:        newErrorTableEntry(enumErrorCodes.Fatal, "Failed to scan file!\n%v"),

	enumErrorCodes.BinFileSeekAtEnd:          newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d is at the very end of file so no bytes can be read!"),
	enumErrorCodes.BinFileSeekAfterEnd:       newErrorTableEntry(enumErrorCodes.Fatal, "Seek value of %d goes beyond the size of file by %d byte(s)"),
	enumErrorCodes.BinFileReadBeyondFileSize: newErrorTableEntry(enumErrorCodes.Fatal, "Read value of %d goes beyond the size of file by %d byte(s)"),

	enumErrorCodes.ParserEndOfInput:      newErrorTableEntry(enumErrorCodes.Error, "Parsing error - Unexpected end of input!"),
	enumErrorCodes.ParserUnexpectedToken: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Parsing error - Unexpected token: %v", coloredSymbol("%v"))),

	enumErrorCodes.OperationUNKNOWN:                 newErrorTableEntry(enumErrorCodes.Error, "UNKNOWN OPERATION"),
	enumErrorCodes.OperationBadTokenAfter:           newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after operation: %v"),
	enumErrorCodes.OperationDirectiveUnknown:        newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Unknown directive!"),
	enumErrorCodes.OperationEmpty:                   newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - UNEXPECTED EMPTY OPERATION"),
	enumErrorCodes.OperationLabelBadTokenAfter:      newErrorTableEntry(enumErrorCodes.Error, "ILLEGAL token(s) after labeled operation: %v"),
	enumErrorCodes.OperationLabelMissingColon:       newErrorTableEntry(enumErrorCodes.Error, "Incomplete label/labeled operation - either needs colon or directive"),
	enumErrorCodes.OperationLabeledDirectiveNoSpace: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Labeled directive needs space between label and directive"),
	enumErrorCodes.OperationLabeledDirectiveUnknown: newErrorTableEntry(enumErrorCodes.Error, "Operation parsing failed - Unknown labeled directive!"),

	enumErrorCodes.OperandListStartingComma: newErrorTableEntry(enumErrorCodes.Error, "Operand list cannot start with a comma!"),
	enumErrorCodes.OperandListTooMany:       newErrorTableEntry(enumErrorCodes.Error, "Too many operands for operation! Takes at most: %d"),
	enumErrorCodes.OperandListTooFew:        newErrorTableEntry(enumErrorCodes.Error, "Too few operands for operation! Needs at least: %d"),

	enumErrorCodes.NodeTypeNotNumeric:   newErrorTableEntry(enumErrorCodes.Error, "Value must be numeric!"),
	enumErrorCodes.NodeValueNotPositive: newErrorTableEntry(enumErrorCodes.Error, "Value must be positive!"),
	enumErrorCodes.NodeValueNot8Bit:     newErrorTableEntry(enumErrorCodes.Error, "Value must be 8 bit!"),
	enumErrorCodes.NodeValueNotPowerOf2: newErrorTableEntry(enumErrorCodes.Error, "Value must be a power of 2!"),

	enumErrorCodes.NodeValueNotGT:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol(">"), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotLT:        newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol("<"), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotGTE:       newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol(">="), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotLTE:       newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v", coloredSymbol("<="), coloredNumber("%d"))),
	enumErrorCodes.NodeValueNotGTEandLTE: newErrorTableEntry(enumErrorCodes.Error, fmt.Sprintf("Value must be %v %v and %v %v", coloredSymbol(">="), coloredNumber("%d"), coloredSymbol("<="), coloredNumber("%d"))),

	enumErrorCodes.InvalidValueAlias: newErrorTableEntry(enumErrorCodes.Error, "Invalid value alias!"),
	enumErrorCodes.UnacceptableAlias: newErrorTableEntry(enumErrorCodes.Error, "Unacceptable value alias!"),

	enumErrorCodes.INESValueAlreadySet: newErrorTableEntry(enumErrorCodes.Error, "%v value has already been set!!!"),

	enumErrorCodes.InstUnsupportedMode: newErrorTableEntry(enumErrorCodes.Error, "Mode is not supported by instruction!"),
	enumErrorCodes.InstXYUnusableMode:  newErrorTableEntry(enumErrorCodes.Error, "X or Y indexes cannot be used with target mode"),
	enumErrorCodes.InstBadAccumMode:    newErrorTableEntry(enumErrorCodes.Error, "No tokens can follow A for accumulator mode."),
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

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) error {
	errData, tableKeyExists := errorTable[errorTableKey]
	if tableKeyExists {
		errMsg := fmt.Sprintf(errData.description, args...)
		entry := NewErrorEntry(errorTableKey, errMsg, errData.severity)

		colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", entry.fileName), "red", true)
		fmt.Println("░", colorizedFileName)

		//Line number and content
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), 7, ' '), "blue", true)
		fmt.Println("▓", colorizedLineNumber, entry.lineContent)

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
		return errors.New(fmt.Sprintf("SEVERITY_%d", entry.severity))
	}
	return errors.New("Non-error-code error???")
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func IsErrorCode(err error) (bool, enumErrorCodes.Def) {
	errorCode := err.Error()
	errorTableKey := enumErrorCodes.Def(errorCode)
	_, isValidErrorCode := errorTable[errorTableKey]
	return isValidErrorCode, errorTableKey
}
