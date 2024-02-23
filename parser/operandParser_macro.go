package parser

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/interpreter/environment/macroTable"
)

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

	macroTable.AppendToReplacementStack()

	mop.capturedLinesToUnpack, err = macroTable.LookupAndGetMacroInEnvironment(macroName, macroTable.Macro)
	if err != nil {
		return err
	}

	operandList, err := mop.GetOperandList(0, 16, true, []string{"macro"})
	if err != nil {
		return err // ❌ Fails
	}

	for i, operand := range operandList {
		macroTable.AddToReplacementListOnTopOfStack(fmt.Sprintf("\\%d", i+1), operand.NodeValue)
	}

	return nil
}

func (mop *MacroOperandParser) GetUnpackLinesRef() *[]blockStack.CapturedLine {
	return &mop.capturedLinesToUnpack
}

func (mop *MacroOperandParser) ApplyReplacementsToCapturedLine(capturedLineIndex int) blockStack.CapturedLine {
	replacedCapturedLine := mop.capturedLinesToUnpack[capturedLineIndex]
	replacementList := macroTable.GetReplacementListOnTopOfStack()

	for _, replacementListItem := range *replacementList {
		replacedCapturedLine.OriginalLine = replacementListItem.ReplaceRegex.ReplaceAllString(replacedCapturedLine.OriginalLine, replacementListItem.ReplaceString)
	}

	return replacedCapturedLine
}

func (mop *MacroOperandParser) PopFromStack() {
	macroTable.PopFromReplacementStack()
}
