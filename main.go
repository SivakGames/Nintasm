package main

import (
	"fmt"
	"misc/nintasm/assemble"
	"misc/nintasm/assemble/errorHandler"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/romBuilder"
	"misc/nintasm/util"
	"os"
	"time"
)

func main() {
	var err error

	util.ClearTerminal()

	fmt.Println("Nintasm")

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename> [-s]")
		return
	}

	// ---------------------------------------------------------

	baseInputFileName := os.Args[1]

	//	sFlag := flag.Bool("s", false, "A S boolean flag")
	//	rFlag := flag.Bool("r", false, "A R boolean flag")
	//
	//	flag.CommandLine.SetOutput(ioutil.Discard)
	//	err = flag.CommandLine.Parse(os.Args[2:])
	//
	//	fmt.Println("File:", filename)
	//	fmt.Println("Command:", *sFlag)
	//	fmt.Println("Command:", *rFlag)

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
	romBuilder.OutputSegmentUsage()
	fmt.Println("Assembly has succeeded!")
	fmt.Println("Output file can be found at:", util.Colorize(outFileName, enumTerminalColors.AnsiSeaGreen, false))
	fmt.Println()

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
	return
}
