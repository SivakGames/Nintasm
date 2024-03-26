package assemble

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/predefSymbols"
	"misc/nintasm/interpreter/environment/unresolvedTable"
	"misc/nintasm/parser"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var lineInitParser = parser.NewInitialLineParser()
var lineOperationParser = parser.NewOperationParser()

// ============================================================

// Main process starts - open input primary input file
func Start(initialInputFile string) error {
	var err error

	predefSymbols.AddPregensToMacroTable()
	err = fileHandler.GetFirstInputFile(initialInputFile)
	if err != nil {
		return err
	}

	fmt.Println("=========Pass 1=========")

	err = startReadingLinesTopFileStack()
	if err != nil {
		return err
	}

	if errorHandler.GetErrorCount() > 0 {
		return errors.New("Pass 1 failed")
	}

	fmt.Println("=========Pass 2=========")

	err = unresolvedTable.ResolvedUnresolvedSymbols()
	if err != nil {
		return err
	}
	err = unresolvedTable.ResolvedUnresolvedRomEntries()
	if err != nil {
		return err
	}

	if errorHandler.GetErrorCount() > 0 {
		return errors.New("Pass 2 failed")
	}

	return nil
}

// Will get whatever's on the top of the file stack and prep the loop
func startReadingLinesTopFileStack() error {
	fileData := fileStack.GetTopOfFileStack()
	err := ReadLines(&fileData.ProcessedLines, &fileData.CurrentLineNumber)
	if err != nil {
		return err
	}
	fileStack.PopFromFileStack()
	return nil
}

func ReadLines(lines *[]string, lineCounter *uint) error {
	// Iterate over all lines
	for i, rawLine := range *lines {
		*lineCounter += 1

		//Step 1 - Reformat line
		reformattedLine, lineInitErr := lineInitParser.Process(rawLine)
		(*lines)[i] = reformattedLine

		if lineInitErr != nil {
			return lineInitErr
		}
		if len(reformattedLine) == 0 {
			continue
		}

		//Step 2 - determine line op
		lineOperationErr := lineOperationParser.Process(reformattedLine)
		if lineOperationErr != nil {
			lineOperationErr := errorHandler.CheckErrorContinuesUpwardPropagation(lineOperationErr, enumErrorCodes.Error)
			if lineOperationErr != nil {
				return lineOperationErr // ❌❌ CONTINUES Failing!
			}
			continue
		}

		lineOperationParsedValues := lineOperationParser.GetLineOperationValues()

		//Intermediate - determine if capturing things in a block stack
		currentBlockStack := blockStack.GetCurrentCaptureBlockStack()
		if len(*currentBlockStack) > 0 {
			err := handleBlockStack(reformattedLine, &lineOperationParsedValues)
			if err != nil {
				return err
			}
			continue
		}

		//Do regular operand parsing/processing
		err := parseAndProcessOperandString(
			reformattedLine,
			&lineOperationParsedValues,
		)
		if err != nil {
			return err
		}

		// See if a new source file was opened via an include directive
		if fileHandler.TriggerNewStackCall {
			fileHandler.TriggerNewStackCall = false
			err := startReadingLinesTopFileStack()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
