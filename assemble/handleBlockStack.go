package assemble

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/fileStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/util"
	"strings"
)

var nestedBlockOps []string

// Things are in a block operation of some sort
func handleBlockStack(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {
	// End ops close the block
	// However, blocks can be nested so this also doubles as a "close the nested block"
	isEndOp := blockStack.NEW_IsEndOperation(lineOperationParsedValues)

	//End ops will close
	if isEndOp {
		if len(nestedBlockOps) > 0 {
			lastStartOp := nestedBlockOps[len(nestedBlockOps)-1]
			matches := blockStack.NEW_CheckEndOpVsStartOp(lineOperationParsedValues.OperationTokenValue, lastStartOp)
			if matches {
				nestedBlockOps = nestedBlockOps[:len(nestedBlockOps)-1]
				err := blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Unmatched ")
			}
			return nil
		}

		//Process end op
		err := parseAndProcessOperandString(
			reformattedLine,
			lineOperationParsedValues,
		)
		if err != nil {
			return err // ❌ Fails
		}

		//Process captured lines
		err = preProcessBlockStack()
		if err != nil {
			return err // ❌ Fails
		}
		//Mainly set by namespaces - will clear the overriding parent label
		if interpreter.PopParentLabelWhenBlockOpDone {
			interpreter.PopParentLabel()
			interpreter.PopParentLabelWhenBlockOpDone = false
		}
	} else {
		//See if parent allows capturing this block
		err := blockStack.CheckOperationIsCapturable(reformattedLine, lineOperationParsedValues)
		if err != nil {
			return err // ❌ Fails
		}

		isNewStartOp := blockStack.NEW_IsStartOperation(lineOperationParsedValues)

		if !isNewStartOp || len(nestedBlockOps) > 0 {
			_ = blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
		} else {
			endOpExists := blockStack.NEW_CheckEndOpExistsForStartOp(strings.ToUpper(lineOperationParsedValues.OperationTokenValue))
			if endOpExists {
				_ = blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
				nestedBlockOps = append(nestedBlockOps, strings.ToUpper(lineOperationParsedValues.OperationTokenValue))
			} else {
				err := parseAndProcessOperandString(
					reformattedLine,
					lineOperationParsedValues,
				)
				if err != nil {
					return err // ❌ Fails
				}
			}
		}
	}

	//isStartOrEndOperation := blockStack.CheckIfNewStartEndOperation(lineOperationParsedValues)

	//	if isStartOrEndOperation {

	// Operation will always process
	/*err := parseAndProcessOperandString(
		reformattedLine,
		lineOperationParsedValues,
	)
	if err != nil {
		return err // ❌ Fails
	}*/

	//If ending, iterate bottom of stack and parse all captured operations (if any)
	/*
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
		} */

	//	} else {
	//If in forced eval mode, evaluate the node right here
	/* if blockStack.GetCaptureBlockListEvalFlag() {
		err := blockStack.CheckOperationIsCapturable(reformattedLine, lineOperationParsedValues)
		if err != nil {
			return err // ❌ Fails
		}
		err = parseAndProcessOperandString(
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
	} */
	//	}

	return nil
}

func preProcessBlockStack() error {
	currentOp := blockStack.GetCurrentOpPtr()
	blockStack.AddNewCaptureBlockListNode() //Create new temp stack

	currentBlockEntries := blockStack.GetBlockEntriesWithPtr(currentOp)
	if len(*currentBlockEntries) == 0 {
		return nil
	}

	tempNewOp := blockStack.GetCurrentOpPtr()
	currentOpName := blockStack.GetCapturedLinesOpNameWithPtr(currentOp)
	err := readCapturedLines(currentOp, tempNewOp)
	if err != nil {
		return err
	}
	blockStack.DestroyCaptureBlockListNodeWithPointer(tempNewOp) //Remove upper level buffer stack
	blockStack.ClearBlockEntriesWithPtr(currentOp)
	fileStack.ClearSubOp()
	if blockStack.GetExitOpName() != "" && currentOpName == blockStack.GetExitOpName() {
		blockStack.SetExitOpName("")
	}

	return nil
}

func readCapturedLines(
	currentOp *blockStack.CaptureBlockListNode,
	tempNewOp *blockStack.CaptureBlockListNode) error {
	var processCapturedErr error

	lines := blockStack.GetCapturedLinesWithPtr(currentOp)
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
			if blockStack.GetExitOpName() != "" {
				break
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
		processCapturedErr = parseAndProcessOperandString(
			b.OriginalLine,
			&processOperandArguments,
		)
		if processCapturedErr != nil {
			return processCapturedErr // ❌ Fails
		}
		if blockStack.GetExitOpName() != "" {
			break
		}
	}

	return nil
}
