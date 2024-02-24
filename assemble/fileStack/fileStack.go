package fileStack

import (
	"bufio"
	"errors"
	"os"
)

var InputFileLines [][]string

func GetFirstInputFile(inputFileName string) error {
	err := openInputFileAndAppendLinesToStack(inputFileName)
	if err != nil {
		return err
	}

	//Add Code for making relative file path

	return nil
}

func openInputFileAndAppendLinesToStack(inputFileName string) error {
	file, err := os.Open(inputFileName)
	if err != nil {
		return errors.New("Failed to open file.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var processedLines []string
	for scanner.Scan() {
		processedLines = append(processedLines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return errors.New("Failed to read line in file.")
	}

	InputFileLines = append(InputFileLines, processedLines)
	return nil
}

func GetTopOfFileStack() []string {
	return InputFileLines[len(InputFileLines)-1]
}
