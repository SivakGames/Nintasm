package assemble

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/fileStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/util"
)

func handleBlockStack(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	isStartOrEndOperation := blockStack.CheckIfNewStartEndOperation(lineOperationParsedValues)

	if isStartOrEndOperation {
		// Operation will always process
		err := parseOperandStringAndProcess(
			reformattedLine,
			lineOperationParsedValues,
		)
		if err != nil {
			return err // ❌ Fails
		}

		//If ending, iterate bottom of stack and parse all captured operations (if any)
		if blockStack.CheckIfEndOperationAndGoesToProcessing(lineOperationParsedValues) {
			err := preProcessBlockStack()
			if err != nil {
				return err // ❌ Fails
			}
			//Mainly set by namespaces - will clear the overriding parent label
			if interpreter.PopParentLabelWhenBlockOpDone {
				interpreter.PopParentLabel()
				interpreter.PopParentLabelWhenBlockOpDone = false
			}
		}

	} else {
		//If in forced eval mode, evaluate the node right here
		if blockStack.GetCaptureBlockListEvalFlag() {
			err := blockStack.CheckOperationIsCapturable(reformattedLine, lineOperationParsedValues)
			if err != nil {
				return err // ❌ Fails
			}
			err = parseOperandStringAndProcess(
				reformattedLine,
				lineOperationParsedValues,
			)
			if err != nil {
				return err // ❌ Fails
			}

		} else {
			err := blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func preProcessBlockStack() error {
	currentOp := blockStack.GetCurrentOpPtr()
	blockStack.AddNewCaptureBlockListNode() //Create new temp stack
	tempNewOp := blockStack.GetCurrentOpPtr()
	err := readCapturedLines(currentOp, tempNewOp)
	if err != nil {
		return err
	}
	blockStack.DestroyCaptureBlockListNodeWithPointer(tempNewOp) //Remove upper level buffer stack
	blockStack.ClearBlockEntriesWithPtr(currentOp)
	return nil
}

func readCapturedLines(
	currentOp *blockStack.CaptureBlockListNode,
	tempNewOp *blockStack.CaptureBlockListNode) error {
	var processCapturedErr error

	lines := blockStack.GetLinesWithPtr(currentOp)
	monitorStack := blockStack.GetBlockEntriesWithPtr(tempNewOp)

	//Iterate over captured lines
	for i, b := range *lines {
		if len(*monitorStack) > 0 {
			newlyEvaluatedParsedValues := util.LineOperationParsedValues{
				OperandStartPosition: b.OperandStartPosition,
				OperationLabel:       b.OperationLabel,
				OperationTokenEnum:   b.OperationTokenEnum,
				OperationTokenValue:  b.OperationTokenValue,
				ParentParserEnum:     b.ParentParserEnum,
			}
			processCapturedErr = handleBlockStack(b.OriginalLine, &newlyEvaluatedParsedValues)
			if processCapturedErr != nil {
				return processCapturedErr // ❌ Fails
			}
			continue
		}

		fileStack.AddSubOp(uint(i+1), b.OriginalLine)

		processOperandArguments := util.NewLineOperationParsedValues(b.OperandStartPosition,
			b.OperationLabel,
			b.OperationTokenEnum,
			b.OperationTokenValue,
			b.ParentParserEnum,
		)
		processCapturedErr = parseOperandStringAndProcess(
			b.OriginalLine,
			&processOperandArguments,
		)
		if processCapturedErr != nil {
			return processCapturedErr // ❌ Fails
		}
	}

	return nil
}
