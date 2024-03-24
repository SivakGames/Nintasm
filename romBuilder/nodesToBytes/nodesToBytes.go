package nodesToBytes

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/operandFactory"
	"unicode/utf8"
)

type Node = operandFactory.Node

//------------------------------------------

// Convert into bytes for ROM data
func ConvertNodeValueToUInts(node Node, operandByteSize int, isBigEndian bool, isExact bool) ([]uint8, error) {
	var lowByte, highByte int = 0, 0

	convertedValue := make([]uint8, 0)

	if !node.Resolved {
		switch operandByteSize {
		case 1:
			convertedValue = append(convertedValue, 0)
		case 2:
			convertedValue = append(convertedValue, 0, 0)
		default:
			panic("🛑 Something is very wrong with unresolved byte conversion!")
		}
		return convertedValue, nil
	}

	switch node.NodeType {
	case enumNodeTypes.NumericLiteral:
		if !operandFactory.ValidateNodeIsInt(&node) {
			errorHandler.AddNew(enumErrorCodes.ResolvedValueNotInt, node.AsNumber, int(node.AsNumber))

		}

		highByte = (int(node.AsNumber) & 0x0ff00) >> 8
		lowByte = int(node.AsNumber) & 0x000ff

		switch operandByteSize {
		case 1:
			if node.AsNumber < -0x000ff || node.AsNumber > 0x000ff {
				return nil, errorHandler.AddNew(enumErrorCodes.ResolvedValueNot8Bit, int(node.AsNumber)) // ❌ Fails
			}
			convertedValue = append(convertedValue, uint8(lowByte))
		case 2:
			if node.AsNumber < -0x0ffff || node.AsNumber > 0x0ffff {
				return nil, errorHandler.AddNew(enumErrorCodes.ResolvedValueNot16Bit, int(node.AsNumber)) // ❌ Fails
			}
			if !isBigEndian {
				convertedValue = append(convertedValue, uint8(lowByte))
				convertedValue = append(convertedValue, uint8(highByte))
			} else {
				convertedValue = append(convertedValue, uint8(highByte))
				convertedValue = append(convertedValue, uint8(lowByte))
			}
		default:
			panic("🛑 Something is very wrong with numeric byte conversion!")
		}
	case enumNodeTypes.BooleanLiteral:
		if node.AsBool {
			lowByte = 1
		} else {
			lowByte = 0
		}

		switch operandByteSize {
		case 1:
			errorHandler.AddNew(enumErrorCodes.ResolvedValueIsBool, lowByte) // ⚠️ Warns
			convertedValue = append(convertedValue, uint8(lowByte))
		case 2:
			return convertedValue, errorHandler.AddNew(enumErrorCodes.ResolvedValue16BitBool) // ❌ Fails
		default:
			panic("🛑 Something is very wrong with boolean byte conversion!")
		}
	case enumNodeTypes.StringLiteral:
		switch operandByteSize {
		case 1:
			convertedStringAsBytes := make([]uint8, 0, len(node.NodeValue))
			for _, c := range node.NodeValue {
				runeLen := utf8.RuneLen(c)
				if runeLen > 1 {
					errorHandler.AddNew(enumErrorCodes.ResolvedValueMultiByteChar, c, runeLen) // ⚠️ Warns
					for i := 0; i < runeLen; i++ {
						writeRune := (rune(c) >> (i * 8)) & 0x000ff
						convertedStringAsBytes = append(convertedStringAsBytes, uint8(writeRune))
					}
				} else {
					convertedStringAsBytes = append(convertedStringAsBytes, uint8(rune(c)))
				}
			}
			convertedValue = append(convertedValue, convertedStringAsBytes...)
		case 2:
			return convertedValue, errorHandler.AddNew(enumErrorCodes.ResolvedValue16BitString) // ❌ Fails
		default:
			panic("🛑 Something is very wrong with string byte conversion!")
		}

	case enumNodeTypes.MultiByte:
		for _, n := range *node.ArgumentList {
			subValue, err := ConvertNodeValueToUInts(n, operandByteSize, isBigEndian, false)
			if err != nil {
				return nil, err
			}
			convertedValue = append(convertedValue, subValue...)
		}

	default:
		panic("🛑 Something is very wrong with operand conversion!")
	}

	//Multibytes are technially OK, but have to fit
	if isExact && len(convertedValue) > operandByteSize {
		return convertedValue, errorHandler.AddNew(enumErrorCodes.ResolvedValueTooBig, len(convertedValue), operandByteSize) // ❌ Fails
	}

	return convertedValue, nil
}
