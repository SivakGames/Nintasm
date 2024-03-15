package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
)

func evalSwitch(directiveName string, operandList *[]Node) error {
	switchOperand := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsNumeric(switchOperand) &&
		!operandFactory.ValidateNodeIsString(switchOperand) {
		return errorHandler.AddNew(enumErrorCodes.SwitchStatementBadOperand)
	}

	blockStack.PushOntoTopEntry(directiveName, *operandList)
	return nil
}

func evalCase(directiveName string, operandList *[]Node) error {
	err := checkProperCaseDefaultNesting("CASE")
	if err != nil {
		return err
	}

	switchOperand := getOriginalSwitchOperand()
	caseOperand := &(*operandList)[0]

	if switchOperand.NodeType != caseOperand.NodeType {
		return errorHandler.AddNew(enumErrorCodes.SwitchStatementMismatchedCaseType)
	}
	blockStack.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalDefault(directiveName string, operandList *[]Node) error {
	err := checkProperCaseDefaultNesting("DEFAULT")
	if err != nil {
		return err
	}

	blockStack.CreateNewAlternateForTopEntry(directiveName, *operandList)
	return nil
}

func evalEndSwitch() error {
	//currentStackOperation := blockStack.GetCurrentBlockEntry()
	//var trueStatementCapturedLines *[]blockStack.CapturedLine

	//switchOperand := getOriginalSwitchOperand()

	return nil
}

// Certain shared conditions for case and default statements need to be met before they
// are deemed valid
func checkProperCaseDefaultNesting(childOp string) error {
	//Must actually be in a block
	entries := blockStack.GetCurrentBlockEntries()
	if len(*entries) == 0 {
		if childOp == "CASE" {
			return errorHandler.AddNew(enumErrorCodes.CaseNoSwitch)
		}
		return errorHandler.AddNew(enumErrorCodes.DefaultNoSwitch)
	}

	//Block must actually be a switch block
	topBlockOp := blockStack.GetCurrentBlockEntry()
	switchOpName := topBlockOp.BlockOperationName
	if switchOpName != "SWITCH" {
		if childOp == "CASE" {
			return errorHandler.AddNew(enumErrorCodes.CaseNoSwitch)
		}
		return errorHandler.AddNew(enumErrorCodes.DefaultNoSwitch)
	}

	// Default after default also not allowed
	lastOpName := blockStack.GetCurrentBlockEntryOperationName()
	if lastOpName == "DEFAULT" {
		if childOp == "CASE" {
			return errorHandler.AddNew(enumErrorCodes.SwitchStatementCaseAfterDefault)
		}
		return errorHandler.AddNew(enumErrorCodes.SwitchStatementDuplicateDefault)
	}

	return nil
}

func getOriginalSwitchOperand() Node {
	switchBlockOp := blockStack.GetCurrentBlockEntry()
	return switchBlockOp.OperandList[0]
}
