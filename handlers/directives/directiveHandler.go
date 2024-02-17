package handlerDirective

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/handlers/blockStack"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

func Process(operationTokenEnum enumTokenTypes.Def, directiveName string, operationLabel string, operandList *[]Node) error {

	if operationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		operationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd {
		if len(blockStack.Stack) == 0 {
			errMsg := fmt.Sprintf("%v with no opening operation found!", directiveName)
			return errors.New(errMsg)
		}

		if !blockStack.CheckIfEndOpMatchesOpeningOp(directiveName) {
			errMsg := fmt.Sprintf("Non-matching closing block with parent operation, %v", directiveName)
			return errors.New(errMsg)
		}
	}

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
		case "ELSE":
			return evalElse(directiveName, operandList)
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

	case enumTokenTypes.DIRECTIVE_labeledBlockStart:
		switch directiveName {
		case "MACRO":
			return evalMacro(directiveName, operationLabel, operandList)
		default:
			return errors.New("BAD LABELED BLOCK START DIRECTIVE!!!")
		}

	case enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		switch directiveName {
		case "ENDM":
			return evalEndMacro(directiveName, operandList)
		default:
			return errors.New("BAD LABELED BLOCK END DIRECTIVE!!!")
		}

	default:
		errMsg := fmt.Sprintf("BAD DIRECTIVE OPERATION TYPE!!! %v", directiveName)
		return errors.New(errMsg)
	}
}
