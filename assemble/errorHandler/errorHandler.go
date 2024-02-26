package errorHandler

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
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

type ErrorEntry struct {
	code       enumErrorCodes.Def
	lineNumber uint
	fileName   string
	message    string
	hint       string
}

func NewErrorEntry(code enumErrorCodes.Def) ErrorEntry {
	fileData := fileStack.GetTopOfFileStack()
	return ErrorEntry{
		code:       code,
		fileName:   fileData.FileName,
		lineNumber: fileData.CurrentLineNumber,
	}
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func AddNew(errorTableKey enumErrorCodes.Def, args ...interface{}) {
	errData, ok := errorTable[errorTableKey]
	if ok {
		errMsg := fmt.Sprintf(errData.description, args...)
		fmt.Println(errMsg)
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
