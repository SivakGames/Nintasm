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

var fileStackEntries []fileStackEntry

// +++++++++++++++++++++++++++++++++

func GetTopOfFileStack() *fileStackEntry {
	if len(fileStackEntries) > 0 {
		return &fileStackEntries[len(fileStackEntries)-1]
	}
	return nil
}

func PushToTopOfStack(inputFileName string, processedLines []string) {
	fileStackEntries = append(fileStackEntries, newFileStackEntry(inputFileName, processedLines))
}

func PopFromFileStack() {
	fileStackEntries = fileStackEntries[:len(fileStackEntries)-1]
}
