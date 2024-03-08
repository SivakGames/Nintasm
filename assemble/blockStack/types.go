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

type BlockOperationStack struct {
	BlockOperationName  string
	OperandList         []Node
	CapturedLines       []CapturedLine
	AlternateStackBlock *BlockOperationStack
}

func newBlockOperationStack(operationName string, operandList []Node) BlockOperationStack {
	return BlockOperationStack{
		BlockOperationName: operationName,
		OperandList:        operandList,
	}
}

// ++++++++++++++++++++++++++++++++++++

type mainStack struct {
	Flag1               bool
	flag2               bool
	blockOperationStack []BlockOperationStack
}

func newMainStack() mainStack {
	return mainStack{
		Flag1:               false,
		flag2:               false,
		blockOperationStack: []BlockOperationStack{},
	}
}
