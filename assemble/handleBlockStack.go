package assemble

import (
	"misc/nintasm/assemble/blockStack2"
	"misc/nintasm/interpreter"
	"misc/nintasm/util"
)

func handleBlockStack(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	isStartOrEndOperation := blockStack2.CheckIfNewStartEndOperation(lineOperationParsedValues)

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
		if blockStack2.CheckIfEndOperationAndGoesToProcessing(lineOperationParsedValues) {
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
		if blockStack2.GetCurrentInvokeOperationEvalFlag() {
			err := parseOperandStringAndProcess(
				reformattedLine,
				lineOperationParsedValues,
			)

			if err != nil {
				return err // ❌ Fails
			}
		} else {
			err := blockStack2.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func preProcessBlockStack() error {
	currentOp := blockStack2.GetCurrentOpPtr()
	blockStack2.AddNewInvokeOperationCollection() //Create new temp stack
	tempNewOp := blockStack2.GetCurrentOpPtr()
	err := readCapturedLines(currentOp, tempNewOp)
	if err != nil {
		return err
	}
	blockStack2.DestroyTempCollection(tempNewOp) //Remove upper level buffer stack
	blockStack2.ClearBlockEntriesWithPtr(currentOp)
	return nil
}

func readCapturedLines(
	currentOp *blockStack2.InvokeOperation,
	tempNewOp *blockStack2.InvokeOperation) error {
	var processCapturedErr error

	lines := blockStack2.GetLinesWithPtr(currentOp)
	monitorStack := blockStack2.GetBlockEntriesWithPtr(tempNewOp)

	//Iterate over captured lines
	for _, b := range *lines {
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

		// Ordinarily things will come here

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
