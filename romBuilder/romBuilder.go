package romBuilder

import (
	"errors"
	"fmt"
	"math"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

var romSegments = make([][][]uint8, 0)

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

	romSegments = append(romSegments, newSegment)

	currentRomSegmentIndex = len(romSegments) - 1
	currentBankIndex = 0
	currentInsertionIndex = 0
	return nil
}

func AddBytesToRom(insertions []uint8) error {
	currentRomSegment := &romSegments[currentRomSegmentIndex]
	currentBankSegment := &(*currentRomSegment)[currentBankIndex]

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

func ConvertNodeValueToUInts(node Node, needBytes int) ([]uint8, error) {
	if needBytes < 0 {
		return nil, errors.New("Ruh roh")
	}

	convertedValue := make([]uint8, 0)

	switch node.NodeType {
	case enumNodeTypes.NumericLiteral:
		switch needBytes {
		case 1:
			if node.AsNumber < -0x000ff || node.AsNumber > 0x000ff {
				return nil, errors.New("Instruction operand for mode must resolve to an 8 bit value")
			}
			convertedValue = append(convertedValue, uint8(node.AsNumber))
		case 2:
			if node.AsNumber < -0x0ffff || node.AsNumber > 0x0ffff {
				return nil, errors.New("Instruction operand for mode must resolve to a 16 bit value")
			}
			word := uint16(node.AsNumber)
			highByte := (word & 0x0ff00) >> 8
			lowByte := word & 0x000ff

			convertedValue = append(convertedValue, uint8(lowByte))
			convertedValue = append(convertedValue, uint8(highByte))
		}
	case enumNodeTypes.BooleanLiteral:
		if node.AsBool {
			convertedValue = append(convertedValue, 1)
		} else {
			convertedValue = append(convertedValue, 0)
		}
		fmt.Println("\x1b[33mWARNING\x1b[0m: Value is boolean; Resolving to 1 or 0...")
	}

	return convertedValue, nil
}

func GetOrg() int {
	return currentOrg + currentInsertionIndex
}
