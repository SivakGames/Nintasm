package util

import (
	enumParserTypes "misc/nintasm/constants/enums/parserTypes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
)

type LineOperationParsedValues struct {
	OperandStartPosition int
	OperationLabel       string
	OperationTokenEnum   enumTokenTypes.Def
	OperationTokenValue  string
	ParentParserEnum     enumParserTypes.Def
}

func NewLineOperationParsedValues(
	operandStartPosition int,
	operationLabel string,
	operationTokenEnum enumTokenTypes.Def,
	operationTokenValue string,
	parentParserEnum enumParserTypes.Def,
) LineOperationParsedValues {
	return LineOperationParsedValues{
		OperandStartPosition: operandStartPosition,
		OperationLabel:       operationLabel,
		OperationTokenEnum:   operationTokenEnum,
		OperationTokenValue:  operationTokenValue,
		ParentParserEnum:     parentParserEnum,
	}
}
