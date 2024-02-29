package romSegmentation

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

func ValidateAndSetBank(bankNode *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(bankNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	}
	newBankIndex := bankNode.AsNumber
	currentBankIndex := romBuilder.GetBankIndex()
	if currentBankIndex+1 != newBankIndex {
		return errorHandler.AddNew(enumErrorCodes.BankNotSequential)
	}
	totalBanksInRomSegment := romBuilder.GetTotalBanksInCurrentRomSegment()
	if newBankIndex >= totalBanksInRomSegment {
		return errorHandler.AddNew(enumErrorCodes.BankNumberTooHigh)
	}

	romBuilder.SetBankIndex(newBankIndex)

	return nil
}

func ValidateAndSetOrg(orgNode *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(orgNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsPositive(orgNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
	} else if !operandFactory.ValidateNumericNodeIs16BitValue(orgNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNot16Bit)
	}

	newOrg := orgNode.AsNumber
	err := romBuilder.SetOrg(newOrg)

	return err
}
