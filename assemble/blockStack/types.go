package blockStack

import (
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/util"
)

type Node = operandFactory.Node

// ++++++++++++++++++++++++++++++++++++

type CapturedLine struct {
	OriginalLine string
	util.LineOperationParsedValues
}

func newCapturedLine(originalLine string,
	operationLabel string,
	operationTokenEnum enumTokenTypes.Def,
	operationTokenValue string,
	operandStartPosition int,
	parentParserEnum enumParserTypes.Def,
) CapturedLine {
	return CapturedLine{
		OriginalLine: originalLine,
		LineOperationParsedValues: util.LineOperationParsedValues{
			OperationLabel:       operationLabel,
			OperationTokenEnum:   operationTokenEnum,
			OperationTokenValue:  operationTokenValue,
			OperandStartPosition: operandStartPosition,
			ParentParserEnum:     parentParserEnum,
		},
	}
}

// ++++++++++++++++++++++++++++++++++++

type blockEntry struct {
	BlockOperationName  string
	CapturedLines       []CapturedLine
	OperandList         []Node
	AlternateStackBlock *blockEntry
}

func newBlockEntry(blockOperationName string, operandList []Node) blockEntry {
	return blockEntry{
		BlockOperationName:  blockOperationName,
		CapturedLines:       []CapturedLine{},
		OperandList:         operandList,
		AlternateStackBlock: nil,
	}
}

// ++++++++++++++++++++++++++++++++++++

type InvokeOperation struct {
	blockEntries []blockEntry
	//Setting where the operation evaluates things while capturing
	evalutesInsteadOfCapturing bool
	//Mainly for macros - Will always capture nodes except for a corresponding ending block
	forcedCapturing bool
	nextCollection  *InvokeOperation
}

func newInvokeOperation() InvokeOperation {
	return InvokeOperation{
		blockEntries:               []blockEntry{},
		evalutesInsteadOfCapturing: false,
		forcedCapturing:            false,
		nextCollection:             nil,
	}
}
