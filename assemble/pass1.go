package assemble

import (
	enumParserTypes "misc/nintasm/enums/parserTypes"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/parser"
	"misc/nintasm/util"
)

var directiveOperandParser = parser.NewDirectiveOperandParser()
var instructionOperandParser = parser.NewInstructionOperandParser()
var labelOperandParser = parser.NewLabelOperandParser()

func Start(lines []string) error {
	instructionOperandParser.ShouldParseInstructions = true

	lineInitParser := parser.NewInitialLineParser()
	lineOperationParser := parser.NewOperationParser()

	var lineCounter uint = 0

	// Iterate over all lines

	for _, rawLine := range lines {
		lineCounter++

		//Step 1 - Reformat line
		reformattedLine, lineInitErr := lineInitParser.Process(rawLine)
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
			if isStartEndOp {
				blockStackErr = parseOperandStringAndProcess(
					reformattedLine,
					&lineOperationParsedValues,
				)
				if blockStackErr != nil {
					return blockStackErr
				}
				if blockStack.CheckIfEndOperationAndClearStack(&lineOperationParsedValues) {
					for _, b := range blockStack.Stack[0].CapturedLines {
						temp := util.NewLineOperationParsedValues(b.OperandStartPosition,
							b.OperationLabel,
							b.OperationTokenEnum,
							b.OperationTokenValue,
							b.ParentParserEnum,
						)
						blockStackErr = parseOperandStringAndProcess(
							b.OriginalLine,
							&temp,
						)
						if blockStackErr != nil {
							return blockStackErr
						}
					}
					blockStack.ClearStack()
				}
			} else {
				err := blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, &lineOperationParsedValues)
				if err != nil {
					return err
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

	//	case enumParserTypes.Macro:
	//		fmt.Println("Mack")
	//		fmt.Println(optype, opval, opPos)

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

	default:
		panic("🛑 Parent parsing operation could not be determined!")
	}
	return nil
}
