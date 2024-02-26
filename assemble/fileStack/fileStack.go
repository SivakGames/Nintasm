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

var InputFileLines []fileStackEntry

// +++++++++++++++++++++++++++++++++

func GetTopOfFileStack() fileStackEntry {
	return InputFileLines[len(InputFileLines)-1]
}

func PushToTopOfStack(inputFileName string, processedLines []string) {
	InputFileLines = append(InputFileLines, newFileStackEntry(inputFileName, processedLines))
}

func PopFromFileStack() {
	InputFileLines = InputFileLines[:len(InputFileLines)-1]
}
