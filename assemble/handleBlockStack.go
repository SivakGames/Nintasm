package assemble

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/util"
)

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

	var flag bool
	if fromMainLoop {
		flag = blockStack.GetCurrentOperationEvaluatesCapturedNodesFlag()
	} else {
		flag = true
	}

	fmt.Println(flag)

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
		if !flag {
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

	//Push new entry one onto main stack in case a new stack will need to be evaluated
	blockStack.PushOntoMainStack()
	defer blockStack.PopFromMainStack()

	//Make a pointer to what was jus pushed onto the main stack
	newlyEvaluatedBlockStack := blockStack.GetCurrentStack()

	err := processCapturedLines(processingBlockStack, newlyEvaluatedBlockStack)
	if err != nil {
		return err
	}

	return nil
}

func processCapturedLines(
	currentLevel *[]blockStack.BlockOperationStack,
	upperLevel *[]blockStack.BlockOperationStack) error {
	var processCapturedErr error

	//Iterate over captured lines
	for _, b := range (*currentLevel)[0].CapturedLines {
		if len(*upperLevel) > 0 {
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
