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

type ProcessLineScope map[string]Node

type ProcessLine struct {
	Scope         ProcessLineScope
	CapturedLines []CapturedLine
}

func newProcessLine(scope ProcessLineScope, capturedLines []CapturedLine) ProcessLine {
	return ProcessLine{
		Scope:         scope,
		CapturedLines: capturedLines,
	}
}

// ++++++++++++++++++++++++++++++++++++

type CaptureBlock struct {
	BlockOperationName    string
	CapturedLines         []CapturedLine
	ProcessLines          []ProcessLine
	OperandList           []Node
	AlternateCaptureBlock *CaptureBlock
}

func newCaptureBlock(blockOperationName string, operandList []Node) CaptureBlock {
	return CaptureBlock{
		BlockOperationName:    blockOperationName,
		CapturedLines:         []CapturedLine{},
		OperandList:           operandList,
		AlternateCaptureBlock: nil,
	}
}

// ++++++++++++++++++++++++++++++++++++

type CaptureBlockListNode struct {
	captureBlockStack []CaptureBlock
	//Setting where the operation evaluates things while capturing
	evalutesInsteadOfCapturing bool
	//Mainly for macros - Will always capture nodes except for a corresponding ending block
	forcedCapturing        bool
	overwriteForcedCapture *map[string]bool
	nextNode               *CaptureBlockListNode
}

func newCaptureBlockListNode() CaptureBlockListNode {
	return CaptureBlockListNode{
		captureBlockStack:          []CaptureBlock{},
		evalutesInsteadOfCapturing: false,
		forcedCapturing:            false,
		nextNode:                   nil,
	}
}
