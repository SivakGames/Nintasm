package errorHandler

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/util"
)

type errorOutput struct {
	fileName     string
	errorEntries []ErrorEntry
}

func newErrorOutput(fileName string) errorOutput {
	return errorOutput{
		fileName: fileName,
	}
}

// =============================================================

func CheckAndPrintErrors() {
	if totalErrors == 0 {
		return
	}

	errorOutputs := make([]errorOutput, 0)
	fileNameOutputOrderIndexes := make(map[string]int)

	// Will first group by file

	for i := 0; i < int(totalErrors); i++ {
		entry := errorStack[i]
		_, exists := fileNameOutputOrderIndexes[entry.fileName]
		if !exists {
			fileNameOutputOrderIndexes[entry.fileName] = len(errorOutputs)
			errorOutputs = append(errorOutputs, newErrorOutput(entry.fileName))
		}
		index := fileNameOutputOrderIndexes[entry.fileName]
		errorOutputs[index].errorEntries = append(errorOutputs[index].errorEntries, entry)
	}

	for _, output := range errorOutputs {
		printErrorFileName(output.fileName)
		for _, entry := range output.errorEntries {
			printError(&entry)
		}
	}
}

func printErrorFileName(fileName string) {
	colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", fileName), "red", true)
	fmt.Println("░", colorizedFileName)
}

func printError(entry *ErrorEntry) {

	//Line number and content
	if entry.lineNumber > 0 {
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), ERROR_CAPTION_MIN_WIDTH, ' '), "blue", true)
		fmt.Println("▓", colorizedLineNumber, entry.lineContent)
	}

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
