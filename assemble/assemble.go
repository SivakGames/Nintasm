package assemble

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/fileHandler"
	"misc/nintasm/assemble/fileStack"
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	"misc/nintasm/interpreter/environment/predefSymbols"
	"misc/nintasm/interpreter/environment/unresolvedTable"
	"misc/nintasm/parser"
	"misc/nintasm/util"
)

var directiveOperandParser = parser.NewDirectiveOperandParser()
var instructionOperandParser = parser.NewInstructionOperandParser()
var labelOperandParser = parser.NewLabelOperandParser()
var macroOperandParser = parser.NewMacroOperandParser()

// Main process starts - open input primary input file
func Start(initialInputFile string) error {
	var err error

	predefSymbols.AddPregensToMacroTable()

	fmt.Println("=========Pass 1=========")

	err = fileHandler.GetFirstInputFile(initialInputFile)
	if err != nil {
		return err
	}
	err = startReadingLinesTopFileStack()
	if err != nil {
		return err
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
	lineInitParser := parser.NewInitialLineParser()
	lineOperationParser := parser.NewOperationParser()

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
			return lineOperationErr
		}

		lineOperationParsedValues := lineOperationParser.GetLineOperationValues()

		//Intermediate - determine if capturing things in a block stack

		currentBlockStack := blockStack.GetCurrentStack()
		if len(*currentBlockStack) > 0 {
			err := handleBlockStack(reformattedLine, &lineOperationParsedValues)
			if err != nil {
				return err
			}
			continue
		}

		//Do regular operand parsing/processing

		err := parseOperandStringAndProcess(
			reformattedLine,
			&lineOperationParsedValues,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

//==========================================================

// Main operand parsing...
func parseOperandStringAndProcess(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {

	switch lineOperationParsedValues.ParentParserEnum {

	// -------------------
	case enumParserTypes.Instruction:
		operandParserErr := instructionOperandParser.SetupOperandParser(
			reformattedLine,
			lineOperationParsedValues.OperandStartPosition,
		)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = instructionOperandParser.Process(lineOperationParsedValues.OperationTokenValue)
		if operandParserErr != nil {
			return operandParserErr
		}

	// -------------------
	case enumParserTypes.Directive:
		operandParserErr := directiveOperandParser.SetupOperandParser(
			reformattedLine,
			lineOperationParsedValues.OperandStartPosition,
		)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = directiveOperandParser.Process(
			lineOperationParsedValues.OperationTokenEnum,
			lineOperationParsedValues.OperationTokenValue,
			lineOperationParsedValues.OperationLabel,
		)
		if operandParserErr != nil {
			return operandParserErr
		}

	// -------------------
	case enumParserTypes.Label:
		operandParserErr := labelOperandParser.SetupOperandParser(
			reformattedLine,
			lineOperationParsedValues.OperandStartPosition,
		)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = labelOperandParser.Process(
			lineOperationParsedValues.OperationTokenEnum,
			lineOperationParsedValues.OperationTokenValue,
			lineOperationParsedValues.OperationLabel,
		)
		if operandParserErr != nil {
			return operandParserErr
		}

	// -------------------
	case enumParserTypes.Macro:
		operandParserErr := macroOperandParser.SetupOperandParser(
			reformattedLine,
			lineOperationParsedValues.OperandStartPosition,
		)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = macroOperandParser.Process(lineOperationParsedValues.OperationTokenValue)
		if operandParserErr != nil {
			return operandParserErr
		}
		linesToUnpack := macroOperandParser.GetUnpackLinesRef()
		for i := range *linesToUnpack {
			replacedCapturedLine := macroOperandParser.ApplyReplacementsToCapturedLine(i)
			temp := util.NewLineOperationParsedValues(replacedCapturedLine.OperandStartPosition,
				replacedCapturedLine.OperationLabel,
				replacedCapturedLine.OperationTokenEnum,
				replacedCapturedLine.OperationTokenValue,
				replacedCapturedLine.ParentParserEnum,
			)
			err := parseOperandStringAndProcess(replacedCapturedLine.OriginalLine, &temp)
			if err != nil {
				return err
			}
		}
		macroOperandParser.PopFromStack()

	default:
		panic("🛑 Parent parsing operation could not be determined!")
	}
	if fileHandler.TriggerNewStackCall {
		fileHandler.TriggerNewStackCall = false
		err := startReadingLinesTopFileStack()
		if err != nil {
			return err
		}
	}

	return nil
}
