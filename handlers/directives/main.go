package handlerDirective

import (
	"errors"
	"fmt"
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
		switch directiveName {
		case "IF":
			return evalIf(directiveName, operandList)
		case "ELSEIF":
			return evalElseIf(directiveName, operandList)
		case "REPEAT":
			return evalRepeat(directiveName, operandList)
		default:
			return errors.New("BAD BLOCK START DIRECTIVE!!!")
		}
	case enumTokenTypes.DIRECTIVE_blockEnd:
		switch directiveName {
		case "ENDIF":
			return evalEndIf(directiveName, operandList)
		case "ENDREPEAT":
			return evalEndRepeat(directiveName, operandList)
		default:
			return errors.New("BAD BLOCK END DIRECTIVE!!!")
		}
		//return evalRomBuildingOperands(directiveName, operandList)
	default:
		errMsg := fmt.Sprintf("BAD DIRECTIVE OPERATION TYPE!!! %v", directiveName)
		return errors.New(errMsg)
	}
}
