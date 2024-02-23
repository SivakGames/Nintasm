package directiveHandler

import (
	"errors"
	"fmt"
	"misc/nintasm/assemble/blockStack"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
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

	//TODO: Remove directive name args being passed where unneeded

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
	case enumTokenTypes.DIRECTIVE_defCharMap:
		return evalDefChar(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_defExprMap:
		return evalDefExpr(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_setting:
		return evalSettingChange(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_invokeKeyVal:
		return evalKv(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_settingReset:
		return evalSettingReset(directiveName)

	case enumTokenTypes.DIRECTIVE_blockStart:
		switch directiveName {
		case "IF":
			return evalIf(directiveName, operandList)
		case "ELSEIF":
			return evalElseIf(directiveName, operandList)
		case "ELSE":
			return evalElse(directiveName, operandList)
		case "IKV":
			return evalIkv(directiveName, operandList)
		case "REPEAT":
			return evalRepeat(directiveName, operandList)
		default:
			return errors.New("BAD BLOCK START DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_blockEnd:
		switch directiveName {
		case "ENDIF":
			return evalEndIf(directiveName, operandList)
		case "ENDIKV":
			return evalEndIkv(directiveName, operandList)
		case "ENDREPEAT":
			return evalEndRepeat(directiveName, operandList)
		default:
			return errors.New("BAD BLOCK END DIRECTIVE!!!" + directiveName)
		}
		//return evalRomBuildingOperands(directiveName, operandList)

	case enumTokenTypes.DIRECTIVE_labeledBlockStart:
		err := ProcessOpenLabelBlock(operationLabel)
		if err != nil {
			return err
		}

		switch directiveName {
		case "CHARMAP":
			return evalCharmap(directiveName, operationLabel, operandList)
		case "EXPRMAP":
			return evalExprmap(directiveName, operationLabel, operandList)
		case "KVMACRO":
			return evalKVMacro(directiveName, operationLabel, operandList)
		case "MACRO":
			return evalMacro(directiveName, operationLabel, operandList)
		case "NAMESPACE":
			return evalNamespace(directiveName, operationLabel, operandList)

		default:
			return errors.New("BAD LABELED BLOCK START DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		switch directiveName {
		case "ENDCHARMAP":
			return evalEndCharmap(directiveName)
		case "ENDEXPRMAP":
			return evalEndExprmap(directiveName)
		case "ENDKVM":
			return evalEndKVMacro(directiveName)
		case "ENDM":
			return evalEndMacro(directiveName)
		case "ENDNAMESPACE":
			return evalEndNamespace(directiveName)
		default:
			return errors.New("BAD LABELED BLOCK END DIRECTIVE!!!" + directiveName)
		}

	default:
		errMsg := fmt.Sprintf("BAD DIRECTIVE OPERATION TYPE!!! %v", directiveName)
		return errors.New(errMsg)
	}
}

//--------------------------------------------

// Helper for opening labeled block directives.
// Will see if the stack is empty, if the symbol is already defined,
// and will set the current operation label if possible
func ProcessOpenLabelBlock(openBlockLabel string) error {
	var err error

	if len(blockStack.Stack) > 0 {
		errMsg := fmt.Sprintf("Cannot define a labeled block when in another block statement!")
		return errors.New(errMsg) // ❌ Fails
	}
	err = environment.CheckIfSymbolAlreadyDefined(openBlockLabel)
	if err != nil {
		return err // ❌ Fails
	}
	err = blockStack.SetCurrentOperationLabel(openBlockLabel)
	if err != nil {
		return err // ❌ Fails
	}
	return nil
}
