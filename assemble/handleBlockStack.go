package assemble

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/util"
	"strings"
)

var nestedBlockOps []string

// Things are in a block operation of some sort
func handleBlockStack(
	reformattedLine string,
	lineOperationParsedValues *util.LineOperationParsedValues,
) error {

	// End operations, if valid, will close the currently open block
	// However, blocks can also be nested so this doubles as a "close the nested block" too.
	if blockStack.NEW_IsEndOperation(lineOperationParsedValues) {
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
				return errorHandler.AddNew(enumErrorCodes.DirectiveUnmatchedEndBlock, lastStartOp)
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
		return nil

	}

	//----------------------------------------------
	//----------------------------------------------
	//All other non-ending ops drop down to here...
	//----------------------------------------------
	//----------------------------------------------

	//See if parent allows capturing this block
	err := blockStack.CheckOperationIsCapturable(reformattedLine, lineOperationParsedValues)
	if err != nil {
		return err // ❌ Fails
	}

	// -------------------------------------------
	// Check for inline evaluation
	// Will evalute the line right here instead of processing at the end
	// Mainly for KV statements in IKV blocks

	isInlineEval := blockStack.CheckInlineEval()

	if isInlineEval {
		err := parseAndProcessOperandString(
			reformattedLine,
			lineOperationParsedValues,
		)
		if err != nil {
			return err // ❌ Fails
		}
		return nil
	}

	// -------------------------------------------
	// If whatever's incoming is not the start of a new nested block operation, simply
	// append it for later processing and exit

	if !blockStack.NEW_IsStartOperation(lineOperationParsedValues) {
		_ = blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
		return nil
	}

	// -------------------------------------------
	// If this is reached, then whatever's incoming is
	// either the start of a NEW nested block operation
	// OR an alternate block starter for the current one
	// -------------------------------------------

	// See if there's NO corresponding ending operation. For directives like ELSEIF, these are
	// suboperations of the parent IF and thus they themselves have no such ops.
	// In such cases, the operation should be processed now PROVIDED it's not a nested one.

	if !blockStack.NEW_CheckEndOpExistsForStartOp(strings.ToUpper(lineOperationParsedValues.OperationTokenValue)) {

		//If this is nested, append and exit
		if len(nestedBlockOps) > 0 {
			_ = blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
			return nil
		}

		//If not nested, process now
		err := parseAndProcessOperandString(
			reformattedLine,
			lineOperationParsedValues,
		)
		if err != nil {
			return err // ❌ Fails
		}
		return nil
	}

	// -------------------------------------------
	// All exceptions have failed!
	// Simply append the operation to the end

	_ = blockStack.CheckOperationIsCapturableAndAppend(reformattedLine, lineOperationParsedValues)
	nestedBlockOps = append(nestedBlockOps, strings.ToUpper(lineOperationParsedValues.OperationTokenValue))

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

	//lines := blockStack.GetCapturedLinesWithPtr(currentOp)
	processedLines := blockStack.GetProcessedLinesWithPtr(currentOp)
	monitorStack := blockStack.GetBlockEntriesWithPtr(tempNewOp)
	postFn := blockStack.GetPostFnWithPtr(currentOp)
	if postFn != nil {
		defer postFn()
	}

	for _, pl := range *processedLines {
		lines := &pl.CapturedLines
		scope := pl.Scope
		err := processInner(lines, scope, monitorStack)
		if err != nil {
			return err
		}
	}

	return nil
}

func processInner(lines *[]blockStack.CapturedLine, scope blockStack.ProcessLineScope, monitorStack *[]blockStack.CaptureBlock) error {
	var processCapturedErr error

	symbolAsNodeTable.AddChildBlockScope(scope)
	defer symbolAsNodeTable.PopChildBlockScope()

	for j, b := range *lines {
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

		fileStack.AddSubOp(uint(j+1), b.OriginalLine)

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
