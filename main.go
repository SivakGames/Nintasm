package main

import (
	"fmt"
	"misc/nintasm/parser"
	"misc/nintasm/parser/parserTypes"
	"misc/nintasm/romBuilder"
	"time"
)

func main() {
	process()
}

func process() {
	start := time.Now()
	lines := []string{"JAJA = 11", "JAJ2A equ 42", "RESET:", " sei ", " cld", " lda JAJA"}
	//lines := make([]string, 0x2000)
	//for i := range lines {
	//	lines[i] = " lda [555], y "
	//}

	lineInitParzival := parser.NewInitialLineParser()
	lineOperationParzival := parser.NewOperationParser()
	instructionOperandParzival := parser.NewInstructionOperandParser()
	labelOperandParzival := parser.NewLabelOperandParser()

	instructionOperandParzival.ShouldParseInstructions = true

	err := romBuilder.AddNewRomSegment(0x20000, 0x4000)
	err = romBuilder.AddNewRomSegment(0x20000, 0x4000)
	if err != nil {
		fmt.Println(err)
		return
	}

	var lineCounter uint = 0

	// Iterate over all lines
	for _, l := range lines {
		lineCounter++
		refoLine, lineInitErr := lineInitParzival.Process(l)
		if lineInitErr != nil {
			fmt.Println(lineInitErr)
			return
		}
		lineOperationErr := lineOperationParzival.Process(refoLine)
		if lineOperationErr != nil {
			fmt.Println(lineOperationErr)
			return
		}

		operationType, operationSimpleType, operationValue, operationLabel, operandStartPosition := lineOperationParzival.GetOperationDetails()

		switch operationSimpleType {
		//	case parserTypes.Directive:
		//		fmt.Println("Directive")
		//		fmt.Println(optype, opval, opPos)
		case parserTypes.Instruction:
			operandParserErr := instructionOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				fmt.Println(operandParserErr)
				return
			}
			operandParserErr = instructionOperandParzival.Process(operationValue)
			if operandParserErr != nil {
				fmt.Println(operandParserErr)
				return
			}

		//	case parserTypes.Macro:
		//		fmt.Println("Mack")
		//		fmt.Println(optype, opval, opPos)
		case parserTypes.Label:
			operandParserErr := labelOperandParzival.SetupOperandParser(l, operandStartPosition)
			if operandParserErr != nil {
				fmt.Println(operandParserErr)
				return
			}
			labelOperandParzival.Process(operationType, operationValue, operationLabel)
		}

	}

	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)

}
