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
			//blockStack.ClearCurrentOperationEvaluatesCapturedNodesFlag()
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

/*

//==========================================================

// Main handler for when in a block operation that will capture multiple lines
// and process them.
func handleBlockStack(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
	fromMainLoop bool,
) error {
	var blockStackErr error

	//See if the incoming operation is for starting/ending a block
	isStartOrEndOperation := blockStack.CheckIfNewStartEndOperation(lineOperationParsedValues)

	if isStartOrEndOperation {
		blockStackErr = parseOperandStringAndProcess(
			reformattedLine,
			lineOperationParsedValues,
		)
		if blockStackErr != nil {
			return blockStackErr // ❌ Fails
		}

		//If ending, iterate bottom of stack and parse all captured operations (if any)
		if blockStack.CheckIfEndOperationAndClearStack(lineOperationParsedValues) {
			blockStackErr = reassignPointers()

			if blockStackErr != nil {
				return blockStackErr // ❌ Fails
			}
			blockStack.ClearCurrentOperationEvaluatesCapturedNodesFlag()
			//Mainly set by namespaces - will clear the overriding parent label
			if interpreter.PopParentLabelWhenBlockOpDone {
				interpreter.PopParentLabel()
				interpreter.PopParentLabelWhenBlockOpDone = false
			}
			if fromMainLoop {
				blockStack.PopFromStack()
			}
		}

	} else {
		//Either append the operation to the stack's captured lines or evaluate them now
		if !blockStack.GetCurrentOperationEvaluatesCapturedNodesFlag() {
			err := blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
			if err != nil {
				return err
			}
		} else {
			blockStackErr = parseOperandStringAndProcess(
				reformattedLine,
				lineOperationParsedValues,
			)

			if blockStackErr != nil {
				return blockStackErr // ❌ Fails
			}
		}
	}

	return nil
}

func reassignPointers() error {

	//Make a pointer to the current stack
	processingBlockStack := blockStack.GetCurrentStack()
	lines := (*processingBlockStack)[0].CapturedLines

	fmt.Println("\x1b[35m", lines, "\x1b[0m")

	blockStack.PushOntoMainStack() //Create new temp stack

	//Make a pointer to what was jus pushed onto the main stack
	newlyEvaluatedBlockStack := blockStack.GetCurrentStack()

	err := processCapturedLines(&lines, newlyEvaluatedBlockStack)
	if err != nil {
		return err
	}

	blockStack.PopFromMainStack() //Remove upper level buffer stack

	return nil
}

func processCapturedLines(
	lines *[]blockStack.CapturedLine,
	upperLevel *[]blockStack.BlockOperationStack) error {
	var processCapturedErr error

	fmt.Println("\x1b[34m", *lines, "\x1b[0m")
	fmt.Println("\x1b[33m", *upperLevel, "\x1b[0m")

	//Iterate over captured lines
	for _, b := range *lines {
		if len(*upperLevel) > 0 {
			fmt.Println("\x1b[33m", *upperLevel, "\x1b[0m")
			newlyEvaluatedParsedValues := util.LineOperationParsedValues{
				OperandStartPosition: b.OperandStartPosition,
				OperationLabel:       b.OperationLabel,
				OperationTokenEnum:   b.OperationTokenEnum,
				OperationTokenValue:  b.OperationTokenValue,
				ParentParserEnum:     b.ParentParserEnum,
			}
			processCapturedErr = handleBlockStack(b.OriginalLine, &newlyEvaluatedParsedValues, false)
			if processCapturedErr != nil {
				return processCapturedErr // ❌ Fails
			}
			blockStack.PopFromStack()
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
*/
