package logHandler

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
)

// ++++++++++++++++++++++++++++++++++++++++

var totalQueuedLogs uint = 0

const LOG_STACK_CAPACITY = 500

type LogEntry struct {
	fileName    string
	lineNumber  uint
	lineContent string
	message     string
}

var logStack = make([]LogEntry, LOG_STACK_CAPACITY)

func newLogEntry(
	message string,
) LogEntry {
	entry := LogEntry{
		message:     message,
		fileName:    "???",
		lineNumber:  0,
		lineContent: "",
	}
	fileData := fileStack.GetTopOfFileStack()
	if fileData != nil {
		entry.fileName = fileData.FileName
		entry.lineNumber = fileData.CurrentLineNumber
		entry.lineContent = fileData.ProcessedLines[fileData.CurrentLineNumber-1]
	}
	return entry
}

// ----------------------------------------

func AddLog(message string) {
	logEntry := newLogEntry(message)
	logStack[totalQueuedLogs] = logEntry
	totalQueuedLogs++
}

func OutputLogs() {
	if totalQueuedLogs > 0 {
		fmt.Println()
		fmt.Println(util.Colorize(" --- START LOG OUTPUT --- ", enumTerminalColors.AnsiDarkSeaGreen, true))
	}
	for i := uint(0); i < totalQueuedLogs; i++ {
		colorizedFileName := util.Colorize(fmt.Sprintf(" >> %v ", logStack[i].fileName), enumTerminalColors.AnsiGray2, true)
		fmt.Println("░", colorizedFileName)
		colorizedLineNumber := util.Colorize(util.PadStringLeft(fmt.Sprintf(" %d ", logStack[i].lineNumber), 7, ' '), enumTerminalColors.AnsiGray3, true)
		fmt.Println("▓", colorizedLineNumber, logStack[i].message)
	}
	if totalQueuedLogs > 0 {
		fmt.Println(util.Colorize(" ---- END LOG OUTPUT ---- ", enumTerminalColors.AnsiDarkSeaGreen, true))
		fmt.Println()
	}
}
