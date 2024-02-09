package assemble

import (
	enumParserTypes "misc/nintasm/enums/parserTypes"
	"misc/nintasm/parser"
)

func Start(lines []string) error {

	lineInitParzival := parser.NewInitialLineParser()
	lineOperationParzival := parser.NewOperationParser()

	directiveOperandParzival := parser.NewDirectiveOperandParser()
	instructionOperandParzival := parser.NewInstructionOperandParser()
	labelOperandParzival := parser.NewLabelOperandParser()

	instructionOperandParzival.ShouldParseInstructions = true

	var lineCounter uint = 0

	// Iterate over all lines
	for _, l := range lines {
		lineCounter++
		reformattedLine, lineInitErr := lineInitParzival.Process(l)
		if lineInitErr != nil {
			return lineInitErr
		}
		if len(reformattedLine) == 0 {
			continue
		}
		lineOperationErr := lineOperationParzival.Process(reformattedLine)
		if lineOperationErr != nil {
			return lineOperationErr
		}

		operationType, operationSimpleType, operationValue, operationLabel, operandStartPosition := lineOperationParzival.GetOperationDetails()

		switch operationSimpleType {
		case enumParserTypes.Directive:
			operandParserErr := directiveOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				return operandParserErr
			}
			operandParserErr = directiveOperandParzival.Process(operationType, operationValue)
			if operandParserErr != nil {
				return operandParserErr
			}

		case enumParserTypes.Instruction:
			operandParserErr := instructionOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				return operandParserErr
			}
			operandParserErr = instructionOperandParzival.Process(operationValue)
			if operandParserErr != nil {
				return operandParserErr
			}

		//	case enumParserTypes.Macro:
		//		fmt.Println("Mack")
		//		fmt.Println(optype, opval, opPos)
		case enumParserTypes.Label:
			operandParserErr := labelOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				return operandParserErr
			}
			operandParserErr = labelOperandParzival.Process(operationType, operationValue, operationLabel)
			if operandParserErr != nil {
				return operandParserErr
			}
		}

	}
	return nil
}
