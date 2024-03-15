package fileStack

type fileStackEntry struct {
	FileName          string
	CurrentLineNumber uint
	ProcessedLines    []string
}

func newFileStackEntry(fileName string, processedLines []string) fileStackEntry {
	return fileStackEntry{
		FileName:          fileName,
		CurrentLineNumber: 0,
		ProcessedLines:    processedLines,
	}
}

type subOp struct {
	SubOpLineNumber uint
	LineContent     string
}

func newSubOp(lineNumber uint, lineContent string) subOp {
	return subOp{
		SubOpLineNumber: lineNumber,
		LineContent:     lineContent,
	}
}

var fileStackEntries []*fileStackEntry
var subOperation *subOp = nil

// +++++++++++++++++++++++++++++++++

func GetTopOfFileStack() *fileStackEntry {
	if len(fileStackEntries) > 0 {
		return fileStackEntries[len(fileStackEntries)-1]
	}
	return nil
}

func PushToTopOfStack(inputFileName string, processedLines []string) {
	entry := newFileStackEntry(inputFileName, processedLines)
	fileStackEntries = append(fileStackEntries, &entry)
}

func PopFromFileStack() {
	fileStackEntries = fileStackEntries[:len(fileStackEntries)-1]
}

// +++++++++++++++++++++++++++++++++

func AddSubOp(lineNumber uint, lineContent string) {
	s := newSubOp(lineNumber, lineContent)
	subOperation = &s
}

func ClearSubOp() {
	subOperation = nil
}

func GetSubOp() *subOp {
	return subOperation
}
