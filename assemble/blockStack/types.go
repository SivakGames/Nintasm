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

type captureBlock struct {
	BlockOperationName    string
	CapturedLines         []CapturedLine
	OperandList           []Node
	AlternateCaptureBlock *captureBlock
}

func newCaptureBlock(blockOperationName string, operandList []Node) captureBlock {
	return captureBlock{
		BlockOperationName:    blockOperationName,
		CapturedLines:         []CapturedLine{},
		OperandList:           operandList,
		AlternateCaptureBlock: nil,
	}
}

// ++++++++++++++++++++++++++++++++++++

type CaptureBlockList struct {
	captureBlockStack []captureBlock
	//Setting where the operation evaluates things while capturing
	evalutesInsteadOfCapturing bool
	//Mainly for macros - Will always capture nodes except for a corresponding ending block
	forcedCapturing bool
	nextList        *CaptureBlockList
}

func newCaptureBlockList() CaptureBlockList {
	return CaptureBlockList{
		captureBlockStack:          []captureBlock{},
		evalutesInsteadOfCapturing: false,
		forcedCapturing:            false,
		nextList:                   nil,
	}
}
