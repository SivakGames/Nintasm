package assemble

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/fileHandler"
	"misc/nintasm/assemble/fileStack"
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	"misc/nintasm/interpreter"
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

	err = fileHandler.GetFirstInputFile(initialInputFile)
	if err != nil {
		return err
	}
	err = startReadingLinesTopFileStack()
	if err != nil {
		return err
	}

	unresolvedTable.ResolvedUnresolved()

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

		//Intermediate - determine if in stack capturing
		if len(blockStack.Stack) > 0 {
			var blockStackErr error
			isStartEndOp := blockStack.CheckIfNewStartEndOperation(&lineOperationParsedValues)

			//See if the incoming operation is for starting/ending a block
			if isStartEndOp {
				blockStackErr = parseOperandStringAndProcess(
					reformattedLine,
					&lineOperationParsedValues,
				)
				if blockStackErr != nil {
					return blockStackErr
				}

				//If ending, iterate bottom of stack and parse all captured operations (if any)
				if blockStack.CheckIfEndOperationAndClearStack(&lineOperationParsedValues) {
					for _, b := range blockStack.Stack[0].CapturedLines {
						processOperandArguments := util.NewLineOperationParsedValues(b.OperandStartPosition,
							b.OperationLabel,
							b.OperationTokenEnum,
							b.OperationTokenValue,
							b.ParentParserEnum,
						)
						blockStackErr = parseOperandStringAndProcess(
							b.OriginalLine,
							&processOperandArguments,
						)
						if blockStackErr != nil {
							return blockStackErr
						}
					}
					//Mainly set by namespaces - will clear the overriding parent label
					if interpreter.PopParentLabelWhenBlockOpDone {
						interpreter.PopParentLabel()
						interpreter.PopParentLabelWhenBlockOpDone = false
					}
					blockStack.ClearStack()
				}

			} else {
				//Either append the operation to the stack's captured lines or evaluate them now
				if !blockStack.GetCurrentOperationEvaluatesCapturedNodesFlag() {
					err := blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, &lineOperationParsedValues)
					if err != nil {
						return err
					}
				} else {
					blockStackErr = parseOperandStringAndProcess(
						reformattedLine,
						&lineOperationParsedValues,
					)
					if blockStackErr != nil {
						return blockStackErr
					}
				}
			}
			continue
		}

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
