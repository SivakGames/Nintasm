package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
)

type Node = operandFactory.Node

func EvaluateDirective(operationTokenEnum enumTokenTypes.Def, directiveName string, operationLabel string, operandList *[]Node) error {

	// Check if end block and if it's actually closing something
	if operationTokenEnum == enumTokenTypes.DIRECTIVE_blockEnd ||
		operationTokenEnum == enumTokenTypes.DIRECTIVE_labeledBlockEnd {
		currentBlockEntries := blockStack.GetCurrentCaptureBlockStack()
		if len(*currentBlockEntries) == 0 {
			return errorHandler.AddNew(enumErrorCodes.DirectiveUnopenedEndBlock, directiveName)
		}

		if !blockStack.CheckIfEndOpMatchesOpeningOp(directiveName) {
			return errorHandler.AddNew(enumErrorCodes.DirectiveUnopenedEndBlock, directiveName)
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
	case enumTokenTypes.DIRECTIVE_invokeKeyVal:
		return evalKv(operandList)
	case enumTokenTypes.DIRECTIVE_setting:
		return evalSettingChange(directiveName, operandList)
	case enumTokenTypes.DIRECTIVE_settingReset:
		return evalSettingReset(directiveName)
	case enumTokenTypes.DIRECTIVE_throw:
		return evalThrow(operandList)
	case enumTokenTypes.DIRECTIVE_deleteSymbol:
		return evalDeleteMacro(operandList)

	case enumTokenTypes.DIRECTIVE_include:
		switch directiveName {
		case "INCLUDE":
			return evalInclude(operandList)
		case "INCBIN":
			return evalIncbin(operandList)
		default:
			panic("🛑 BAD INCLUDE DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_labeled:
		switch directiveName {
		case "FUNC":
			return evalFunc(operationLabel, operandList)
		case "GNSI":
			return evalGnsi(operationLabel, operandList)
		case "RS":
			return evalRs(operationLabel, operandList)
		default:
			panic("🛑 BAD LABELED DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_blockStart:
		switch directiveName {
		case "IF":
			return evalIf(directiveName, operandList)
		case "ELSEIF":
			return evalElseIf(directiveName, operandList)
		case "ELSE":
			return evalElse(directiveName, operandList)
		case "SWITCH":
			return evalSwitch(directiveName, operandList)
		case "CASE":
			return evalCase(directiveName, operandList)
		case "DEFAULT":
			return evalDefault(directiveName, operandList)
			//return evalIf(directiveName, operandList)
		case "IKV":
			return evalIkv(directiveName, operandList)
		case "REPEAT":
			return evalRepeat(directiveName, operandList)
		default:
			panic("🛑 BAD BLOCK START DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_blockEnd:
		switch directiveName {
		case "ENDIF":
			return evalEndIf(operandList)
		case "ENDIKV":
			return evalEndIkv(operandList)
		case "ENDREPEAT":
			return evalEndRepeat()
		case "ENDSWITCH":
			return evalEndSwitch()
		default:
			panic("🛑 BAD BLOCK END DIRECTIVE!!!" + directiveName)
		}

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
			panic("🛑 BAD LABELED BLOCK START DIRECTIVE!!!" + directiveName)
		}

	case enumTokenTypes.DIRECTIVE_labeledBlockEnd:
		switch directiveName {
		case "ENDCHARMAP":
			return evalEndCharmap()
		case "ENDEXPRMAP":
			return evalEndExprmap()
		case "ENDKVM":
			return evalEndKVMacro()
		case "ENDM":
			return evalEndMacro()
		case "ENDNAMESPACE":
			return evalEndNamespace()
		default:
			panic("🛑 BAD LABELED BLOCK END DIRECTIVE!!!" + directiveName)
		}

	default:
		errMsg := fmt.Sprintf("🛑 BAD DIRECTIVE OPERATION TYPE!!! %v", directiveName)
		panic(errMsg)
	}
}

//--------------------------------------------

// Helper for opening labeled block directives.
// Will see if the stack is empty, if the symbol is already defined,
// and will set the current operation label if possible
func ProcessOpenLabelBlock(openBlockLabel string) error {
	var err error

	currentStack := blockStack.GetCurrentCaptureBlockStack()
	if len(*currentStack) > 0 {
		return errorHandler.AddNew(enumErrorCodes.DirectiveNestedLabelBlock) // ❌ Fails
	}
	err = environment.CheckIfAlreadyDefinedInMasterTable(openBlockLabel)
	if err != nil {
		return err // ❌ Fails
	}
	err = blockStack.SetCurrentOperationLabel(openBlockLabel)
	if err != nil {
		return err // ❌ Fails
	}
	return nil
}
