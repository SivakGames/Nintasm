package fileHandler

import (
	"bufio"
	"io"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"os"
	"path/filepath"
)

var relativeFileDirectory string
var TriggerNewStackCall bool = false

// ---------------------------------------------------

// Get and process the primary source file from the command line
func GetFirstInputFile(inputFileName string) error {
	var err error

	_, err = os.Stat(inputFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return errorHandler.AddNew(enumErrorCodes.IncludeFileNotExist, inputFileName) // ❌ Fails
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
		return errorHandler.AddNew(enumErrorCodes.FailOpenFile) // ❌ Fails
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var processedLines []string
	for scanner.Scan() {
		processedLines = append(processedLines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return errorHandler.AddNew(enumErrorCodes.FailScanFile) // ❌ Fails
	}

	fileStack.PushToTopOfStack(inputFileName, processedLines)
	return nil
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
		return []byte{}, errorHandler.AddNew(enumErrorCodes.BinFileSeekAtEnd, seekPosition) // ❌ Fails

	} else if seekFileSizeDifference < 0 {
		return []byte{}, errorHandler.AddNew(enumErrorCodes.BinFileSeekAfterEnd, seekPosition, -seekFileSizeDifference) // ❌ Fails
	}
	_, err = file.Seek(seekPosition, 0)
	if err != nil {
		return []byte{}, err
	}

	if read > int(seekFileSizeDifference) {
		diff := read - int(seekFileSizeDifference)
		return []byte{}, errorHandler.AddNew(enumErrorCodes.BinFileReadBeyondFileSize, read, diff) // ❌ Fails
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
