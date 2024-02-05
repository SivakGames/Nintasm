package assemble

import (
	"misc/nintasm/parser"
	"misc/nintasm/parser/parserTypes"
	"misc/nintasm/romBuilder"
)

func Start(lines []string) error {

	lineInitParzival := parser.NewInitialLineParser()
	lineOperationParzival := parser.NewOperationParser()

	directiveOperandParzival := parser.NewDirectiveOperandParser()
	instructionOperandParzival := parser.NewInstructionOperandParser()
	labelOperandParzival := parser.NewLabelOperandParser()

	instructionOperandParzival.ShouldParseInstructions = true

	err := romBuilder.AddNewRomSegment(0x20000, 0x4000)
	err = romBuilder.AddNewRomSegment(0x20000, 0x4000)
	if err != nil {
		return err
	}

	var lineCounter uint = 0

	// Iterate over all lines
	for _, l := range lines {
		lineCounter++
		refoLine, lineInitErr := lineInitParzival.Process(l)
		if lineInitErr != nil {
			return lineInitErr
		}
		lineOperationErr := lineOperationParzival.Process(refoLine)
		if lineOperationErr != nil {
			return lineOperationErr
		}

		operationType, operationSimpleType, operationValue, operationLabel, operandStartPosition := lineOperationParzival.GetOperationDetails()

		switch operationSimpleType {
		case parserTypes.Directive:
			operandParserErr := directiveOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				return operandParserErr
			}
			operandParserErr = directiveOperandParzival.Process(operationType, operationValue)
			if operandParserErr != nil {
				return operandParserErr
			}

		case parserTypes.Instruction:
			operandParserErr := instructionOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				return operandParserErr
			}
			operandParserErr = instructionOperandParzival.Process(operationValue)
			if operandParserErr != nil {
				return operandParserErr
			}

		//	case parserTypes.Macro:
		//		fmt.Println("Mack")
		//		fmt.Println(optype, opval, opPos)
		case parserTypes.Label:
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
