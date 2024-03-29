package main

import (
	"flag"
	"fmt"
	"io"
	"misc/nintasm/assemble"
	"misc/nintasm/assemble/errorHandler"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/romBuilder"
	"misc/nintasm/util"
	"misc/nintasm/util/commandLine"
	"os"
	"time"
)

func main() {
	var err error

	util.ClearTerminal()

	util.DrawBox("Nintasm v1.000",
		enumTerminalColors.AnsiOrange,
		enumTerminalColors.LightYellow,
		enumTerminalColors.None,
		true)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename> [-s]")
		return
	}

	// ---------------------------------------------------------

	baseInputFileName := os.Args[1]
	rFlag := flag.Bool("r", false, "Raw rom file without INES Header")
	sFlag := flag.Bool("s", false, "Show segment usage")

	flag.CommandLine.SetOutput(io.Discard)
	flag.CommandLine.Parse(os.Args[2:])

	commandLine.SettingRawROMFile = *rFlag
	commandLine.SettingShowSegmentUsage = *sFlag

	start := time.Now()
	err = assemble.Start(baseInputFileName)
	if err != nil {
		errorHandler.CheckAndPrintErrors()
		errorHandler.PrintTotalErrorMessage()
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

	fmt.Println()
	if commandLine.SettingShowSegmentUsage {
		romBuilder.OutputSegmentUsage()
	}
	fmt.Println("Assembly has succeeded!")
	fmt.Println("Output file can be found at:", util.Colorize(outFileName, enumTerminalColors.AnsiSeaGreen, false))
	fmt.Println()

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
	return
}
