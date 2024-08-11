package main

import (
	"flag"
	"fmt"
	"io"
	"misc/nintasm/assemble"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/logHandler"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/romBuilder"
	"misc/nintasm/util"
	"misc/nintasm/util/commandLine"
	"os"
	"time"
)

const VERSION_NUMBER = "1.000"
const BUILD_DATE = "24.08.11"

func main() {
	var err error

	util.ClearTerminal()

	util.DrawBox(fmt.Sprintf("Nintasm %s / Rev %s", VERSION_NUMBER, BUILD_DATE),
		enumTerminalColors.AnsiOrange,
		enumTerminalColors.LightYellow,
		enumTerminalColors.None,
		false)

	if len(os.Args) < 2 {
		fmt.Println()
		fmt.Println(util.Colorize("*** No input file detected ***", enumTerminalColors.AnsiGray4, false))
		fmt.Println()
		fmt.Println(util.Colorize("Usage:", enumTerminalColors.Cyan, false), "nintasm <filename> [-arguments]")
		fmt.Println(util.Colorize("For help:", enumTerminalColors.Cyan, false), "nintasm -h")
		return
	}

	// ---------------------------------------------------------

	baseInputFileName := os.Args[1]
	hFlag := *(flag.Bool("h", false, "Display help")) || baseInputFileName == "-h"
	oFlag := flag.String("o", "", "Change output file name (path relative to input file name)")
	rFlag := flag.Bool("r", false, "Raw rom file without INES Header")
	sFlag := flag.Bool("s", false, "Show segment usage")

	flag.CommandLine.SetOutput(io.Discard)
	err = flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		fmt.Println("Errors while parsing command line")
		fmt.Println(err)
		return
	}

	if hFlag {
		util.ClearTerminal()
		util.DrawBox(fmt.Sprintf("Nintasm Help"),
			enumTerminalColors.AnsiGray3,
			enumTerminalColors.Cyan,
			enumTerminalColors.None,
			false)

		fmt.Println("Command Line Arguments")
		fmt.Println(util.Colorize("-o", enumTerminalColors.Green, false), flag.Lookup("o").Usage)
		fmt.Println(util.Colorize("-r", enumTerminalColors.Green, false), flag.Lookup("r").Usage)
		fmt.Println(util.Colorize("-s", enumTerminalColors.Green, false), flag.Lookup("s").Usage)

		return
	}

	commandLine.SettingOverwriteOutFileName = *oFlag
	commandLine.SettingRawROMFile = *rFlag
	commandLine.SettingShowSegmentUsage = *sFlag

	start := time.Now()
	err = assemble.Start(baseInputFileName)
	if err != nil {
		errorHandler.CheckAndPrintErrors()
		errorHandler.PrintTotalErrorMessage(false)
		logHandler.OutputLogs()
		return // ❌ Assembly Fails
	}

	// ---------------------------------------------------------
	// Generate rom file

	outFileName, romBuildErr := romBuilder.GenerateRomBinFile()
	if romBuildErr != nil {
		fmt.Println("No errors with code, but output file could not be generated...")
		fmt.Println(romBuildErr)
		return
	}

	// ---------------------------------------------------------
	// 🟢 Assembly Succeeds!

	hasShownWarnings := errorHandler.GetWarningCount() > 0
	if hasShownWarnings {
		errorHandler.CheckAndPrintErrors()
		errorHandler.PrintTotalErrorMessage(true)
	}

	fmt.Println()
	logHandler.OutputLogs()

	if commandLine.SettingShowSegmentUsage {
		romBuilder.OutputSegmentUsage()
	}
	if !hasShownWarnings {
		fmt.Println(util.Colorize("Assembly has succeeded!", enumTerminalColors.AnsiGreen, false))
	}
	fmt.Println("Output file can be found at:", util.Colorize(outFileName, enumTerminalColors.AnsiSeaGreen, false))
	if commandLine.SettingRawROMFile {
		fmt.Println(util.Colorize("*Note:", enumTerminalColors.Yellow, false), "Raw file flag detected so no iNES header generated")
	}
	fmt.Println()

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
	return
}
