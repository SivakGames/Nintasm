package assemble

import (
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	"misc/nintasm/parser"
	"misc/nintasm/util"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var directiveOperandParser = parser.NewDirectiveOperandParser()
var instructionOperandParser = parser.NewInstructionOperandParser()
var labelOperandParser = parser.NewLabelOperandParser()
var macroOperandParser = parser.NewMacroOperandParser()

// ============================================================

// Main operand parsing...
func parseAndProcessOperandString(
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

		preProcessBlockStack()
		macroOperandParser.EndInvokeMacro()

	default:
		panic("🛑 Parent parsing operation could not be determined!")
	}

	return nil
}
