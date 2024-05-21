package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
)

func evalRepeat(directiveName string, operandList *[]Node) error {
	// Check and validate repeat amount
	numRepeatsNode, err := interpreter.EvaluateNode((*operandList)[0])
	if err != nil {
		return err
	}
	if !operandFactory.ValidateNodeIsNumeric(&numRepeatsNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric) // ❌ Fails
	} else if !operandFactory.ValidateNumericNodeIsGTEValue(&numRepeatsNode, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTE, 1) // ❌ Fails
	}

	evaluatedNodes := []Node{numRepeatsNode}

	// Check if iterator node was set
	if len(*operandList) > 1 {
		iterNameNode := (*operandList)[1]
		if !operandFactory.ValidateNodeIsIdentifier(&iterNameNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[1])
	}

	blockStack.PushCaptureBlock(directiveName, evaluatedNodes)

	return nil
}

func evalEndRepeat() error {
	capturedLines, operandList := blockStack.GetCurrentCaptureBlockCapturedLinesAndOperandList()

	//Extract repeatAmount
	repeatAmountNumber := (*operandList)[0].AsNumber
	repeatAmount := int(repeatAmountNumber)

	//Extract iterator name (if any)
	iteratorName := ""
	if len(*operandList) > 1 {
		iteratorName = (*operandList)[1].NodeValue
	}

	processedLines := []blockStack.ProcessLine{}

	for i := 0; i < repeatAmount; i++ {
		repeatScope := blockStack.ProcessLineScope{}
		if iteratorName != "" {
			repeatScope[iteratorName] = operandFactory.CreateNumericLiteralNode(float64(i))
		}

		pl := blockStack.GenerateProcessedLine(repeatScope, *capturedLines)
		processedLines = append(processedLines, pl)
	}

	/*
		replacedLines := make([]blockStack.CapturedLine, len(*capturedLines)*repeatAmount)
		replacedIndex := 0

		if iteratorName == "" {
			for i := 0; i < repeatAmount; i++ {
				for _, j := range *capturedLines {
					replacedLines[replacedIndex] = j
					replacedIndex++
				}
			}
		} else {
			iterNameAsRegex := regexp.MustCompile(`\` + iteratorName + `\b`)
			for i := 0; i < repeatAmount; i++ {
				replaceNum := strconv.Itoa(i)
				for _, j := range *capturedLines {
					replacedOriginalLine := iterNameAsRegex.ReplaceAllString(j.OriginalLine, replaceNum)
					replacedOperationLabel := iterNameAsRegex.ReplaceAllString(j.OperationLabel, replaceNum)
					j.OriginalLine = replacedOriginalLine
					j.OperationLabel = replacedOperationLabel
					replacedLines[replacedIndex] = j
					replacedIndex++
				}
			}
		}
		if len(replacedLines) == 0 {
			errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
		}
	*/
	blockStack.NEW_PopCaptureBlockPrepProcessBlock(processedLines)

	//blockStack.PopCaptureBlockThenExtendCapturedLines(replacedLines)
	return nil
}
