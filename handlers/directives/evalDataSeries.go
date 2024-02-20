package handlerDirective

import (
	"errors"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

// +++++++++++++++++++++++++

// For .ds
func evalDataSeriesOperands(directiveName string, operandList *[]Node) error {
	seriesValue := uint8(0xff)
	repetitionNode := &(*operandList)[0]
	if !(operandFactory.ValidateNodeIsNumeric(repetitionNode) &&
		operandFactory.ValidateNumericNodeIsGTZero(repetitionNode)) {
		return errors.New("DS/PAD directive repeat value must be a number that is > 0")
	}

	repetitionNumber := repetitionNode.AsNumber
	numRepetitions := repetitionNumber

	if len(*operandList) == 2 {
		padNode := &(*operandList)[1]
		if !(operandFactory.ValidateNodeIsNumeric(padNode) &&
			operandFactory.ValidateNumericNodeIsPositive(padNode) &&
			operandFactory.ValidateNumericNodeIs8BitValue(padNode)) {
			return errors.New("DS/PAD directive fill value must be a non-negative 8-bit number")
		}
		seriesValue = uint8(padNode.AsNumber)
	}

	asRomData := make([]uint8, numRepetitions)

	for i := range asRomData {
		asRomData[i] = seriesValue
	}

	err := romBuilder.AddBytesToRom(asRomData)
	return err
}
