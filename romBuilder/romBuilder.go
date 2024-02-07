package romBuilder

import (
	"errors"
	"fmt"
	"log"
	"math"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	"misc/nintasm/parser/operandFactory"
	"unicode/utf8"
)

type Node = operandFactory.Node

// The final ROM that will be built
var romLayout = make([][][]uint8, 0)

var currentRomSegmentIndex = -1
var currentBankIndex = -1
var currentInsertionIndex = -1
var currentOrg = 0x8000

func AddNewRomSegment(totalSize int, bankSize int) error {
	floatResult := float64(totalSize) / float64(bankSize)
	numBanks, deci := math.Modf(floatResult)
	if deci != 0 {
		return errors.New("Bank size is not evenly distributable")
	}

	newSegment := make([][]uint8, int(numBanks))

	for i := range newSegment {
		newSegment[i] = make([]uint8, bankSize)
	}

	romLayout = append(romLayout, newSegment)

	currentRomSegmentIndex = len(romLayout) - 1
	currentBankIndex = 0
	currentInsertionIndex = 0
	return nil
}

func AddBytesToRom(insertions []uint8) error {
	currentBankSegment := GetCurrentBankSegment()

	toInsertSpace := currentInsertionIndex + len(insertions)
	overflowByteTotal := toInsertSpace - len(*currentBankSegment)

	if overflowByteTotal > 0 {
		errMsg := fmt.Sprintf("Will overflow by: %d byte(s) here", overflowByteTotal)
		return errors.New(errMsg)
	}
	for i := range insertions {
		(*currentBankSegment)[currentInsertionIndex] = insertions[i]
		currentInsertionIndex++
	}

	return nil
}

// Convert into bytes for ROM data
func ConvertNodeValueToUInts(node Node, neededBytes int, isBigEndian bool) ([]uint8, error) {
	var lowByte, highByte int = 0, 0

	convertedValue := make([]uint8, 0)

	if !node.Resolved {
		switch neededBytes {
		case 1:
			convertedValue = append(convertedValue, 0)
		case 2:
			convertedValue = append(convertedValue, 0, 0)
		default:
			panic("Something is very wrong with unresolved byte conversion!")
		}
		return convertedValue, nil
	}

	switch node.NodeType {
	case enumNodeTypes.NumericLiteral:
		highByte = (node.AsNumber & 0x0ff00) >> 8
		lowByte = node.AsNumber & 0x000ff

		switch neededBytes {
		case 1:
			if node.AsNumber < -0x000ff || node.AsNumber > 0x000ff {
				return nil, errors.New("Instruction operand for mode must resolve to an 8 bit value")
			}
			convertedValue = append(convertedValue, uint8(lowByte))
		case 2:
			if node.AsNumber < -0x0ffff || node.AsNumber > 0x0ffff {
				return nil, errors.New("Instruction operand for mode must resolve to a 16 bit value")
			}
			if !isBigEndian {
				convertedValue = append(convertedValue, uint8(lowByte))
				convertedValue = append(convertedValue, uint8(highByte))
			} else {
				convertedValue = append(convertedValue, uint8(highByte))
				convertedValue = append(convertedValue, uint8(lowByte))
			}
		default:
			panic("Something is very wrong with numeric byte conversion!")
		}
	case enumNodeTypes.BooleanLiteral:
		if node.AsBool {
			lowByte = 1
		} else {
			lowByte = 0
		}

		switch neededBytes {
		case 1:
			fmt.Println("\x1b[33mWARNING\x1b[0m: Value is boolean; Resolving to", lowByte)
			convertedValue = append(convertedValue, uint8(lowByte))
		case 2:
			return convertedValue, errors.New("Boolean value cannot be used in 16 bit operations")
		default:
			panic("Something is very wrong with boolean byte conversion!")
		}
	case enumNodeTypes.StringLiteral:
		switch neededBytes {
		case 1:
			convertedStringAsBytes := make([]uint8, 0, len(node.NodeValue))
			for _, c := range node.NodeValue {
				runeLen := utf8.RuneLen(c)
				if runeLen > 1 {
					log.Println("\x1b[43m WARN \x1b[0mCharacter", c, "encoding requires more than a single byte. Using", runeLen, "bytes")
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
			return convertedValue, errors.New("String values cannot be used in 16 bit operations")
		default:
			panic("Something is very wrong with string byte conversion!")
		}
	default:
		panic("Something is very wrong with operand conversion!")
	}

	return convertedValue, nil
}

func GetCurrentRomSegment() *[][]uint8 {
	return &romLayout[currentRomSegmentIndex]
}
func GetCurrentBankSegment() *[]uint8 {
	currentRomSegment := GetCurrentRomSegment()
	return &(*currentRomSegment)[currentBankIndex]
}

func GetOrg() int {
	return currentOrg + currentInsertionIndex
}

// TODO: Set ORG upper/lower bounds checks
func SetOrg(newOrg int) {
	currentOrg = newOrg
	return
}
