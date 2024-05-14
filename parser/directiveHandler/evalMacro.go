package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/operandFactory"
	"regexp"
)

var symbolRegex = regexp.MustCompile(`^\\[A-Za-z_][0-9A-Za-z_]*`)

func evalMacro(directiveName string, macroLabel string, operandList *[]Node) error {
	useArguments := []string{}
	if len(*operandList) > 0 {
		firstOperand := (*operandList)[0]
		if operandFactory.ValidateNodeIsNumeric(&firstOperand) {
			if len(*operandList) > 1 {
				return errorHandler.AddNew(enumErrorCodes.Other, "Only 1 numeric value if using numbers")
			}
			if !operandFactory.ValidateNumericNodeIsPositive(&firstOperand) {
				return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
			}

		} else {
			definedArguments := map[string]bool{}
			for _, o := range *operandList {
				if !operandFactory.ValidateNodeIsSubstitutionID(&o) {
					return errorHandler.AddNew(enumErrorCodes.Other, "Must use a substitution type value")
				}
				_, exists := definedArguments[o.NodeValue]
				if exists {
					return errorHandler.AddNew(enumErrorCodes.Other, "Already exists!")
				}
				if !symbolRegex.MatchString(o.NodeValue) {
					return errorHandler.AddNew(enumErrorCodes.Other, "Substition value must start with a letter or underscore after the backslash and be followed by numbers, underscores, or numbers")
				}

				definedArguments[o.NodeValue] = true
				useArguments = append(useArguments, o.NodeValue)
			}
		}
	}

	blockStack.PushCaptureBlock(directiveName, *operandList)
	environment.AddOtherIdentifierToMasterTable(macroLabel, enumSymbolTableTypes.Macro)
	if len(*operandList) > 0 {
		macroTable.AddArgumentsToMacroTable(macroLabel, &useArguments)
	}
	return nil
}

// End the macro definition and add to environment
func evalEndMacro() error {
	macroLabel := blockStack.GetCurrentOperationLabel()
	capturedLines := blockStack.GetCurrentCaptureBlockCapturedLines()
	if len(*capturedLines) == 0 {
		errorHandler.AddNew(enumErrorCodes.BlockIsEmpty) // ⚠️ Warns
	}

	macroTable.AddCapturedLinesToMacro(macroLabel, macroTable.Macro, *capturedLines)
	blockStack.ProcessEndLabeledDirective()
	return nil
}
