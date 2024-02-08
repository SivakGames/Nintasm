package romBuilder

import (
	"errors"
	"math"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

// The final ROM that will be built
var romLayout = make([][][]uint8, 0)

var allocatedRomSize = 0

var currentRomSegmentIndex = -1
var currentBankIndex = -1
var CurrentInsertionIndex = -1
var currentOrg = 0x8000

//-------------------------------------------

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
	currentBankIndex = -1
	CurrentInsertionIndex = -1
	return nil
}

// The entire ROM layout (array of segments)
func GetRomLayout() *[][][]uint8 {
	return &romLayout
}

// The current ROM segment (array of bank segments)
func GetCurrentRomSegment() *[][]uint8 {
	return &romLayout[currentRomSegmentIndex]
}

// The current bank segment (array of uint8 )
func GetCurrentBankSegment() *[]uint8 {
	currentRomSegment := GetCurrentRomSegment()
	return &(*currentRomSegment)[currentBankIndex]
}

func SetBank(newBankIndex int) error {
	currentBankIndex = newBankIndex
	return nil
}

func GetOrg() int {
	return currentOrg + CurrentInsertionIndex
}

// TODO: Set ORG upper/lower bounds checks
func SetOrg(newOrg int) {
	currentOrg = newOrg
	return
}
