package errorHandler

import enumErrorCodes "misc/nintasm/constants/enums/errorCodes"

type ErrorEntry struct {
	code       enumErrorCodes.Def
	lineNumber int
	fileName   string
	message    string
	hint       string
}

func NewErrorEntry(code enumErrorCodes.Def, lineNumber int) ErrorEntry {
	return ErrorEntry{}
}
