package main

import (
	"bufio"
	"fmt"
	"log"
	"misc/nintasm/assemble"
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

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Failed to open file.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var file_lines []string
	for scanner.Scan() {
		file_lines = append(file_lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		log.Println("Failed to read line in file.")
		return
	}

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
	err = assemble.Start(file_lines)
	if err != nil {
		fmt.Println(err)
	}

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
	return
}
