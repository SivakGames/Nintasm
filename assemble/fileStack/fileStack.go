package fileStack

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

var InputFileLines []fileStackEntry
var relativeFileDirectory string
var TriggerNewStackCall bool = false

// The primary source file from the command line
func GetFirstInputFile(inputFileName string) error {
	var err error

	_, err = os.Stat(inputFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Source file does not exist!")
		} else {
			return err
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	absFilePath := filepath.Join(cwd, inputFileName)
	relativeFileDirectory = filepath.Dir(absFilePath)

	err = OpenInputFileAndPushLinesToStack(absFilePath)
	if err != nil {
		return err
	}
	return nil
}

// Ensure newly included file names have the complete path preceding them
func AddRelativePathIncludeFileName(inputFileName string) string {
	if len(inputFileName) > 0 && inputFileName[0] != '/' {
		inputFileName = "/" + inputFileName
	}
	inputFileName = relativeFileDirectory + inputFileName
	return inputFileName
}

// Open new include file
func OpenInputFileAndPushLinesToStack(inputFileName string) error {
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

	InputFileLines = append(InputFileLines, newFileStackEntry(inputFileName, processedLines))
	return nil
}

// +++++++++++++++++++++++++++++++++

func GetTopOfFileStack() fileStackEntry {
	return InputFileLines[len(InputFileLines)-1]
}

func PopFromFileStack() {
	InputFileLines = InputFileLines[:len(InputFileLines)-1]
}

// --------------------------------

// For incbin directives
func ProcessBinFile(binFileName string, seek int, read int) ([]byte, error) {
	var buffer []byte

	file, err := os.Open(binFileName)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return []byte{}, err
	}

	seekPosition := int64(seek)
	seekFileSizeDifference := fileInfo.Size() - seekPosition

	if seekFileSizeDifference == 0 {
		errMsg := fmt.Sprintf("Seek value of %d is at the very end of file so no bytes can be read!", seekPosition)
		return []byte{}, errors.New(errMsg)
	} else if seekFileSizeDifference < 0 {
		errMsg := fmt.Sprintf("Seek value of %d goes beyond the size of file by %d byte(s)", seekPosition, -seekFileSizeDifference)
		return []byte{}, errors.New(errMsg)
	}
	_, err = file.Seek(seekPosition, 0)
	if err != nil {
		return []byte{}, err
	}

	if read > int(seekFileSizeDifference) {
		diff := read - int(seekFileSizeDifference)
		errMsg := fmt.Sprintf("Read value of %d goes beyond the size of file by %d byte(s)", read, diff)
		return []byte{}, errors.New(errMsg)
	}

	if read != 0 {
		buffer = make([]byte, read)
	} else {
		buffer = make([]byte, seekFileSizeDifference)
	}

	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}
