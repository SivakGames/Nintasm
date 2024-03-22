package main

import (
	"fmt"
	"misc/nintasm/assemble"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/romBuilder"
	"os"
	"time"
)

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename> [-s]")
		return
	}

	// --------------------

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
		fmt.Println("\x1b[31mERRORZ HAPPENED\x1b[0m")
		errorHandler.CheckAndPrintErrors()
	}

	romBuilder.GenerateRomBinFile()

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
	return
}
