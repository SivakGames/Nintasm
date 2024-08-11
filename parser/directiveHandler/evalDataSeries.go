package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

// +++++++++++++++++++++++++

// For .ds
func evalDataSeriesOperands(directiveName string, operandList *[]Node) error {
	seriesValue := uint8(0xff)
	repetitionNode := &(*operandList)[0]
	if !operandFactory.ValidateNodeIsNumeric(repetitionNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsGTZero(repetitionNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGT, 0)
	}

	repetitionNumber := repetitionNode.AsNumber
	numRepetitions := int(repetitionNumber)

	if len(*operandList) == 2 {
		padNode := &(*operandList)[1]
		if !operandFactory.ValidateNodeIsNumeric(padNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
		} else if !operandFactory.ValidateNumericNodeIsPositive(padNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
		} else if !operandFactory.ValidateNumericNodeIs8BitValue(padNode) {
			return errorHandler.AddNew(enumErrorCodes.NodeValueNot8Bit)
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
