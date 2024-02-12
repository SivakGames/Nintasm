package assemble

import (
	enumParserTypes "misc/nintasm/enums/parserTypes"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/parser"
	"misc/nintasm/util"
)

var directiveOperandParzival = parser.NewDirectiveOperandParser()
var instructionOperandParzival = parser.NewInstructionOperandParser()
var labelOperandParzival = parser.NewLabelOperandParser()

func Start(lines []string) error {
	instructionOperandParzival.ShouldParseInstructions = true

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
	case enumParserTypes.Directive:
		operandParserErr := directiveOperandParzival.SetupOperandParser(reformattedLine, lineOperationParsedValues.OperandStartPosition)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = directiveOperandParzival.Process(lineOperationParsedValues.OperationTokenEnum, lineOperationParsedValues.OperationTokenValue)
		if operandParserErr != nil {
			return operandParserErr
		}

	case enumParserTypes.Instruction:
		operandParserErr := instructionOperandParzival.SetupOperandParser(reformattedLine, lineOperationParsedValues.OperandStartPosition)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = instructionOperandParzival.Process(lineOperationParsedValues.OperationTokenValue)
		if operandParserErr != nil {
			return operandParserErr
		}

	//	case enumParserTypes.Macro:
	//		fmt.Println("Mack")
	//		fmt.Println(optype, opval, opPos)
	case enumParserTypes.Label:
		operandParserErr := labelOperandParzival.SetupOperandParser(reformattedLine, lineOperationParsedValues.OperandStartPosition)
		if operandParserErr != nil {
			return operandParserErr
		}
		operandParserErr = labelOperandParzival.Process(lineOperationParsedValues.OperationTokenEnum, lineOperationParsedValues.OperationTokenValue, lineOperationParsedValues.OperationLabel)
		if operandParserErr != nil {
			return operandParserErr
		}

	default:
		panic("Ruh roh")
	}
	return nil
}
