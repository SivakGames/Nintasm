package romSegmentation

import (
	"errors"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

func ValidateAndSetBank(bankNode *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(bankNode) {
		return errors.New("MUST be a number!")
	}
	newBankIndex := bankNode.AsNumber
	currentBankIndex := romBuilder.GetBankIndex()
	if currentBankIndex+1 != newBankIndex {
		return errors.New("Bank declarations must be sequentially incrementing")
	}
	totalBanksInRomSegment := romBuilder.GetTotalBanksInCurrentRomSegment()
	if newBankIndex >= totalBanksInRomSegment {
		return errors.New("Too high bank number")
	}

	romBuilder.SetBankIndex(newBankIndex)

	return nil
}

func ValidateAndSetOrg(orgNode *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(orgNode) {
		return errors.New("MUST be a number!")
	}
	if !operandFactory.ValidateNumericNodeIsPositive(orgNode) {
		return errors.New("ORG MUST be a positive number!")
	}
	if !operandFactory.ValidateNumericNodeIs16BitValue(orgNode) {
		return errors.New("ORG MUST be a 16 bit value")
	}

	newOrg := orgNode.AsNumber
	err := romBuilder.SetOrg(newOrg)

	return err
}
