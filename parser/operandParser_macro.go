package parser

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
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

// Main macro invoke parser
func (mop *MacroOperandParser) Process(macroName string) error {
	var err error
	var capturedLinesToProcess []blockStack.CapturedLine

	macroTable.AppendToReplacementStack()
	mop.capturedLinesToUnpack, err = macroTable.LookupAndGetMacroInEnvironment(macroName, macroTable.Macro)
	if err != nil {
		return err
	}

	operandList, err := mop.GetOperandList(
		MACRO_MIN_OPERANDS, MACRO_MAX_OPERANDS, MACRO_MANAULLY_EVALS,
		[]string{"macro"},
	)
	if err != nil {
		return err // ❌ Fails
	}

	// Put numeric replacements on stack
	for i, operand := range operandList {
		macroTable.AddToReplacementListOnTopOfStack(fmt.Sprintf("\\%d", i+1), operand.NodeValue)
	}
	macroTable.AddNumToReplacementListOnTopOfStack(fmt.Sprintf("%d", len(operandList)))

	//Iterate over macro lines and apply replacements
	for i := range mop.capturedLinesToUnpack {
		capturedLinesToProcess = append(capturedLinesToProcess, mop.ApplyReplacementsToCapturedLine(i))
	}

	blockStack.PushCaptureBlock("IM", operandList)

	currentStackOp := blockStack.GetCurrentCaptureBlock()
	currentStackOp.CapturedLines = capturedLinesToProcess

	return nil
}

func (mop *MacroOperandParser) ApplyReplacementsToCapturedLine(capturedLineIndex int) blockStack.CapturedLine {
	replacedCapturedLine := mop.capturedLinesToUnpack[capturedLineIndex]
	replacementList := macroTable.GetReplacementListOnTopOfStack()

	for _, replacementListItem := range *replacementList {
		replacedCapturedLine.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllLiteralString(replacedCapturedLine.OriginalLine, replacementListItem.ReplaceString)
	}

	return replacedCapturedLine
}

func (mop *MacroOperandParser) EndInvokeMacro() {
	macroTable.PopFromReplacementStack()
}
