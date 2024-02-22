package parser

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/blockStack"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/environment/macroTable"
	"regexp"
)

type replacer struct {
	replaceRegex  *regexp.Regexp
	replaceString string
}

func newReplacer(index string, replaceString string) replacer {
	return replacer{
		replaceRegex:  regexp.MustCompile(`\\` + index + `\b`),
		replaceString: replaceString,
	}
}

type MacroOperandParser struct {
	OperandParser
	capturedLinesToUnpack []blockStack.CapturedLine
	replacementsList      []replacer
}

func NewMacroOperandParser() MacroOperandParser {
	return MacroOperandParser{}
}

// Main macro invoke parser
func (mop *MacroOperandParser) Process(macroName string) error {
	var err error

	mop.replacementsList = mop.replacementsList[:0]

	mop.capturedLinesToUnpack, err = macroTable.LookupAndGetMacroInEnvironment(macroName, macroTable.Macro)
	if err != nil {
		return err
	}

	err = mop.getMacroReplacementsList()
	if err != nil {
		return err // ❌ Fails
	}

	return nil
}

func (mop *MacroOperandParser) getMacroReplacementsList() error {
	var err error

	//No operands at all
	if mop.lookaheadType == enumTokenTypes.None {
		return nil // 🟢 Succeeds
	}

	//No commas at the beginning...
	if mop.lookaheadType == enumTokenTypes.DELIMITER_comma {
		return errors.New("Operand list \x1b[38;5;202mCANNOT\x1b[0m start with a comma!") // ❌ Fails
	}

	firstReplacement, err := mop.getMacroReplacement()
	mop.replacementsList = append(mop.replacementsList, newReplacer("1", firstReplacement))
	i := 2

	for mop.lookaheadType != enumTokenTypes.None && mop.lookaheadType == enumTokenTypes.DELIMITER_comma {
		err = mop.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_comma)
		if err != nil {
			return err // ❌ Fails
		}
		subsequentReplacement, err := mop.getMacroReplacement()
		if err != nil {
			return err // ❌ Fails
		}
		mop.replacementsList = append(mop.replacementsList, newReplacer(fmt.Sprintf("%d", i), subsequentReplacement))
		i++
	}

	return nil
}

func (mop *MacroOperandParser) getMacroReplacement() (string, error) {
	var closingTokenEnum []enumTokenTypes.Def

	replacement := ""
	closingTokenEnum = append(closingTokenEnum, enumTokenTypes.DELIMITER_comma)

	for len(closingTokenEnum) > 0 && mop.lookaheadType != enumTokenTypes.None {
		topOfStackEnum := closingTokenEnum[len(closingTokenEnum)-1]
		switch mop.lookaheadType {
		case topOfStackEnum:
			closingTokenEnum = closingTokenEnum[:len(closingTokenEnum)-1]
			if len(closingTokenEnum) > 0 {
				err := mop.eatFreelyAndAdvance(topOfStackEnum)
				if err != nil {
					return replacement, err
				}
			}

		case enumTokenTypes.DELIMITER_leftCurlyBrace:
			if topOfStackEnum == enumTokenTypes.DELIMITER_rightCurlyBrace {
				return replacement, errors.New("Macro args - Must close curly brace before opening another!")
			}
			err := mop.eatFreelyAndAdvance(enumTokenTypes.DELIMITER_leftCurlyBrace)
			if err != nil {
				return replacement, err
			}
			closingTokenEnum = append(closingTokenEnum, enumTokenTypes.DELIMITER_rightCurlyBrace)

		default:
			replacement += mop.lookaheadValue
			err := mop.eatFreelyAndAdvance(mop.lookaheadType)
			if err != nil {
				return replacement, err
			}
		}
	}
	if len(closingTokenEnum) > 1 {
		return replacement, errors.New("Unclosed???")
	}

	return replacement, nil
}

func (mop *MacroOperandParser) GetUnpackLinesRef() *[]blockStack.CapturedLine {
	return &mop.capturedLinesToUnpack
}

func (mop *MacroOperandParser) ApplyReplacementsToCapturedLine(capturedLineIndex int) blockStack.CapturedLine {
	replacedCapturedLine := mop.capturedLinesToUnpack[capturedLineIndex]
	for _, replacementListItem := range mop.replacementsList {
		replacedCapturedLine.OriginalLine = replacementListItem.replaceRegex.ReplaceAllString(replacedCapturedLine.OriginalLine, replacementListItem.replaceString)
	}

	return replacedCapturedLine
}
