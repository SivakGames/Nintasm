package fileHandler

import (
	"bufio"
	"fmt"
	"io"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
	"misc/nintasm/util/commandLine"
	"os"
	"path/filepath"
	"strings"
)

var relativeFileDirectory string
var TriggerNewStackCall bool = false
var outFileName string = ""

// ---------------------------------------------------

// Get and process the primary source file from the command line
func GetFirstInputFile(inputFileName string) error {
	var err error

	// Get the directory part of the input file path
	dir, base, inFileExtension := filepath.Dir(inputFileName), filepath.Base(inputFileName), filepath.Ext(inputFileName)

	// Check if the directory exists
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return errorHandler.AddNew(enumErrorCodes.IncludeFileNotExist, inputFileName) //❌☠️ FATAL ERROR
	} else if err != nil {
		return errorHandler.AddNew(enumErrorCodes.OtherFatal, err) //❌☠️ FATAL ERROR
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	fmt.Println(util.Colorize("▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁", enumTerminalColors.AnsiGray5, false))
	fmt.Println(fmt.Sprintf(" Attempting to assemble: %v", util.Colorize(base, enumTerminalColors.AnsiGreen, false)))
	fmt.Println(util.Colorize("▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔", enumTerminalColors.AnsiGray5, false))

	relativeFileDirectory = absPath
	err = OpenInputFileAndPushLinesToStack(relativeFileDirectory + "/" + base)
	if err != nil {
		return err
	}

	if commandLine.SettingOverwriteOutFileName != "" {
		outFileName = commandLine.SettingOverwriteOutFileName
	} else {
		outFileName = strings.TrimSuffix(base, inFileExtension) + ".nes"

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
		return errorHandler.AddNew(enumErrorCodes.FailOpenFile, inputFileName) // ❌ Fails
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

func GenerateOutFileName() string {
	return AddRelativePathIncludeFileName(outFileName)
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
