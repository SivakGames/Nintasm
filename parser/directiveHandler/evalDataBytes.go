package directiveHandler

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/romBuilder/addDataToRom"
)

// +++++++++++++++++++++++++

var directiveOperandByteSizes = map[string]int{
	"DB":   1,
	"DW":   2,
	"DWBE": 2,
}

// For .db, .dw, .dwbe
func evalDataBytesOperands(directiveName string, operandList *[]Node) error {
	isBigEndian := directiveName == "DWBE"
	operandByteSize := directiveOperandByteSizes[directiveName]

	for _, operand := range *operandList {
		evalOperand, err := interpreter.EvaluateNode(operand)
		if err != nil {
			err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
			if err != nil {
				return err // ❌❌ CONTINUES Failing!
			}
		}

		err = addDataToRom.AddRawBytesToRom(evalOperand, operandByteSize, isBigEndian)
		if err != nil {
			return err // ❌ Fails
		}
	}
	return nil
}
