package errorHandler

import (
	"fmt"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
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
	if totalQueuedErrors == 0 {
		return
	}

	errorOutputs := make([]errorOutput, 0)
	fileNameOutputOrderIndexes := make(map[string]int)

	// Will first group by file

	for i := 0; i < int(totalQueuedErrors); i++ {
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
	colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", fileName), enumTerminalColors.Red, true)
	fmt.Println("░", colorizedFileName)
}

func printError(entry *ErrorEntry) {

	//Line number and content
	if entry.lineNumber > 0 {
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.lineNumber), ERROR_CAPTION_MIN_WIDTH, ' '), enumTerminalColors.Blue, true)
		fmt.Println("▓", colorizedLineNumber, entry.lineContent)
	}

	if entry.subOpLineNumber > 0 {
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", entry.subOpLineNumber), ERROR_CAPTION_MIN_WIDTH-2, ' '), enumTerminalColors.Blue, true)
		fmt.Println("▓▓▓", colorizedLineNumber, entry.subOpLineContent)
	}

	if highlightEnd > highlightStart {
		marginSpacer := util.Colorize(util.PadStringLeft("", ERROR_CAPTION_MIN_WIDTH+2, '░'), enumTerminalColors.Cyan, false)
		leadingSpace := util.PadStringLeft("", highlightStart, ' ')
		arrows := util.Colorize(util.PadStringLeft("", highlightEnd-highlightStart, '~'), enumTerminalColors.LightRed, false)
		fmt.Println(marginSpacer, fmt.Sprintf("%v%v", leadingSpace, arrows))
	}

	//Severity and description
	severityDescription, severityColor := "", enumTerminalColors.Red

	switch entry.severity {
	case enumErrorCodes.Warning:
		severityColor = enumTerminalColors.Yellow
		severityDescription = "WARN"
	case enumErrorCodes.Error:
		severityColor = enumTerminalColors.Red
		severityDescription = "ERROR"
	case enumErrorCodes.Fatal:
		severityColor = enumTerminalColors.AnsiRed
		severityDescription = "FATAL ERROR"
	}
	colorizedSeverity := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %v ", severityDescription), ERROR_CAPTION_MIN_WIDTH, ' '), severityColor, true)
	fmt.Println("▓", colorizedSeverity, entry.message)

	if entry.hint != "" {
		colorizedHint := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %v ", "HINT"), ERROR_CAPTION_MIN_WIDTH, ' '), enumTerminalColors.AnsiGray3, false)
		fmt.Println("▓", colorizedHint, entry.hint)
	}

}

func PrintTotalErrorMessage() {
	fmt.Println()
	fmt.Println(util.Colorize("Assembly has terminated!", enumTerminalColors.AnsiOrange, false))
	totalErrorText := fmt.Sprintf("Total errors: %d", totalErrors)
	totalWarningText := fmt.Sprintf("Total warnings: %d", totalWarnings)
	fmt.Println(fmt.Sprintf("%v / %v", totalErrorText, totalWarningText))
}
