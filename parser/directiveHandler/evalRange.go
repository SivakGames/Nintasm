package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
)

func evalRange(directiveName string, operandList *[]Node) error {
	// Check and validate repeat amount
	arrayNode, err := interpreter.EvaluateNode((*operandList)[0])
	if err != nil {
		return err
	}
	if !operandFactory.ValidateNodeIsArray(&arrayNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotArray) // ❌ Fails
	}

	indexNameNode := (*operandList)[1]
	if !operandFactory.ValidateNodeIsIdentifier(&indexNameNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
	}

	evaluatedNodes := []Node{arrayNode, indexNameNode}

	// Check if iterator node was set
	if len(*operandList) > 2 {
		valueNameNode := (*operandList)[2]
		if !operandFactory.ValidateNodeIsIdentifier(&valueNameNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotSubstitutionID) // ❌ Fails
		}
		evaluatedNodes = append(evaluatedNodes, (*operandList)[2])
	}

	blockStack.PushCaptureBlock(directiveName, evaluatedNodes)
	return nil
}

func evalEndRange() error {
	capturedLines, operandList := blockStack.GetCurrentCaptureBlockCapturedLinesAndOperandList()

	//Extract the base array
	rangeArray := (*operandList)[0].ArgumentList

	//Extract index name
	indexName := (*operandList)[1].NodeValue

	//Extract value name
	valueName := ""
	if len(*operandList) > 2 {
		valueName = (*operandList)[2].NodeValue
	}

	processedLines := []blockStack.ProcessLine{}

	for i, v := range *rangeArray {
		repeatScope := blockStack.ProcessLineScope{}
		repeatScope[indexName] = operandFactory.CreateNumericLiteralNode(float64(i))
		if valueName != "" {
			repeatScope[valueName] = v
		}

		pl := blockStack.GenerateProcessedLine(repeatScope, *capturedLines)
		processedLines = append(processedLines, pl)
	}

	/*
		replacedLines := make([]blockStack.CapturedLine, len(*capturedLines)*len(*rangeArray))
		replacedIndex := 0

		indexNameAsRegex := regexp.MustCompile(`\` + indexName + `\b`)
		var valueNameAsRegex *regexp.Regexp
		if len(valueName) > 0 {
			valueNameAsRegex = regexp.MustCompile(`\` + valueName + `\b`)
		}

		for i, v := range *rangeArray {
			replaceIndex := strconv.Itoa(i)
			replaceValue := ""
			if valueNameAsRegex != nil {
				replaceValue = getReplaceValueAsString(&v)
			}
			for _, j := range *capturedLines {
				replacedOriginalLine := indexNameAsRegex.ReplaceAllString(j.OriginalLine, replaceIndex)
				replacedOperationLabel := indexNameAsRegex.ReplaceAllString(j.OperationLabel, replaceIndex)
				if valueNameAsRegex != nil {
					replacedOriginalLine = valueNameAsRegex.ReplaceAllString(replacedOriginalLine, replaceValue)
					replacedOperationLabel = valueNameAsRegex.ReplaceAllString(replacedOperationLabel, replaceValue)
				}
				j.OriginalLine = replacedOriginalLine
				j.OperationLabel = replacedOperationLabel
				replacedLines[replacedIndex] = j
				replacedIndex++
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

func getReplaceValueAsString(targetNode *Node) string {
	replaceString := ""

	switch targetNode.NodeType {
	case enumNodeTypes.NumericLiteral, enumNodeTypes.BooleanLiteral:
		replaceString = targetNode.NodeValue
	case enumNodeTypes.StringLiteral:
		replaceString = fmt.Sprintf("\"%s\"", targetNode.NodeValue)
	case enumNodeTypes.MultiByte:
		accumulatedString := ""
		for i, v := range *targetNode.ArgumentList {
			returnedString := getReplaceValueAsString(&v)
			accumulatedString += returnedString
			if i < len(*targetNode.ArgumentList)-1 {
				accumulatedString += ","
			}
		}
		replaceString = fmt.Sprintf("[%s]", accumulatedString)
	}

	return replaceString
}
