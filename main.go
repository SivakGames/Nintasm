package main

import (
	"fmt"
	"misc/nintasm/parser"
	"misc/nintasm/parser/parserTypes"
	"time"
)

func main() {
	process()
}

func process() {
	start := time.Now()
	lines := []string{" lsr a bank(5,7,234)"}
	//lines := make([]string, 84000)
	//for i := range lines {
	//	lines[i] = " lda 5"
	//}

	lineInitParzival := parser.NewInitialLineParser()
	lineOperationParzival := parser.NewOperationParser()
	instructionOperandParzival := parser.NewInstructionOperandParser()
	labelOperandParzival := parser.NewLabelOperandParser()

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
			instructionOperandParzival.Process(operationValue)
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
