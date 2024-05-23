package parser

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/environment/macroTable"
)

const MACRO_MIN_OPERANDS = 0
const MACRO_MAX_OPERANDS = 16
const MACRO_MANAULLY_EVALS = true

type MacroOperandParser struct {
	OperandParser
	capturedLinesToUnpack []blockStack.CapturedLine
}

func NewMacroOperandParser() MacroOperandParser {
	return MacroOperandParser{}
}

var macroCaptureMask []string

func init() {
	captureMasks := make([]string, MACRO_MAX_OPERANDS)
	for i := range captureMasks {
		captureMasks[i] = "macro"
	}
	macroCaptureMask = captureMasks
}

// Main macro invoke parser
func (mop *MacroOperandParser) Process(macroName string) error {
	var err error
	var capturedLinesToProcess []blockStack.CapturedLine
	var validArguments *[]string

	macroTable.AppendToReplacementStack()
	mop.capturedLinesToUnpack, validArguments, err = macroTable.LookupAndGetMacroInEnvironment(macroName)
	if err != nil {
		return err
	}

	operandList, err := mop.GetOperandList(
		MACRO_MIN_OPERANDS, MACRO_MAX_OPERANDS, MACRO_MANAULLY_EVALS,
		macroCaptureMask,
	)
	if err != nil {
		return err // ❌ Fails
	}

	if validArguments != nil {
		suppliedArguments := len(operandList)
		neededArguments := len(*validArguments)

		if suppliedArguments < neededArguments {
			return errorHandler.AddNew(enumErrorCodes.MacroInvokeTooFewArgs, suppliedArguments, neededArguments)
		}
		if suppliedArguments > neededArguments {
			return errorHandler.AddNew(enumErrorCodes.MacroInvokeTooManyArgs, suppliedArguments, neededArguments)
		}
		for i, operand := range operandList {
			macroTable.AddToReplacementListOnTopOfStack((*validArguments)[i], operand.NodeValue)
		}
	} else {
		// Put numeric replacements on stack
		for i, operand := range operandList {
			macroTable.AddToReplacementListOnTopOfStack(fmt.Sprintf("\\%d", i+1), operand.NodeValue)
		}
	}
	macroTable.AddNumArgsToReplacementListOnTopOfStack(fmt.Sprintf("%d", len(operandList)))

	//Iterate over macro lines and apply replacements
	for i := range mop.capturedLinesToUnpack {
		capturedLinesToProcess = append(capturedLinesToProcess, mop.ApplyReplacementsToCapturedLine(i))
	}

	blockStack.PushCaptureBlock("IM", operandList)

	blockStack.CopyPresetCapturedLinesToProcessedWithEmptyScope(&capturedLinesToProcess)

	//currentStackOp := blockStack.GetCurrentCaptureBlock()
	//currentStackOp.CapturedLines = capturedLinesToProcess

	return nil
}

func (mop *MacroOperandParser) ApplyReplacementsToCapturedLine(capturedLineIndex int) blockStack.CapturedLine {
	replacedCapturedLine := mop.capturedLinesToUnpack[capturedLineIndex]
	replacementList := macroTable.GetReplacementListOnTopOfStack()

	for _, replacementListItem := range *replacementList {
		replacedCapturedLine.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllLiteralString(replacedCapturedLine.OriginalLine, replacementListItem.ReplaceString)
		replacedCapturedLine.OperationLabel = replacementListItem.ReplaceRegex.ReplaceAllLiteralString(replacedCapturedLine.OperationLabel, replacementListItem.ReplaceString)
	}

	return replacedCapturedLine
}

func (mop *MacroOperandParser) EndInvokeMacro() {
	macroTable.PopFromReplacementStack()
}
