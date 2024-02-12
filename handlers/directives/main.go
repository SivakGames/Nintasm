package handlerDirective

import (
	"errors"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

func Process(operationTokenEnum enumTokenTypes.Def, directiveName string, operandList *[]Node) error {

	switch operationTokenEnum {

	// 🟢/❌ Could be either
	case enumTokenTypes.DIRECTIVE_dataBytes:
		return evalDataBytesOperands(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_dataSeries:
		return evalDataSeriesOperands(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_mixedData:
		return evalMixedDataBytesOperands(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_INES:
		return evalInesOperands(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_romBuilding:
		return evalRomBuildingOperands(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_blockStart:
		return evalRepeat(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_blockEnd:
		return evalEndRepeat(directiveName, operandList)
		//return evalRomBuildingOperands(directiveName, operandList)
	default:
		return errors.New("BAD DIRECTIVE OPERATION TYPE!!!")
	}
}
